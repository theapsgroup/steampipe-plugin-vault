# Table: vault_aws_role

AWS Roles contained within Vault Mountpoints.

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