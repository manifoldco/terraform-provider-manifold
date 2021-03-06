provider "manifold" {
  // uses the default and loads the API key from the ENV `MANIFOLD_API_TOKEN`.
}

// This example shows how to get all credentials for a project in one go. The
// setup contains 2 custom resources, `custom-resource1-1` and
// `custom-resource2-1`.
// `custom_resource1` has 2 credentials, `TOKEN_ID` and `TOKEN_SECRET`.
// `custom_resource2` has 2 credentials, `USERNAME` and `PASSWORD`.
data "manifold_project" "no_resource_selected" {
  project = "manifold-terraform"
}

output "nr_token_id" {
  value = "${data.manifold_project.no_resource_selected.credentials.TOKEN_ID}"
}

output "nr_token_secret" {
  value = "${data.manifold_project.no_resource_selected.credentials.TOKEN_SECRET}"
}

output "nr_username" {
  value = "${data.manifold_project.no_resource_selected.credentials.USERNAME}"
}

output "nr_password" {
  value = "${data.manifold_project.no_resource_selected.credentials.USERNAME}"
}

// This example shows how to get all credentials for a specific set of
// resources, in this case only 1.
// The difference with the example above is that the `TOKEN_ID` and
// `TOKEN_SECRET` are not available anymore.
data "manifold_project" "resource_selected" {
  project = "manifold-terraform"

  resource {
    resource = "custom-resource2-1"
  }
}

output "sr_username" {
  value = "${data.manifold_project.resource_selected.credentials.USERNAME}"
}

output "sr_password" {
  value = "${data.manifold_project.resource_selected.credentials.PASSWORD}"
}

// This example shows how to filter for specific credentials across different
// resources. The setup is the same as with selecting credentials through the
// `manifold_resource` data source.
// Here we'll fetch the `USERNAME` value from `custom-resource2-1` and alias it
// with `my-alias`. We also get a credential that hasn't been set,
// `NON_EXISTING`and give it a default value `my-default-secret`. Lastly, we
// just select the `TOKEN_ID` from `custom-resource1-1`.
data "manifold_project" "credential_selected" {
  project = "manifold-terraform"

  resource {
    resource = "custom-resource2-1"

    credential {
      name = "my_alias"
      key  = "USERNAME"
    }

    credential {
      key     = "NON_EXISTING"
      default = "my-default-secret"
    }
  }

  resource {
    resource = "custom-resource1-1"

    credential {
      key = "TOKEN_ID"
    }
  }
}

output "cs_username" {
  value = "${data.manifold_project.credential_selected.credentials.my_alias}"
}

output "cs_non_existing" {
  value = "${data.manifold_project.credential_selected.credentials.NON_EXISTING}"
}

output "cs_token_id" {
  value = "${data.manifold_project.credential_selected.credentials.TOKEN_ID}"
}
