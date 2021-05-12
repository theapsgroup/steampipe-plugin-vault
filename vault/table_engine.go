package vault

import (
	"context"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type Engine struct {
	Path        string
	Type        string
	Description string
	Accessor    string
	Version     int64
	Local       bool
	SealWrap    bool
	DefaultTtl  int
	MaxTtl      int
}

func tableEngine() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_engine",
		Description: "Vault secrets engines",
		List: &plugin.ListConfig{
			Hydrate: listEngines,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("path"),
			Hydrate:    getEngine,
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path (mount point) of the secrets engine"},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the secrets engine"},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description associated to mounted engine"},
			{Name: "accessor", Type: proto.ColumnType_STRING, Description: "The accessor used by the secrets engine"},
			{Name: "version", Type: proto.ColumnType_INT, Description: "The secrets engine version"},
			{Name: "local", Type: proto.ColumnType_BOOL, Description: "Is Local Mount (Local mounts are not replicated across clusters)"},
			{Name: "seal_wrap", Type: proto.ColumnType_BOOL, Description: "Is the secrets engine running seal wrap (https://www.vaultproject.io/docs/enterprise/sealwrap)"},
			{Name: "default_ttl", Type: proto.ColumnType_INT, Description: "Default TTL of Secrets within Engine"},
			{Name: "max_ttl", Type: proto.ColumnType_INT, Description: "Max TTL of Secrets within Engine"},
		},
	}
}

func listEngines(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)

	if err != nil {
		return nil, err
	}

	data, err := conn.Sys().ListMounts()
	for path := range data {
		ver, err := strconv.ParseInt(data[path].Options["version"], 0, 32)
		if err != nil {
			ver = 0
		}

		d.StreamListItem(ctx, &Engine{
			Type:        data[path].Type,
			Path:        path,
			Description: data[path].Description,
			Accessor:    data[path].Accessor,
			Version:     ver,
			Local:       data[path].Local,
			SealWrap:    data[path].SealWrap,
			DefaultTtl:  data[path].Config.DefaultLeaseTTL,
			MaxTtl:      data[path].Config.MaxLeaseTTL})
	}

	return nil, nil
}

func getEngine(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)

	if err != nil {
		return nil, err
	}

	data, err := conn.Sys().ListMounts()

	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	path := quals["path"].GetStringValue()

	result := data[path]
	if result == nil {
		// TODO: figure out if this is expected to be error
		return nil, nil
	}

	ver, err := strconv.ParseInt(data[path].Options["version"], 0, 32)
	if err != nil {
		ver = 0
	}

	return &Engine{
		Type:        data[path].Type,
		Path:        path,
		Description: data[path].Description,
		Accessor:    data[path].Accessor,
		Version:     ver,
		Local:       data[path].Local,
		SealWrap:    data[path].SealWrap,
		DefaultTtl:  data[path].Config.DefaultLeaseTTL,
		MaxTtl:      data[path].Config.MaxLeaseTTL}, nil
}
