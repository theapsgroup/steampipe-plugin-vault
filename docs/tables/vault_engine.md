# Table: vault_engine

Vault Secrets Engines currently mounted.

## Columns

| Column | Description |
| - | - |
| path | The path at which an engine is mounted - for example `apples/` |
| type | The type of engine used by the mountpoint, such as `kv`, `aws`, etc |
| description | Description associated with the mountpoint of the engine |
| accessor | The accessor used by the engine |
| version | Version of the secrets engine (can be null for unversionable engine types) |
| local | Determines if the mountpoint is local only, local mountpoints are not replicated across clusters |
| seal_wrap | Determines if the engine at the mountpoint is using [seal wrap](https://www.vaultproject.io/docs/enterprise/sealwrap) |
| default_ttl | Default Lease TTL of Secrets Engine (if set) |
| max_ttl | Maximum Lease TTL of Secrets Engine (if set) |


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