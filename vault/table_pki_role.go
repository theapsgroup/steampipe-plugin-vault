package vault

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type PkiRole struct {
	Path             string
	Name             string
	AllowAnyName     bool
	AllowIpSans      bool
	AllowLocalhost   bool
	AllowSubDomains  bool
	AllowedDomains   []string
	AllowedUriSans   []string
	AllowedOtherSans []string
	ClientFlag       bool
	CodeSigningFlag  bool
	KeyBits          int64
	KeyType          string
	Ttl              int64
	MaxTtl           int64
	ServerFlag       bool
}

// Table Function
func tablePkiRole() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_pki_role",
		Description: "Vault PKI Roles",
		List: &plugin.ListConfig{
			Hydrate: listPkiRoles,
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path (mount point) of the engine containing PKI roles"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The PKI role"},
			{Name: "allow_any_name", Type: proto.ColumnType_BOOL, Description: "Allow any name"},
			{Name: "allow_ip_sans", Type: proto.ColumnType_BOOL, Description: "Allow IP based subject alternative names", Transform: transform.FromField("AllowIpSans")},
			{Name: "allow_localhost", Type: proto.ColumnType_BOOL, Description: "Allow localhost"},
			{Name: "allow_sub_domains", Type: proto.ColumnType_BOOL, Description: "Allow subdomains"},
			{Name: "allowed_domains", Type: proto.ColumnType_JSON, Description: "Array of allowed domain names"},
			{Name: "allowed_uri_sans", Type: proto.ColumnType_JSON, Description: "Array of allowed URI based subject alternative names", Transform: transform.FromField("AllowedUriSans")},
			{Name: "allowed_other_sans", Type: proto.ColumnType_JSON, Description: "Array of allowed other subject alternative names", Transform: transform.FromField("AllowedOtherSans")},
			{Name: "client_flag", Type: proto.ColumnType_BOOL, Description: "Can generate client based certificates"},
			{Name: "code_signing_flag", Type: proto.ColumnType_BOOL, Description: "Can generate code-signing based certificates"},
			{Name: "key_bits", Type: proto.ColumnType_INT, Description: "Length of key in bits"},
			{Name: "key_type", Type: proto.ColumnType_STRING, Description: "Type of key used, example 'RSA'"},
			{Name: "ttl", Type: proto.ColumnType_INT, Description: "Default TTL", Transform: transform.FromField("Ttl")},
			{Name: "max_ttl", Type: proto.ColumnType_INT, Description: "Maximum TTL", Transform: transform.FromField("MaxTtl")},
			{Name: "server_flag", Type: proto.ColumnType_BOOL, Description: "Can generate server based certificates"},
		},
	}
}

// Hydrate Functions
func listPkiRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	allMounts, err := conn.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	mounts := filterMounts(allMounts, "pki")
	for mount := range mounts {
		roles, err := getPkiRoleDetails(ctx, conn, mount)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			d.StreamListItem(ctx, role)
		}
	}
	return nil, nil
}

// Data Obtaining Functions
func getPkiRoleDetails(ctx context.Context, client *api.Client, engine string) ([]*PkiRole, error) {
	data, err := client.Logical().List(replaceDoubleSlash(fmt.Sprintf("/%s/roles", engine)))
	if err != nil {
		return nil, err
	}

	out := []*PkiRole{}

	roles := getSecretAsStrings(data)
	for _, role := range roles {
		roleDetails, err := client.Logical().Read(replaceDoubleSlash(fmt.Sprintf("/%s/roles/%s", engine, role)))
		if err != nil {
			return nil, err
		}

		r := &PkiRole{Path: engine, Name: role}
		r.AllowAnyName = roleDetails.Data["allow_any_name"].(bool)
		r.AllowIpSans = roleDetails.Data["allow_ip_sans"].(bool)
		r.AllowLocalhost = roleDetails.Data["allow_localhost"].(bool)
		r.AllowSubDomains = roleDetails.Data["allow_subdomains"].(bool)
		r.AllowedDomains = getValues(roleDetails.Data, "allowed_domains")
		r.AllowedUriSans = getValues(roleDetails.Data, "allowed_uri_sans")
		r.AllowedOtherSans = getValues(roleDetails.Data, "allowed_other_sans")
		r.ClientFlag = roleDetails.Data["client_flag"].(bool)
		r.CodeSigningFlag = roleDetails.Data["code_signing_flag"].(bool)
		r.KeyBits, _ = roleDetails.Data["key_bits"].(json.Number).Int64()
		r.KeyType = roleDetails.Data["key_type"].(string)
		r.Ttl, _ = roleDetails.Data["ttl"].(json.Number).Int64()
		r.MaxTtl, _ = roleDetails.Data["max_ttl"].(json.Number).Int64()

		out = append(out, r)
	}

	return out, nil
}
