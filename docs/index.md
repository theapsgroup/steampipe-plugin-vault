---
org: The APS Group
category: ["security"]
icon_url: "/images/plugins/theapsgroup/vault.svg"
brand_color: "#003A75"
display_name: "Hashicorp Vault"
short_name: "vault"
description: "Steampipe plugin for querying available secret keys (not values), etc from Hashicorp Vault."
social_about: Query Hashicorp Vault with SQL! Open source CLI. No DB required.
social_preview: "/images/plugins/theapsgroup/vault-social-graphic.png"
---

# Hashicorp Vault + Turbot Steampipe

[Vault](https://www.vaultproject.io/) is an industry-leading Secrets Management & Data Protection solution from [Hashicorp](https://www.hashicorp.com/).

[Steampipe](https://steampipe.io/) is an open source CLI for querying cloud APIs using SQL from [Turbot](https://turbot.com/)

## Getting Started

### Installation

Download and install the latest Vault plugin:

```shell
steampipe plugin install theapsgroup/vault
```

### Prerequisites

- Vault Server
- Vault API Token

### Configuration

The preferred option is to use Environment Variables for configuration as the Vault Token should be rotated frequently, however you can configure in the `~./steampipe/config/vault.spc` (this will take precedence).

Environment Variables (default from Hashicorp Vault):

- `VAULT_ADDR` for the server address (ex: `https://vault.mycorp.com/`)
- `VAULT_TOKEN` for the API token (ex: `s.f7Ea3C3ojOYE0GRLzmhSGNkE`)

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
steampipe query "select * from vault_engine"
```

## Tables

The following tables are available for querying, follow the links for more information.

- [vault_sys_health](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_sys_health.md)
- [vault_engine](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_engine.md)
- [vault_auth](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_auth.md)
- [vault_kv_secret](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_kv_secret.md)
- [vault_aws_role](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_aws_role.md)
- [vault_azure_role](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_azure_role.md)
- [vault_pki_cert](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_pki_cert.md)
- [vault_pki_role](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_pki_role.md)
