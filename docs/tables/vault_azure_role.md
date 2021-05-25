# Table: vault_azure_role

Azure roles contained within Vault.

## Examples

### List All Azure Roles

```sql
select
  *
from
  vault_azure_role;
```

### Roles matching a specific pattern - in example containing `deploy`

```sql
select
  *
from
  vault_azure_role
where
  role like '%deploy%'
```