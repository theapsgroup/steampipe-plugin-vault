# Table: vault_pki_cert

For querying PKI Certificates in the pki [engines](https://github.com/theapsgroup/steampipe-plugin-vault/blob/main/docs/tables/vault_engines.md)

## Examples

### Get all certificates in all pki mounts

```sql
select
  *
from
  vault_pki_cert;
```

### Get certificates from a specific engine mount (example is `pki/`)

```sql
select
  *
from
  vault_pki_cert
where
  path = 'pki/';
```

### Get renewable certificates from the `pki/` mount

```sql
select
  *
from
  vault_pki_cert
where
  path = 'pki/'
and
  renewable = 1;
```