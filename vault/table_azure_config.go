package vault

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

type AzureConfig struct {
	Path           string
	SubscriptionId string
	TenantId       string
	ClientId       string
	Environment    string
}

func tableAzureConfig() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_azure_config",
		Description: "Vault Azure Configurations",
		List: &plugin.ListConfig{
			Hydrate: listAzureConfigs,
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path (mount point) of the Azure Engine"},
			{Name: "subscription_id", Type: proto.ColumnType_STRING, Description: "The Azure subscription identifier"},
			{Name: "tenant_id", Type: proto.ColumnType_STRING, Description: "The Azure tenant identifier"},
			{Name: "client_id", Type: proto.ColumnType_STRING, Description: "The Azure client identifier"},
			{Name: "environment", Type: proto.ColumnType_STRING, Description: "The Azure environment description"},
		},
	}
}

func listAzureConfigs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	allMounts, err := conn.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	mounts := filterMounts(allMounts, "azure")
	for path := range mounts {
		config, err := conn.Logical().Read(replaceDoubleSlash(fmt.Sprintf("/%s/config", path)))
		if err != nil {
			return nil, err
		}

		d.StreamListItem(ctx, &AzureConfig{
			Path:           path,
			SubscriptionId: config.Data["subscription_id"].(string),
			TenantId:       config.Data["tenant_id"].(string),
			ClientId:       config.Data["client_id"].(string),
			Environment:    config.Data["environment"].(string),
		})
	}

	return nil, nil
}
