package vault

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

type AzureRole struct {
	Path string
	Role string
}

func tableAzureRole() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_azure_role",
		Description: "Vault Azure Configurations",
		List: &plugin.ListConfig{
			Hydrate: listAzureRoles,
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path (mount point) of the Azure Engine"},
			{Name: "role", Type: proto.ColumnType_STRING, Description: "The Azure Role"},
		},
	}
}

func listAzureRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		data, err := conn.Logical().List(replaceDoubleSlash(fmt.Sprintf("/%s/roles", path)))
		if err != nil {
			return nil, err
		}
		roles := getSecretAsStrings(data)
		for _, role := range roles {
			d.StreamListItem(ctx, &AzureRole{
				Path: path,
				Role: role,
			})
		}
	}

	return nil, nil
}
