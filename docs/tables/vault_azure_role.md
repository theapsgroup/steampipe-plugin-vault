# Table: vault_azure_role

Azure roles contained within Vault.

## Columns

| Column | Description |
| - | - |
| path | The path at which an engine is mounted - for example `azure/` |
| role | The name of the role |

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