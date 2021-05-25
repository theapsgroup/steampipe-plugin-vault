# Table: vault_sys_health

Allows for displaying system health information.

> Note: This should only every return a single row of data, this is expected.

## Examples

### Get system health

```sql
select
  *
from
  vault_sys_health
```