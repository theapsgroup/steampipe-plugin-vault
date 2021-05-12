# Table: vault_sys_health

Allows for displaying system health information.

> Note: This should only every return a single row of data, this is expected.

## Columns

| Column | Description |
| - | - |
| initialized | Is Initialized (expected `true`) |
| sealed | Is Sealed (expected `false`) |
| standby | Is Standby (bool) |
| performance_standby | Is Performance Standby (bool) |
| replication_performance_mode | Replication Performance Mode (example: `disabled`) |
| replication_dr_mode | Replication Disaster Recovery Mode (example: `disabled`) |
| server_time_utc | Current Server Time (UTC) |
| version | Vault Server Version |
| cluster_name | Vault Cluster Name |
| cluster_id | Vault Cluster Identifier |

## Examples

### Get system health

```sql
select
  *
from
  vault_sys_health
```