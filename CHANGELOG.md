## v0.0.5 [WIP]

_Enhancements_
- Upgraded to golang version 1.17
- Upgraded steampipe sdk version to v1.8.2
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