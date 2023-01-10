package vault

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-vault",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"vault_engine":       tableEngine(),
			"vault_kv_secret":    tableKvSecret(),
			"vault_sys_health":   tableSysHealth(),
			"vault_aws_role":     tableAwsRole(),
			"vault_pki_cert":     tablePkiCert(),
			"vault_pki_role":     tablePkiRole(),
			"vault_auth":         tableAuth(),
			"vault_azure_config": tableAzureConfig(),
			"vault_azure_role":   tableAzureRole(),
		},
	}

	return p
}
