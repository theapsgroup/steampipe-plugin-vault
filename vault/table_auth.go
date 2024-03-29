package vault

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type AuthMethod struct {
	Path                  string
	Type                  string
	Description           string
	Accessor              string
	Local                 bool
	SealWrap              bool
	ExternalEntropyAccess bool
	DefaultTtl            int
	MaxTtl                int
	RequestHeaders        []string
	PluginVersion         string
	DeprecationStatus     string
	Options               map[string]string
}

func tableAuth() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_auth",
		Description: "Vault Authentication Methods",
		List: &plugin.ListConfig{
			Hydrate: listAuth,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("path"),
			Hydrate:    getAuth,
		},
		Columns: authColumns(),
	}
}

func listAuth(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	auths, err := conn.Sys().ListAuth()
	if err != nil {
		return nil, err
	}

	for path, auth := range auths {
		d.StreamListItem(ctx, &AuthMethod{
			Path:                  path,
			Type:                  auth.Type,
			Description:           auth.Description,
			Accessor:              auth.Accessor,
			Local:                 auth.Local,
			SealWrap:              auth.SealWrap,
			ExternalEntropyAccess: auth.ExternalEntropyAccess,
			DefaultTtl:            auth.Config.DefaultLeaseTTL,
			MaxTtl:                auth.Config.MaxLeaseTTL,
			RequestHeaders:        auth.Config.PassthroughRequestHeaders,
			PluginVersion:         auth.PluginVersion,
			DeprecationStatus:     auth.DeprecationStatus,
			Options:               auth.Options,
		})
	}

	return nil, nil
}

func getAuth(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	auths, err := conn.Sys().ListAuth()
	if err != nil {
		return nil, err
	}

	q := d.EqualsQuals
	path := q["path"].GetStringValue()

	auth := auths[path]
	if auth == nil {
		return nil, nil
	}

	return &AuthMethod{
		Path:                  path,
		Type:                  auth.Type,
		Description:           auth.Description,
		Accessor:              auth.Accessor,
		Local:                 auth.Local,
		SealWrap:              auth.SealWrap,
		ExternalEntropyAccess: auth.ExternalEntropyAccess,
		DefaultTtl:            auth.Config.DefaultLeaseTTL,
		MaxTtl:                auth.Config.MaxLeaseTTL,
		RequestHeaders:        auth.Config.PassthroughRequestHeaders,
		PluginVersion:         auth.PluginVersion,
		DeprecationStatus:     auth.DeprecationStatus,
		Options:               auth.Options,
	}, nil
}

func authColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "path",
			Type:        proto.ColumnType_STRING,
			Description: "The path (mount point) of the authentication method",
		},
		{
			Name:        "type",
			Type:        proto.ColumnType_STRING,
			Description: "The type of authentication method",
		},
		{
			Name:        "description",
			Type:        proto.ColumnType_STRING,
			Description: "Description associated to the authentication method",
		},
		{
			Name:        "accessor",
			Type:        proto.ColumnType_STRING,
			Description: "The accessor used by authentication method",
		},
		{
			Name:        "local",
			Type:        proto.ColumnType_BOOL,
			Description: "Local Auth Methods are not replicated across clusters",
		},
		{
			Name:        "seal_wrap",
			Type:        proto.ColumnType_BOOL,
			Description: "Is running seal wrap (https://www.vaultproject.io/docs/enterprise/sealwrap)",
		},
		{
			Name:        "external_entropy_access",
			Type:        proto.ColumnType_BOOL,
			Description: "Does Authentication Method have access to Vaults external entropy source",
		},
		{
			Name:        "default_ttl",
			Type:        proto.ColumnType_INT,
			Description: "Default TTL",
		},
		{
			Name:        "max_ttl",
			Type:        proto.ColumnType_INT,
			Description: "Max TTL",
		},
		{
			Name:        "request_headers",
			Type:        proto.ColumnType_JSON,
			Description: "Allowed Pass-Through Request Headers",
			Transform:   transform.FromField("RequestHeaders"),
		},
		{
			Name:        "plugin_version",
			Type:        proto.ColumnType_STRING,
			Description: "Information about the plugin used for the authentication method",
		},
		{
			Name:        "deprecation_status",
			Type:        proto.ColumnType_STRING,
			Description: "Deprecation status of the authentication method",
		},
		{
			Name:        "options",
			Type:        proto.ColumnType_JSON,
			Description: "The option configuration associated with the authentication method",
		},
	}
}
