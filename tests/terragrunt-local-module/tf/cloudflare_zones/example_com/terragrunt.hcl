include "root" {
  path = find_in_parent_folders()
}

terraform {
  source = "..//_module"
}

inputs = {
  cf_api_token = get_env("CF_API_TOKEN")

  account_id = "1234"
  zone       = "example.com"
}
