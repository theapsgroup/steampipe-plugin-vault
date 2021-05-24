# Table: vault_kv_secret

For working with paths for secrets in the kv [engines](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_engines.md)

> Note: This does not expose the contents of the secrets by design.

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