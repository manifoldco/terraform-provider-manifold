provider "manifold" {
    // uses the default and loads the API key from the ENV `MANIFOLD_API_KEY`.
}

provider "aws" {
    // this uses the default aws configuration
}

data "manifold_project" "marketplace" {
    project = "marketplace-unique-label"

    resource {
        name    = "identity"
        label   = "identity-unique-label"

        credential {
            name    = "username"
            key     = "USERNAME"
        }
    }

    resource {
        name    = "catalog"
        label   = "catalog-unique-label"
    }
}

resource "aws_ec2" "identity" {
    // ec2 setup

    env_vars = {
        "USERNAME=${data.manifold_project.marketplace.resources.identity.credentials.username}",
    }
}

resource "aws_ec2" "catalog" {
    // ec2 setup

    env_vars = {
        "DATABASE_URL=${data.manifold_project.marketplace.resources.catalog.credentials.DATABASE_URL}",
    }
}
