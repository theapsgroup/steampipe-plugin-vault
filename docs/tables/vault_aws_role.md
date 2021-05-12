# Table: vault_aws_role

AWS Roles contained within Vault Mountpoints.

## Columns

| Column | Description |
| - | - |
| path | The path at which an engine is mounted - for example `aws/` |
| role | The aws role name - for example `prod-deploy` |

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