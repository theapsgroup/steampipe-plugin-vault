package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// KvSecret The structure of a KV secret.
// Mount is the mount point of the engine
// Path is the path of the secret within the mount point.
// Name is the name of the secret within the mount point, under the path.
type KvSecret struct {
	Mount        string
	Path         string
	Name         string
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
			KeyColumns: []*plugin.KeyColumn{
				{Name: "mount", Require: plugin.Optional, Operators: []string{"=", "~~"}},
				{Name: "path", Require: plugin.Optional, Operators: []string{"=", "~~"}},
			},
			Hydrate: listSecrets,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"mount", "path", "name"}),
			Hydrate:    getSecret,
		},
		Columns: []*plugin.Column{
			{Name: "mount", Type: proto.ColumnType_STRING, Description: "The mount point of the secrets engine"},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path of the kv secret"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the kv secret"},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time the secret was created"},
			{Name: "deletion_time", Type: proto.ColumnType_TIMESTAMP, Description: "The date and time the secret was destroyed, if destroyed"},
			{Name: "destroyed", Type: proto.ColumnType_BOOL, Description: "Whether the secret was destroyed"},
			{Name: "version", Type: proto.ColumnType_INT, Description: "The current version of the secret"},
		},
	}
}

// Returns the metadata of a secret, or nil if no secret was found
func getSecretMetadata(ctx context.Context, client *api.Client, secret *KvSecret) (*KvSecret, error) {
	logger := plugin.Logger(ctx)

	metadataUrl := fmt.Sprintf("/%smetadata/%s%s", secret.Mount, secret.Path, secret.Name)
	logger.Debug("vault_kv_secret: getSecretMetadata", "metadataUrl", metadataUrl)
	metadata, err := client.Logical().ReadWithContext(ctx, metadataUrl)
	logger.Debug("vault_kv_secret: getSecretMetadata", "metadata", fmt.Sprintf("%#v", metadata))

	if err != nil {
		return nil, err
	}
	if metadata == nil {
		return nil, nil
	}

	createdTime, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%s", metadata.Data["created_time"]))
	if err == nil {
		secret.CreatedTime = createdTime
	}

	deletionTime, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%s", metadata.Data["deletion_time"]))
	if err == nil {
		secret.DeletionTime = deletionTime
	}

	secret.Version, _ = metadata.Data["current_version"].(json.Number).Int64()
	// The returned structure contains a map of versions and their properties. E.g. {..., "versions": { "1": { "destroyed": true } } }
	// This line walks tha tree and fetches the "destroyed" property of the current version
	secret.Destroyed = metadata.Data["versions"].(map[string]interface{})[fmt.Sprintf("%d", secret.Version)].(map[string]interface{})["destroyed"].(bool)

	return secret, nil
}

// Lists all secrets in a secret engine, this has to be done recursively because you only get everything in a "folder"
func listPathSecrets(ctx context.Context, client *api.Client, secret *KvSecret) ([]*KvSecret, error) {
	logger := plugin.Logger(ctx)

	var secrets []*KvSecret
	uri := fmt.Sprintf("/%smetadata/%s%s", secret.Mount, secret.Path, secret.Name)
	logger.Debug("vault_kv_secret: listPathSecrets", "uri", uri)
	data, err := client.Logical().List(uri)
	for _, k := range getSecretAsStrings(data) {
		secrets = append(secrets, &KvSecret{Mount: secret.Mount, Path: fmt.Sprintf("%s%s", secret.Path, secret.Name), Name: k})
	}
	return secrets, err
}

// Worker to receive paths to explore. Folders are explored recursively
// Folders are identified by a trailing slash. Non trailing slash entries are individual secrets
// foldersChan is the channel that will be used to receive paths to still explore from. This is fed by this function as well as the listSecrets one
// secretsChan is the channel that will be used to output received secret metadata, which is the data we're actually interested in
func listKvSecrets(ctx context.Context, client *api.Client, foldersChan chan *KvSecret, secretsChan chan *KvSecret, wg *sync.WaitGroup) {
	logger := plugin.Logger(ctx)

	for k := range foldersChan {
		logger.Debug("vault_kv_secret: listKvSecrets", "k", fmt.Sprintf("%#v", k))
		// time.Sleep(time.Second)

		if k.Name == "" || strings.HasSuffix(k.Name, "/") {
			pathSecrets, _ := listPathSecrets(ctx, client, k)
			logger.Debug("vault_kv_secret: listKvSecrets", "pathSecrets", fmt.Sprintf("%#v", pathSecrets))

			// We use the waitgroup as a counter. Once we've had as many wg.Done() calls as wg.Add, we've processed all trees
			wg.Add(len(pathSecrets))
			for _, p := range pathSecrets {
				foldersChan <- p
			}
		} else {
			secret, err := getSecretMetadata(ctx, client, k)
			if err == nil {
				secretsChan <- secret
			}
		}
		wg.Done()
	}
}

// The function called by steampipe to populate the table. Will recursively fetch all secrets
func listSecrets(ctx context.Context, d *plugin.QueryData, hd *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// used to determine when we've explored all paths
	var wg sync.WaitGroup

	// Fairly large buffers. Because foldersChan is self feeding with recursive paths, it could deadlock if there isn't
	// enough space to actually contain the paths left to explore
	foldersChan := make(chan *KvSecret, 50000)
	secretsChan := make(chan *KvSecret, 50000)

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	// Queue up the mounts to explore
	allMounts, err := conn.Sys().ListMounts()
	if err != nil {
		return nil, err
	}
	logger.Debug("vault_kv_secret: listSecrets", "allMounts", fmt.Sprintf("%#v", allMounts))

	mounts := filterMounts(ctx, allMounts, "kv", d.Quals)
	logger.Debug("vault_kv_secret: listSecrets", "mounts", fmt.Sprintf("%#v", mounts))
	for m := range mounts {
		logger.Debug("vault_kv_secret: listSecrets", "m", m)
		for _, pathQual := range d.Quals["path"].Quals {
			logger.Debug("vault_kv_secret: listSecrets", "pathQual", pathQual)

			path := pathQual.Value.GetStringValue()
			switch pathQual.Operator {
			case "=":
			case "~~":
				path, _, _ = strings.Cut(pathQual.Value.GetStringValue(), "%")
			}

			foldersChan <- &KvSecret{Mount: m, Path: path}
			wg.Add(1)
		}
	}

	// Workers for parallel requests
	go listKvSecrets(ctx, conn, foldersChan, secretsChan, &wg)
	// go listKvSecrets(ctx, conn, foldersChan, secretsChan, &wg)
	// go listKvSecrets(ctx, conn, foldersChan, secretsChan, &wg)
	// go listKvSecrets(ctx, conn, foldersChan, secretsChan, &wg)

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
	logger := plugin.Logger(ctx)

	conn, err := connect(ctx, d)

	if err != nil {
		return nil, err
	}

	quals := d.EqualsQuals
	name := quals["name"].GetStringValue()
	path := quals["path"].GetStringValue()
	mount := quals["mount"].GetStringValue()

	logger.Debug("vault_kv_secret: getSecret", "mount", mount)
	logger.Debug("vault_kv_secret: getSecret", "path", path)
	logger.Debug("vault_kv_secret: getSecret", "name", name)
	secret := KvSecret{
		Mount: mount,
		Path:  path,
		Name:  name,
	}
	data, err := getSecretMetadata(ctx, conn, &secret)

	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return data, nil
}
