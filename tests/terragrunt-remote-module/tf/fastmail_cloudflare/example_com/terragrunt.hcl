terraform {
  source = "tfr://MartiUK/fastmail/cloudflare?version=1.0.3"
}

generate "provider" {
  path      = "_generated_provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
    provider "cloudflare" {
      api_token = var.cf_api_token
    }
    
    variable "cf_api_token" {
      type      = string
      sensitive = true
    }
  EOF
}

inputs = {
  cf_api_token = get_env("CF_API_TOKEN")

  domain  = "example.com"
  zone_id = "1234"
}
