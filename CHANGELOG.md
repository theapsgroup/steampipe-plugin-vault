## v0.3.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29).
- Recompiled plugin with Go version `1.21`.

## v0.2.1 [2023-05-05]

_Enhancements_

- Updated: Recompiled with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05)

## v0.2.0 [2023-01-10]

_What's new?_

- Added `plugin_version`, `deprecation_status` & `options` columns to the `vault_auth` table.
- Added `plugin_version`, `deprecation_status` & `options` columns to the `vault_engine` table.
- Added `last_wal` column to the `vault_health` table.

_Enhancements_

- Updated: Recompiled with [steampipe-plugin-sdk v5.0.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v501-2022-11-30)
- Updated: Recompiled with [AWS SDK v1.44.176](https://github.com/aws/aws-sdk-go/releases/tag/v1.44.176)
- Updated: Recompiled with [Vault API v1.8.0](https://github.com/hashicorp/vault/tree/main/api)

## v0.1.1 [2022-10-08]

_Enhancements_
- Upgraded to golang version 1.19
- Upgraded Steampipe sdk version to v4.1.7

## v0.1.0 [2022-05-05]

_Enhancements_
- Upgraded to golang version 1.18
- Upgraded Steampipe sdk version to v3.1.0

## v0.0.5 [2021-11-29]

_Enhancements_
- Upgraded to golang version 1.17
- Upgraded Steampipe sdk version to v1.8.2
- Upgraded vault version to v1.3.0
- Upgraded aws sdk version to v1.42.11

_Bug fixes_
- Fixed an issue where not setting an `auth_type` would cause an issue [#19](https://github.com/theapsgroup/steampipe-plugin-vault/issues/19)

## v0.0.4 [2021-09-16]

_What's new?_
- Adds a makefile for easier installation while developing

_Bug fixes_
- Fixes an issue where the AWS auth type would incorrectly ignore address configuration

## v0.0.3 [2021-08-17]

_What's new?_

- Allow vault plugin to authenticate using the Vault "AWS" Backend

## v0.0.2 [2021-08-17]

_What's new?_

- Update frontmatter
- Actual first release to steampipe hub

## v0.0.1 [2021-05-25]

_What's new?_

- Adds goreleaser
