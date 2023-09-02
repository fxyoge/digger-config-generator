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

data "cloudflare_ip_ranges" "cloudflare" {}

output "cidr_blocks" {
  value = data.cloudflare_ip_ranges.cloudflare.cidr_blocks
}
