package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// The structure of a KV secret.
// Key is the path within the mountpoint.
// Path is the name of the engine
type KvSecret struct {
	Key  string
	Path string
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
			KeyColumns: plugin.AllColumns([]string{"key", "mountpoint"}),
			Hydrate:    getSecret,
		},
		Columns: []*plugin.Column{
			{Name: "key", Type: proto.ColumnType_STRING, Description: "The key/path of the kv secret"},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path (mount point) of the secrets engine"},
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

// Checks whether a secret exists, used for the get single secret call
// Returns a bool to make sure we don't leak the values of a secret
func secretExists(ctx context.Context, client *api.Client, engine string, keyPath string) (bool, error) {
	data, err := client.Logical().Read(replaceDoubleSlash(fmt.Sprintf("/%s/metadata/%s", engine, keyPath)))
	return data != nil, err
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
				d.StreamListItem(ctx, &KvSecret{Path: path, Key: k})
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

	data, err := secretExists(ctx, conn, mountpoint, keyPath)

	if err != nil {
		return nil, err
	}
	if data {
		return &KvSecret{Path: mountpoint, Key: keyPath}, nil
	}

	return nil, nil
}
