# Table: vault_aws_role

AWS Roles contained within Vault Mountpoints.

## Columns

| Column | Description |
| - | - |
| path | The path at which an engine is mounted - for example `aws/` |
| role | The AWS role name - for example `prod-deploy` |
| credential_type | The type of credential assumed_role, iam, etc |
| default_sts_ttl | Default STS TTL |
| max_sts_ttl | Maximum STS TTL |
| policy_document | JSON AWS policy document associated with the role |
| user_path | Path of the user |
| permissions_boundary_arn | ARN of the permissions boundary |
| role_arns | ARNs associated with the role |
| policy_arns | ARNs associated with the policies of the role |
| iam_groups | IAM groups associated with the role |

## Examples

### List all AWS Roles

```sql
select
  *
from
  vault_aws_role
```

### Roles matching a specific pattern - in example containing `deploy`

```sql
select
  *
from
  vault_aws_role
where
  role like '%deploy%'
```