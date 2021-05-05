package vault

import (
	"context"

	vault "github.com/hashicorp/vault/api"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableEngines() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_engines",
		Description: "Vault secrets engines",
		List: &plugin.ListConfig{
			Hydrate: listUser,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getUser,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the secrets engine"},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the secrets engine"},
		},
	}
}

func listUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var vaultClient *vault.Client
	conn, err := connect(ctx)
	if err != nil {
		return nil, err
	}
	opts := &zendesk.UserListOptions{
		PageOptions: zendesk.PageOptions{
			Page:    1,
			PerPage: 100,
		},
	}
	for true {
		users, page, err := conn.GetUsers(ctx, opts)
		if err != nil {
			return nil, err
		}
		for _, t := range users {
			d.StreamListItem(ctx, t)
		}
		if !page.HasNext() {
			break
		}
		opts.Page++
	}
	return nil, nil
}
func getUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		return nil, err
	}
	quals := d.KeyColumnQuals
	plugin.Logger(ctx).Warn("getUser", "quals", quals)
	id := quals["id"].GetInt64Value()
	plugin.Logger(ctx).Warn("getUser", "id", id)
	result, err := conn.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
