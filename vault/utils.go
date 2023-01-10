package vault

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func connect(ctx context.Context, d *plugin.QueryData) (*api.Client, error) {
	addr := os.Getenv("VAULT_ADDR")
	tkn := os.Getenv("VAULT_TOKEN")
	defaultAuthType := "token"

	// In line with the vault CLI, these values can be set through environment variables.
	vaultConfig := GetConfig(d.Connection)

	if vaultConfig.Address == nil {
		vaultConfig.Address = &addr
	}

	if vaultConfig.Token == nil {
		vaultConfig.Token = &tkn
	}

	if vaultConfig.AuthType == nil {
		vaultConfig.AuthType = &defaultAuthType
	}

	if *vaultConfig.Address == "" {
		return nil, errors.New("Vault Address must be set either in VAULT_ADDR environment variable or in connection configuration file.")
	}

	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	var err error

	apiConfig := &api.Config{Address: *vaultConfig.Address, HttpClient: httpClient}
	client, err := api.NewClient(apiConfig)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	if *vaultConfig.AuthType == "token" && *vaultConfig.Token == "" {
		return nil, errors.New("Token must be set via environment or config file when using token auth_type")
	}

	if *vaultConfig.Token != "" {
		client.SetToken(*vaultConfig.Token)
		return client, nil
	}

	switch *vaultConfig.AuthType {
	case "aws":
		return AwsClient(&vaultConfig, client)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown AuthType %s", *vaultConfig.AuthType))
	}
}

// Util func to replace any double / with single ones, used to make concatenating paths easier
func replaceDoubleSlash(url string) string {
	return strings.ReplaceAll(url, "//", "/")
}

// Util func to obtain filtered mounts from all mounts
func filterMounts(in map[string]*api.MountOutput, mountType string) map[string]*api.MountOutput {
	filtered := map[string]*api.MountOutput{}

	for key, mount := range in {
		if mount.Type == mountType {
			filtered[key] = mount
		}
	}

	return filtered
}

// Util func to obtain []string by key from map[string]interface
func getValues(in map[string]interface{}, key string) []string {
	if in[key] == nil {
		return []string{}
	}

	var out []string
	for _, s := range in[key].([]interface{}) {
		out = append(out, fmt.Sprintf("%s", s.(string)))
	}

	return out
}

// Converts and api.Secret object into a slice of strings containing all secret paths
func getSecretAsStrings(s *api.Secret) []string {
	if s == nil || s.Data["keys"] == nil || len(s.Data["keys"].([]interface{})) == 0 {
		return []string{}
	}
	var secrets []string
	for _, s := range s.Data["keys"].([]interface{}) {
		secrets = append(secrets, fmt.Sprintf("%s", s.(string)))

	}
	return secrets
}

// Transforms
func convertTimestamp(_ context.Context, input *transform.TransformData) (interface{}, error) {
	return time.Unix(input.Value.(int64), 0), nil
}
