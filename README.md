![image](https://hub.steampipe.io/images/plugins/theapsgroup/vault-social-graphic.png)

# Hashicorp Vault Plugin for Steampipe

## Query HashiCorp Vault with SQL

* **[Get started →](https://hub.steampipe.io/plugins/theapsgroup/vault)**
* Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/theapsgroup/vault/tables)
* Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
* Get involved: [Issues](https://github.com/theapsgroup/steampipe-plugin-vault/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io/downloads):

```shell
steampipe plugin install theapsgroup/vault
```

[Configure the plugin](https://hub.steampipe.io/plugins/theapsgroup/vault#configuration) using the configuration file:

```shell
vi ~/.steampipe/vault.spc
```

Or environment variables:

```shell
export VAULT_ADDR=https://vault.mycorp.com/
export VAULT_TOKEN=s.f7Ea3C3ojOYE0GRLzmhSGNkE
```

Start Steampipe:

```shell
steampipe query
```

Run a query:

```sql
select
  path,
  type,
  description
from
  vault_engine;
```

[//]: # (/*)

[//]: # (Or if you prefer, you can clone this repository and build/install from source directly.)

[//]: # ()
[//]: # (```shell)

[//]: # (go build -o steampipe-plugin-vault.plugin)

[//]: # ()
[//]: # (mv steampipe-plugin-vault.plugin ~/.steampipe/plugins/hub.steampipe.io/plugins/theapsgroup/vault@latest/steampipe-plugin-vault.plugin)

[//]: # (```)

[//]: # ()
[//]: # (Alternatively, `make install` will do the same steps as above.)

[//]: # ()
[//]: # (Copy the basic configuration:)

[//]: # ()
[//]: # (```shell)

[//]: # (cp config/vault.spc ~/.steampipe/config/vault.spc)

[//]: # (```)

[//]: # ()
[//]: # (### Configuration)

[//]: # ()
[//]: # (Currently, two authentication methods are supported. Token authentication, and AWS IAM role authentication.)

[//]: # ()
[//]: # (Configuration can be done by following the Vault cli environment variables:)

[//]: # ()
[//]: # (- `VAULT_ADDR` for the address of your Vault Server)

[//]: # (- `VAULT_TOKEN` for the API token used to access Vault)

[//]: # ()
[//]: # (However, these can also be set in the configuration file:)

[//]: # ()
[//]: # (`vi ~/.steampipe/config/vault.spc` )

[//]: # ()
[//]: # (#### AWS Authentication)

[//]: # ()
[//]: # (In addition to the `VAULT_TOKEN` environment variable, the plugin can also authenticate using AWS IAM role credentials. This works via the AWS authentication backend in Vault &#40;[more info]&#40;https://www.vaultproject.io/docs/auth/aws&#41;&#41;. Configuration is done in the `vault.spc` file.)

[//]: # ()
[//]: # (For example:)

[//]: # ()
[//]: # (```)

[//]: # (connection "vault" {)

[//]: # (  plugin = "theapsgroup/vault")

[//]: # ()
[//]: # (  address = "https://your-vault-domain/")

[//]: # ()
[//]: # (  # One of "aws" or "token")

[//]: # (  auth_type = "aws")

[//]: # (  # The name of the aws auth backend role in vault)

[//]: # (  aws_role = "steampipe-role")

[//]: # (  # The mount path of the aws auth backend &#40;default would be "aws"&#41;)

[//]: # (  aws_provider = "awspath")

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (**Note: In line with the vault cli, if you provide `VAULT_TOKEN` that will be used INSTEAD of the aws authentication method**)

[//]: # (*/)
## Developing

Prerequisites:

* [Steampipe](https://steampipe.io/downloads)
* [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/theapsgroup/steampipe-plugin-vault.git
cd steampipe-plugin-vault
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/vault.spc
```

Try it!

```shell
steampipe query
> .inspect vault
```

Further reading:

* [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
* [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-github/blob/main/LICENSE).
