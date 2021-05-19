# Table: vault_azure_config

Configuration settings for Azure mountpoints in Vault.

## Column

| Column | Description |
| - | - |
| path | The path at which an engine is mounted - for example `azure/` |
| subscription_id | The Azure subscription identifier |
| tenant_id | The Azure tenant identifier |
| client_id | The Azure client identifier |
| environment | The Azure environment description |

## Examples

### List all Azure Configurations

```sql
select
  *
from
  vault_azure_config;
```

