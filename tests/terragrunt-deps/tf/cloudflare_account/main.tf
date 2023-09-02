terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "4.13.0"
    }
  }
}

provider "cloudflare" {
  api_token = var.cf_api_token
}

variable "cf_api_token" {
  type      = string
  sensitive = true
}

data "cloudflare_accounts" "main" {
  name = "my_account"
}

output "account_id" {
  value = data.cloudflare_accounts.main.accounts[0].id
}
