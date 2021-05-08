# Table: vault_aws_roles

AWS Roles contained within Vault Mountpoints.

## Columns

| Column | Description |
| - | - |
| mountpoint | The path at which an engine is mounted - for example `aws/` |
| role | The aws role name - for example `prod-deploy` |

## Examples

### List all AWS Roles

```sql
select
  *
from
  vault_aws_roles
```

### Roles matching a specific pattern - in example containing `deploy`

```sql
select
  *
from
  vault_aws_roles
where
  role like '%deploy%'
```