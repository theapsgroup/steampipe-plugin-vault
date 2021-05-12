package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

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
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The path (mount point) of the secrets engine"},
			{Name: "deletion_time", Type: proto.ColumnType_TIMESTAMP, Description: "The path (mount point) of the secrets engine"},
			{Name: "destroyed", Type: proto.ColumnType_BOOL, Description: "The path (mount point) of the secrets engine"},
			{Name: "version", Type: proto.ColumnType_INT, Description: "The path (mount point) of the secrets engine"},
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
// Folders are identified by a trailing slash. Non trailing slash entries are individual secrets
func listKvSecrets(ctx context.Context, client *api.Client, engine string, keyPath string) ([]string, error) {
	var secrets []string
	data, err := client.Logical().List(replaceDoubleSlash(fmt.Sprintf("/%s/metadata/%s", engine, keyPath)))
	for _, k := range getSecretAsStrings(ctx, data) {
		fullPath := replaceDoubleSlash(fmt.Sprintf("%s/%s", keyPath, k))
		if strings.HasSuffix(k, "/") {
			nestedSecrets, _ := listKvSecrets(ctx, client, engine, fullPath)
			secrets = append(secrets, nestedSecrets...)

		} else {
			secrets = append(secrets, fullPath)
		}
	}
	return secrets, err
}

// The function called by steampipe to populate the table. Will recursively fetch all secrets
func listSecrets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	mounts, err := conn.Sys().ListMounts()
	for path := range mounts {
		if mounts[path].Type == "kv" {
			secrets, _ := listKvSecrets(ctx, conn, path, "")
			for _, k := range secrets {
				secret, _ := getSecretMetadata(ctx, conn, path, k)
				d.StreamListItem(ctx, secret)
			}
		}
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
