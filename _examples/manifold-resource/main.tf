provider "manifold" {
  // uses the default and loads the API key from the ENV `MANIFOLD_API_TOKEN`.
}

// This example loads a resource and filters out specific credentials. This way
// you can select only the ones you need, or set up an alias.
// First. we'll just select the `TOKEN_ID` credential, as is.
// In the second credential block, we'll use the `secret` alias, which we can
// use later on to reference our credential.
// In the third example, we'll try and get a non existing key and give it a
// default value, which will be used to populate the credentials map.
data "manifold_resource" "example1" {
  project  = "manifold-terraform"
  resource = "custom-resource1-1"

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

output "TOKEN_ID" {
  value = "${data.manifold_resource.example1.credentials.TOKEN_ID}"
}

output "TOKEN_SECRET" {
  value = "${data.manifold_resource.example1.credentials.secret}"
}

output "DEFAULT" {
  value = "${data.manifold_resource.example1.credentials.default-example}"
}

// In this example we'll select all the credentials for our resource without
// filtering any out.
data "manifold_resource" "example2" {
  project  = "manifold-terraform"
  resource = "custom-resource1-1"
}

output "TOKEN_SECRET_2" {
  value = "${data.manifold_resource.example2.credentials.TOKEN_SECRET}"
}
