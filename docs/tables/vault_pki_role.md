# Table: vault_pki_role

For querying Roles in the pki [engines](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_engines.md)

## Examples

### Get all roles in PKI mounts

```sql
select
  *
from
  vault_pki_role;
```

### Obtain roles which have Code Signing capabilities

```sql
select
  path,
  name,
from
  vault_pki_role
where
  code_signing_flag = 1;
```