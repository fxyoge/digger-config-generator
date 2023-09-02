remote_state {
  backend = "gcs"
  generate = {
    path      = "_generated_backend.tf"
    if_exists = "overwrite_terragrunt"
  }
  config = {
    project     = "myproject"
    bucket      = "bucket"
    location    = "us"
    prefix      = "tfstate/${path_relative_to_include()}"
    credentials = get_env("GOOGLE_APPLICATION_CREDENTIALS")
  }
}
