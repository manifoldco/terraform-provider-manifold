provider "manifold" {
    // uses the default and loads the API key from the ENV `MANIFOLD_API_KEY`.
}

provider "aws" {
    // this uses the default aws configuration
}

data "manifold_resource" "my-service" {
    project = "my-unique-project"
    resource = "my-unique-resource"

    // Specify specific credentials to load for this resource.
    credential {
        name    = "username" // the reference key for later on
        key     = "USERNAME" // the key/value key from the values field
        default = "manifold" // optional, default value
    }

    credential {
        name    = "password" // the reference key for later on
        key     = "PASSWORD" // the key/value key from the values field
    }
}

resource "aws_ec2" "my-service" {
    // ec2 setup

    env_vars = {
        "SERVICE_USERNAME=${data.manifold_resource.my-resource.credentials.username}",
        "SERVICE_PASSWORD=${data.manifold_resource.my-resource.credentials.password}",
        // ...
    }
}

data "manifold_resource" "my-second-service" {
    project = "my-unique-project"
    resource = "my-unique-second-resource"

    // Don't specify credentials to load, load all of them and use their key as
    // reference name.
}

resource "aws_ec2" "my-service" {
    // ec2 setup

    env_vars = {
        "SERVICE_USERNAME=${data.manifold_resource.my-resource.credentials.USERNAME}",
        "SERVICE_PASSWORD=${data.manifold_resource.my-resource.credentials.PASSWORD}",
        // ...
    }
}
