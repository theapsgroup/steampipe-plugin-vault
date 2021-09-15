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
```

Alternatively, `make install` will do the same steps as above.

Copy the basic configuration:

```shell
cp config/vault.spc ~/.steampipe/config/vault.spc
```

### Configuration

Currently, two authentication methods are supported. Token authentication, and AWS IAM role authentication.

Configuration can be done by following the Vault cli environment variables:

- `VAULT_ADDR` for the address of your Vault Server
- `VAULT_TOKEN` for the API token used to access Vault

However, these can also be set in the configuration file:

`vi ~/.steampipe/config/vault.spc` 

#### AWS Authentication

In addition to the `VAULT_TOKEN` environment variable, the plugin can also authenticate using AWS IAM role credentials. This works via the AWS authentication backend in Vault ([more info](https://www.vaultproject.io/docs/auth/aws)). Configuration is done in the `vault.spc` file.

For example:

```
connection "vault" {
  plugin = "theapsgroup/vault"

  address = "https://your-vault-domain/"

  # One of "aws" or "token"
  auth_type = "aws"
  # The name of the aws auth backend role in vault
  aws_role = "steampipe-role"
  # The mount path of the aws auth backend (default would be "aws")
  aws_provider = "awspath"
}
```

**Note: In line with the vault cli, if you provide `VAULT_TOKEN` that will be used INSTEAD of the aws authentication method**

### Testing Installation

```shell
steampipe query "select * from vault_engine"
```

## Documentation

Further documentation can he [found here](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/index.md)