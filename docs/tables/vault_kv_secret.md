# Table: vault_kv_secret

For working with paths for secrets in the kv [engines](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_engines.md)

> Note: This does not expose the contents of the secrets by design.

## Columns

| Column | Description |
| - | - |
| key | The key of the secret within the kv engine |
| path | The path at which an engine is mounted - for example `apples/` |
| created_time | The date and time the secret was created - for example `2020-11-06 11:12:54` |
| deletion_time | The date and time the secret was destroyed - for example `2021-01-06 12:14:56` |
| destroyed | Whether the secret was destroyed - for example `false` |
| version | The current version of the secret - for example `4` |

## Examples

### Get all secret keys from all kv engines

```sql
select
  key,
  path
from
  vault_kv_secret;
```

### Get all secret keys from a specific mounted kv engine (`abc/` in this example)

```sql
select
  key
from
  vault_kv_secret
where
  path = 'abc/';
```

### Search for secret paths based on a fragment/keyword

```sql
select
  *
from
  vault_kv_secret
where
  key like '%myapp%';
```