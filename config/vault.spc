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