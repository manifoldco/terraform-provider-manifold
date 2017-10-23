provider "manifold" {
  // uses the default and loads the API key from the ENV `MANIFOLD_API_KEY`.
}

data "manifold_resource" "example1" {
  project  = "manifold-terraform"
  resource = "custom-resource1"

  credential {
    key = "TOKEN_ID"
  }

  credential {
    name = "secret"
    key  = "TOKEN_SECRET"
  }

  credential {
    name    = "default-example"
    key     = "DEFAULT_EXAMPLE"
    default = "default-value"
  }
}

data "manifold_resource" "example2" {
  project  = "manifold-terraform"
  resource = "custom-resource1"
}

output "TOKEN_ID" {
  value = "${data.manifold_resource.example1.credentials.TOKEN_ID}"
}

output "TOKEN_SECRET" {
  value = "${data.manifold_resource.example1.credentials.secret}"
}

output "DEFAULT" {
  value = "${data.manifold_resource.example1.credentials.default-example}"
}

output "TOKEN_SECRET_2" {
  value = "${data.manifold_resource.example2.credentials.TOKEN_SECRET}"
}
