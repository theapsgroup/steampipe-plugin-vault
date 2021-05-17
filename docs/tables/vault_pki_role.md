# Table: vault_pki_role

For querying Roles in the pki [engines](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_engines.md)

## Columns

| Column | Description |
| - | - |
| path | The path at which an engine is mounted - for example `pki/` |
| name | The name of the role |
| allow_any_name | Is any name allowed on generated certificates |
| allow_ip_sans | Are IP based subject alternative names allowed |
| allow_localhost | Are localhost issued certificates allowed | 
| allow_sub_domains | Are certificates allowed for sub-domains |
| allowed_domains | Array of allowed domain names |
| allowed_uri_sans | Array of allowed uri based subject alternative names |
| allowed_other_sans | Array of other allowed subject alternative names | 
| client_flag | Can generate client based certificates |
| code_signing_flag | Can generate code-signing based certificates |
| key_bits | Length of key in bits |
| key_type | Type of key used - for example `RSA` |
| ttl | Default TTL |
| max_ttl | Maximum TTL |
| server_flag | Can generate server based certificates |

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