# Hashicorp Vault Plugin for Steampipe

## Query HashiCorp Vault with SQL

Use SQL to query Vault. Example:

```sql
select * from vault_engine
```

## Get Started

### Installation

```shell
steampipe plugin install theapsgroup/vault
```

Or if you prefer, you can clone this repository and build/install from source directly.

```shell
go build -o steampipe-plugin-vault.plugin

mv steampipe-plugin-vault.plugin ~/.steampipe/plugins/hub.steampipe.io/plugins/theapsgroup/vault@latest/steampipe-plugin-vault.plugin

cp config/vault.spc ~/.steampipe/config/vault.spc
```

### Configuration

Configuration is preferably done by ensuring you have the default Vault Environment Variables set:

- `VAULT_ADDR` for the address of your Vault Server
- `VAULT_TOKEN` for the API token used to access Vault

However, these can also be set in the configuration file:

`vi ~/.steampipe/config/vault.spc` 

### Testing Installation

```shell
steampipe query "select * from vault_engine"
```

## Documentation

Further documentation can he [found here](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/index.md)