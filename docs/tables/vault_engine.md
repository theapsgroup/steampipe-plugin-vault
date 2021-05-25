# Table: vault_engine

Vault Secrets Engines currently mounted.

## Examples

### List all mounted engines

```sql
select
  *
from
  vault_engine;
```

### Get the path of mounted engines which are of the Key Value (KV) type

```sql
select
  path
from
  vault_engine
where
  type = 'kv';
```

### Get a count of engines by type

```sql
select
  type,
  count(*)
from
  vault_engine
group by
  type;
```