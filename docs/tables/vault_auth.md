# Table: vault_auth

Vault Authentication Methods currently configured.

## Columns

| Column | Description |
| - | - |
| path | The path (mount point) of the authentication method - for example `oidc/` |
| type | The type of authentication method - for example `oidc` |
| description | Description associated to the authentication method |
| accessor | The accessor used by authentication method |
| local | Determines if the Authentication Method is local only, local are not replicated across clusters |
| seal_wrap | Is the Authentication Method using [seal wrap](https://www.vaultproject.io/docs/enterprise/sealwrap) |
| external_entropy_access | Does Authentication Method have access to Vaults external entropy source |
| default_ttl | Default Lease TTL of Authentication Method (if set) |
| max_ttl | Maximum Lease TTL of Authentication Method (if set) |
| request_headers | Allowed Pass-Through Request Headers |

## Examples

### List all authentication methods

```sql
select
  *
from
  vault_auth;
```

### List authentication methods of a specifc type (oidc in example)

```sql
select
  *
from
  vault_auth
where
  type = 'oidc';
```