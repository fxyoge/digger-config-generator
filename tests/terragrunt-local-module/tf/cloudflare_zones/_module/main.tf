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

variable "account_id" {
  type = string
}

variable "zone" {
  type = string
}

resource "cloudflare_zone" "main" {
  account_id = var.account_id
  zone       = var.zone
}

output "zone_id" {
  value = cloudflare_zone.main.id
}
