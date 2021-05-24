# Table: vault_auth

Vault Authentication Methods currently configured.

## Examples

### List all authentication methods

```sql
select
  *
from
  vault_auth;
```

### List authentication methods of a specific type (oidc in example)

```sql
select
  *
from
  vault_auth
where
  type = 'oidc';
```