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
	Path                   string
	Role                   string
	CredentialType         string
	DefaultStsTtl          int64
	MaxStsTtl              int64
	PolicyDocument         string
	UserPath               string
	PermissionsBoundaryArn string
	RoleArns               []string
	PolicyArns             []string
	IamGroups              []string
}

func tableAwsRole() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_aws_role",
		Description: "Vault AWS Roles",
		List: &plugin.ListConfig{
			Hydrate: listRoles,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"path", "role"}),
			Hydrate:    getRole,
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path (mount point) of the engine containing AWS Roles"},
			{Name: "role", Type: proto.ColumnType_STRING, Description: "The AWS Role"},
			{Name: "credential_type", Type: proto.ColumnType_STRING, Description: "The type of Credential assumed_role, iam, etc"},
			{Name: "default_sts_ttl", Type: proto.ColumnType_INT, Description: "Default STS TTL"},
			{Name: "max_sts_ttl", Type: proto.ColumnType_INT, Description: "Maximum STS TTL"},
			{Name: "policy_document", Type: proto.ColumnType_JSON, Description: "AWS Policy Document associated with the Role"},
			{Name: "user_path", Type: proto.ColumnType_JSON, Description: "Path of User"},
			{Name: "permissions_boundary_arn", Type: proto.ColumnType_JSON, Description: "ARN of the Permissions Boundary"},
			{Name: "role_arns", Type: proto.ColumnType_JSON, Description: "ARNs associated with the Role"},
			{Name: "policy_arns", Type: proto.ColumnType_JSON, Description: "ARNs associated with the Policies of the Role"},
			{Name: "iam_groups", Type: proto.ColumnType_JSON, Description: "IAM groups associated with the Role"},
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

		for _, r := range roles {
			role, _ := getRoleDetails(conn, mount, r)
			d.StreamListItem(ctx, role)
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
	mountpoint := quals["path"].GetStringValue()
	role := quals["role"].GetStringValue()

	data, err := getRoleDetails(conn, mountpoint, role)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	return data, nil
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

func listAwsRoles(client *api.Client, engine string) ([]string, error) {
	data, err := client.Logical().List(replaceDoubleSlash(fmt.Sprintf("/%s/roles", engine)))
	if err != nil {
		return []string{}, err
	}

	result := getSecretAsStrings(data)
	return result, nil
}

func getRoleDetails(client *api.Client, engine string, roleName string) (*AwsRole, error) {
	data, err := client.Logical().Read(replaceDoubleSlash(fmt.Sprintf("/%s/roles/%s", engine, roleName)))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	role := &AwsRole{Path: engine, Role: roleName}
	role.CredentialType = data.Data["credential_type"].(string)
	role.DefaultStsTtl, _ = data.Data["default_sts_ttl"].(json.Number).Int64()
	role.MaxStsTtl, _ = data.Data["max_sts_ttl"].(json.Number).Int64()
	role.PolicyDocument = data.Data["policy_document"].(string)
	role.UserPath = data.Data["user_path"].(string)
	role.PermissionsBoundaryArn = data.Data["permissions_boundary_arn"].(string)
	role.RoleArns = getValues(data.Data, "role_arns")
	role.PolicyArns = getValues(data.Data, "policy_arns")
	role.IamGroups = getValues(data.Data, "iam_groups")

	return role, nil
}
