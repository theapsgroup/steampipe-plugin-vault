package vault

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*api.Client, error) {
	addr := os.Getenv("VAULT_ADDR")
	tkn := os.Getenv("VAULT_TOKEN")

	vaultConfig := GetConfig(d.Connection)
	if &vaultConfig != nil {
		if vaultConfig.Address != nil {
			addr = *vaultConfig.Address
		}
		if vaultConfig.Token != nil {
			tkn = *vaultConfig.Token
		}
	}

	if addr == "" {
		return nil, errors.New("Vault Address must be set either in VAULT_ADDR environment variable or in connection configuration file.")
	}

	if tkn == "" {
		return nil, errors.New("Vault Token must be set either in VAULT_TOKEN environment variable or in connection configuration file.")
	}

	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	client, err := api.NewClient(&api.Config{Address: addr, HttpClient: httpClient})

	if err != nil {
		return nil, errors.New(err.Error())
	}

	client.SetToken(tkn)

	return client, nil
}
