![image](https://hub.steampipe.io/images/plugins/theapsgroup/vault-social-graphic.png)

# Hashicorp Vault Plugin for Steampipe

Use SQL to query engines, kv secrets, roles & more from your Hashicorp Vault.

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
