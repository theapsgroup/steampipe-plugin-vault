package vault

import (
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/schema"
)

type vaultConfig struct {
	Address     *string `cty:"address"`
	Token       *string `cty:"token"`
	AuthType    *string `cty:"auth_type"`
	AwsProvider *string `cty:"aws_provider"`
	AwsRole     *string `cty:"aws_role"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"address": {
		Type: schema.TypeString,
	},
	"auth_type": {
		Type: schema.TypeString,
	},
	"token": {
		Type: schema.TypeString,
	},
	"aws_provider": {
		Type: schema.TypeString,
	},
	"aws_role": {
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
