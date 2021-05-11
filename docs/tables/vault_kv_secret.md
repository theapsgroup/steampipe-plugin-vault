# Table: vault_kv_secret

For working with paths for secrets in the kv [engines](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_engines.md)

> Note: This does not expose the contents of the secrets by design.

## Columns

| Column | Description |
| - | - |
| path | The path of the secret within the kv engine |
| mountpoint | The path at which an engine is mounted - for example `apples/` |

## Examples

### Get all secret paths from all kv engines

```sql
select
  path,
  mountpoint
from
  vault_kv_secret;
```

### Get all secret paths from a specific mounted kv engine (`abc/` in this example)

```sql
select
  path
from
  vault_kv_secret
where
  mountpoint = 'abc/';
```

### Search for secret paths based on a fragment/keyword

```sql
select
  *
from
  vault_kv_secret
where
  path like '%myapp%';
```