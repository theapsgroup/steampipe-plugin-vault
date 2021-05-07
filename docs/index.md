# Hashicorp Vault + Turbot Steampipe

[Vault](https://www.vaultproject.io/) is an industry-leading Secrets Management & Data Protection solution from [Hashicorp](https://www.hashicorp.com/).

[Steampipe](https://steampipe.io/) is an open source CLI tool for querying cloud APIs using SQL from [Turbot](https://turbot.com/)

## Getting Started

### Build & Installation

Currently, you will need to build and install this plugin manually via the following:

```shell
go build -o steampipe-plugin-vault.plugin

mv steampipe-plugin-vault.plugin ~/.steampipe/plugins/hub.steampipe.io/plugins/theapsgroup/vault@latest/steampipe-plugin-vault.plugin

cp config/vault.spc ~/.steampipe/config/vault.spc
```

### Prerequisites

- Vault Server
- Vault API Token

### Configuration

The preferred option is to use Environment Variables for configuration as the Vault Token should be rotated frequently, however you can configure in the `~./steampipe/config/vault.spc` (this will take precedence).

Environment Variables (default from Hashicorp Vault):

- `VAULT_ADDR` for the server address (ex: https://vault.mycorp.com/)
- `VAULT_TOKEN` for the API token (ex: s.f7Ea3C3ojOYE0GRLzmhSGNkE)

Configuration File:

```hcl
connection "vault" {
  plugin  = "theapsgroup/vault"
  address = "https://vault.mycorp.com/"
  token   = "s.f7Ea3C3ojOYE0GRLzmhSGNkE"
}
```
### Testing

A quick test can be performed from your terminal with:

```shell
steampipe query "select * from vault_engines"
```

## Tables

The following tables are available for querying, follow the links for more information.

- [vault_sys_health](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_sys_health.md)
- [vault_engines](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_engines.md)
- [vault_kv_secrets](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_kv_secrets.md)