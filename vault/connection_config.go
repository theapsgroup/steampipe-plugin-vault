package vault

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type vaultConfig struct {
	Address *string `cty:address`
	Token   *string `cty:token`
}

var ConfigSchema = map[string]*schema.Attribute{
	"address": {
		Type: schema.TypeString,
	},
	"token": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &vaultConfig{}
}

func GetConfig(connection *plugin.Connection) vaultConfig {
	if connection == nil || connection.Config == nil {
		return vaultConfig{}
	}

	config, _ := connection.Config.(vaultConfig)
	return config
}
