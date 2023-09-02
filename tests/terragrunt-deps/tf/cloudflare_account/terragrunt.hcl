include "root" {
  path = find_in_parent_folders()
}

inputs = {
  cf_api_token = get_env("CF_API_TOKEN")
}
