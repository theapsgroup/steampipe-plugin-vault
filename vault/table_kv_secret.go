package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type SecretPath struct {
	Engine string
	Path   string
}

// The structure of a KV secret.
// Key is the path within the mountpoint.
// Path is the name of the engine
type KvSecret struct {
	Key          string
	Path         string
	CreatedTime  time.Time
	DeletionTime time.Time
	Destroyed    bool
	Version      int64
}

// Defines the table structure and functions to get vault kv secret data
func tableKvSecret() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_kv_secret",
		Description: "Vault kv secret keys",
		List: &plugin.ListConfig{
			Hydrate: listSecrets,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"key", "path"}),
			Hydrate:    getSecret,
		},
		Columns: []*plugin.Column{
			{Name: "key", Type: proto.ColumnType_STRING, Description: "The key/path of the kv secret"},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path (mount point) of the secrets engine"},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time the secret was created"},
			{Name: "deletion_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time the secret was destroyed, if destroyed"},
			{Name: "destroyed", Type: proto.ColumnType_BOOL, Description: "Whether the secret was destroyed"},
			{Name: "version", Type: proto.ColumnType_INT, Description: "The current version of the secret"},
		},
	}
}

// Converts and api.Secret object into a slice of strings containing all secret paths
func getSecretAsStrings(ctx context.Context, s *api.Secret) []string {
	if s == nil || s.Data["keys"] == nil || len(s.Data["keys"].([]interface{})) == 0 {
		return []string{}
	}
	var secrets []string
	for _, s := range s.Data["keys"].([]interface{}) {
		secrets = append(secrets, fmt.Sprintf("%s", s.(string)))

	}
	return secrets
}

// Returns the metadata of a secret, or nil if no secret was found
func getSecretMetadata(ctx context.Context, client *api.Client, engine string, keyPath string) (*KvSecret, error) {
	data, err := client.Logical().Read(replaceDoubleSlash(fmt.Sprintf("/%s/metadata/%s", engine, keyPath)))

	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	secret := &KvSecret{Path: engine, Key: keyPath}

	createdTime, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%s", data.Data["created_time"]))
	if err == nil {
		secret.CreatedTime = createdTime
	}

	deletionTime, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%s", data.Data["deletion_time"]))
	if err == nil {
		secret.DeletionTime = deletionTime
	}

	secret.Version, _ = data.Data["current_version"].(json.Number).Int64()
	// The returned structure contains a map of versions and their properties. E.g. {..., "versions": { "1": { "destroyed": true } } }
	// This line walks tha tree and fetches the "destroyed" property of the current version
	secret.Destroyed = data.Data["versions"].(map[string]interface{})[fmt.Sprintf("%d", secret.Version)].(map[string]interface{})["destroyed"].(bool)

	return secret, nil
}

// Lists all secrets in a secret engine, this has to be done recursively because you only get everything in a "folder"
func listPathSecrets(ctx context.Context, client *api.Client, engine string, keyPath string) ([]string, error) {
	var secrets []string
	data, err := client.Logical().List(replaceDoubleSlash(fmt.Sprintf("/%s/metadata/%s", engine, keyPath)))
	for _, k := range getSecretAsStrings(ctx, data) {
		fullPath := replaceDoubleSlash(fmt.Sprintf("%s/%s", keyPath, k))
		secrets = append(secrets, fullPath)
	}
	return secrets, err
}

// Worker to receive paths to explore. Folders are explored recursively
// Folders are identified by a trailing slash. Non trailing slash entries are individual secrets
// foldersChan is the channel that will be used to receive paths to still explore from. This is fed by this function as well as the listSecrets one
// secretsChan is the channel that will be used to output received secret metadata, which is the data we're actually interested in
func listKvSecrets(ctx context.Context, client *api.Client, foldersChan chan SecretPath, secretsChan chan *KvSecret, wg *sync.WaitGroup) {
	for k := range foldersChan {
		if strings.HasSuffix(k.Path, "/") {
			pathSecrets, _ := listPathSecrets(ctx, client, k.Engine, k.Path)

			// We use the waitgroup as a counter. Once we've had as many wg.Done() calls as wg.Add, we've processed all trees
			wg.Add(len(pathSecrets))
			for _, p := range pathSecrets {
				foldersChan <- SecretPath{Engine: k.Engine, Path: p}
			}
		} else {
			secret, err := getSecretMetadata(ctx, client, k.Engine, k.Path)
			if err == nil {
				secretsChan <- secret
			}
		}
		wg.Done()
	}
}

// The function called by steampipe to populate the table. Will recursively fetch all secrets
func listSecrets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// used to determine when we've explored all paths
	var wg sync.WaitGroup

	// Fairly large buffers. Because foldersChan is self feeding with recursive paths, it could deadlock if there isn't
	// enough space to actually contain the paths left to explore
	foldersChan := make(chan SecretPath, 50000)
	secretsChan := make(chan *KvSecret, 50000)

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	// Queue up the mounts to explore
	mounts, err := conn.Sys().ListMounts()
	for path := range mounts {
		if mounts[path].Type == "kv" {
			foldersChan <- SecretPath{Engine: path, Path: "/"}
			wg.Add(1)
		}
	}

	// Workers for parallel requests
	go listKvSecrets(ctx, conn, foldersChan, secretsChan, &wg)
	go listKvSecrets(ctx, conn, foldersChan, secretsChan, &wg)
	go listKvSecrets(ctx, conn, foldersChan, secretsChan, &wg)
	go listKvSecrets(ctx, conn, foldersChan, secretsChan, &wg)

	// Wait for the waitgroup to be done, once the waitgroup is done we'll have explored all paths and can close the channels
	// This makes the goroutines and stream loop below terminate
	go func() {
		wg.Wait()
		close(foldersChan)
		close(secretsChan)
	}()

	// Stream any items to steampipe that we've received so far
	for s := range secretsChan {
		d.StreamListItem(ctx, s)
	}

	return nil, nil
}

// Fetches a single secret, essentially just a check whether it exists.
func getSecret(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)

	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	keyPath := quals["key"].GetStringValue()
	mountpoint := quals["path"].GetStringValue()

	data, err := getSecretMetadata(ctx, conn, mountpoint, keyPath)

	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return data, nil
}
