---
organization: The APS Group
category: ["security"]
icon_url: "/images/plugins/theapsgroup/vault.svg"
brand_color: "#003A75"
display_name: "Hashicorp Vault"
short_name: "vault"
description: "Steampipe plugin for querying available secret keys (not values), etc from Hashicorp Vault."
og_description: Query Hashicorp Vault with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/theapsgroup/vault-social-graphic.png"
---

# Hashicorp Vault + Steampipe

[Vault](https://www.vaultproject.io/) is an industry-leading Secrets Management & Data Protection solution from [Hashicorp](https://www.hashicorp.com/).

[Steampipe](https://steampipe.io/) is an open source CLI for querying cloud APIs using SQL from [Turbot](https://turbot.com/)

For example:
```sql
select
  path,
  type,
  description
from
  vault_engine;
```

## Documentation

- **[Table definitions & examples →](https://hub.steampipe.io/plugins/theapsgroup/vault/tables)**

## Get started

### Install

Download and install the latest Vault plugin:

```shell
steampipe plugin install theapsgroup/vault
```

### Configuration

Installing the latest vault plugin will create a config file (`~/.steampipe/config/vault.spc`) with a single connection named `vault`:

> Note: If using `token` based authentication, it is preferred that you use Environment Variables as the token should be rotated frequently.

```hcl
connection "vault" {
  plugin = "theapsgroup/vault"

  # The address of your Vault (ignore if VAULT_ADDR env var is set).
  # address = "https://your-vault-domain/"

  # Vault auth type to use, valid options are token and aws
  # auth_type = "token"

  # API Token for Vault (ignore if VAULT_TOKEN env var is set).
  # token = "YOUR_VAULT_TOKEN"

  # For aws authentication
  # auth_type = "aws"
  # The vault role to authenticate as
  # aws_role = "steampipe-role"
  # The name of the aws auth backend to use for authentication
  # aws_provider = "awspath"
}
```

- `token` - [Vault Token](https://developer.hashicorp.com/vault/api-docs/auth/token) for your Vault. This can also be set via the `VAULT_TOKEN` environment variable.
- `address` - The url of your Vault server (e.g. `https://vault.mycorp.com/`). This can also be via the `VAULT_ADDR` environment variable.
- `auth_type` - Should be either `token` to use token based authentication or `aws` to use AWS authentication via the `aws_role` & `aws_provider` properties.
- `aws_role` - The Vault aws role to authenticate as.
- `aws_provider` - The name of the AWS authentication backend to use for authentication.

#### Authentication

Vault supports multiple authentication backends, currently token and AWS IAM are supported.
**Note that in line with the Vault cli behavior, if a vault token is supplied, that will be used instead of your configured authentication method.**

##### Token Example

```hcl
connection "vault" {
  plugin    = "theapsgroup/vault"
  address   = "https://vault.mycorp.com/"
  auth_type = "token"
  token     = "sometoken"
}
```

##### AWS Example

```hcl
connection "vault" {
  plugin    = "theapsgroup/vault"
  address   = "https://vault.mycorp.com/"
  auth_type = "aws"
  aws_role = "steampipe-test-role"
  aws_provider = "aws"
}
```

The vault plugin will resolve the AWS credentials in the normal AWS SDK Credentials chain order.

## Get involved

- Open source: https://github.com/theapsgroup/steampipe-plugin-vault
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)