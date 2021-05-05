# Hashicorp Vault Plugin for Steampipe

## Query HashiCorp Vault with SQL

Use SQL to query Vault. Example:

```sql
select * from vault_engines
```

## Get Started

Build & Installation from Source:

```shell
go build -o steampipe-plugin-vault.plugin

mv steampipe-plugin-vault.plugin ~/.steampipe/plugins/hub.steampipe.io/plugins/theapsgroup/vault@lateststeampipe-plugin-vault.plugin

cp config/vault.spc ~/.steampipe/config/vault.spc
```
## Documentation

Coming soon...