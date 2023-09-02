include "root" {
  path = find_in_parent_folders()
}

terraform {
  source = "..//_module"
}

dependency "account" {
  config_path = "${get_terragrunt_dir()}/../../cloudflare_account"
}

inputs = {
  cf_api_token = get_env("CF_API_TOKEN")

  account_id = dependency.account.outputs.account_id
  zone       = "example.com"
}
