# Table: vault_engine

Vault Secrets Engines currently mounted.

## Columns

| Column | Description |
| - | - |
| path | The path at which an engine is mounted - for example `apples/` |
| type | The type of engine used by the mountpoint, such as `kv`, `aws`, etc |

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