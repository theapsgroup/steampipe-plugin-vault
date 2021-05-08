package vault

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

type AwsRole struct {
	Mountpoint string
	Role       string
}

type SecretData struct {
	Keys []string `json:"keys"`
}

func tableAwsRoles() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_aws_roles",
		Description: "Vault AWS Roles",
		List: &plugin.ListConfig{
			Hydrate: listRoles,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"mountpoint", "role"}),
			Hydrate:    getRole,
		},
		Columns: []*plugin.Column{
			{Name: "mountpoint", Type: proto.ColumnType_STRING, Description: "The mountpoint of the AWS Roles"},
			{Name: "role", Type: proto.ColumnType_STRING, Description: "The AWS Role"},
		},
	}
}

// Function called by Steampipe to populate the table.
func listRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	mounts, err := getAwsMounts(conn.Sys().ListMounts())
	if err != nil {
		return nil, err
	}

	for mount := range mounts {
		roles, err := listAwsRoles(conn, mount)
		if err != nil {
			return nil, err
		}

		for _, role := range roles.Keys {
			d.StreamListItem(ctx, &AwsRole{Mountpoint: mount, Role: role})
		}
	}

	return nil, nil
}

// Fetches a single role, essentially a check on if it exists
func getRole(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	mountpoint := quals["mountpoint"].GetStringValue()
	role := quals["role"].GetStringValue()

	data, err := roleExists(conn, mountpoint, role)
	if err != nil {
		return nil, err
	}
	if data {
		return &AwsRole{Mountpoint: mountpoint, Role: role}, nil
	}

	return nil, nil
}

// Filter mounts for those of type 'aws'
func getAwsMounts(allMounts map[string]*api.MountOutput, err error) (map[string]*api.MountOutput, error) {
	if err != nil {
		return nil, err
	}
	filtered := map[string]*api.MountOutput{}

	for i, mount := range allMounts {
		if mount.Type == "aws" {
			filtered[i] = mount
		}
	}

	return filtered, nil
}

func listAwsRoles(client *api.Client, engine string) (SecretData, error) {
	out := SecretData{}

	data, err := client.Logical().List(replaceDoubleSlash(fmt.Sprintf("/%s/roles", engine)))
	if err != nil {
		return SecretData{}, err
	}

	b, _ := json.Marshal(data.Data)
	_ = json.Unmarshal([]byte(b), &out)
	return out, nil
}

func roleExists(client *api.Client, mountpoint string, role string) (bool, error) {
	data, err := client.Logical().Read(replaceDoubleSlash(fmt.Sprintf("/%s/roles/%s", mountpoint, role)))
	return data != nil, err
}
