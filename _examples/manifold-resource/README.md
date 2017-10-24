# Manifold Resource

This example will show you how to target specific resources using the Manifold
Terraform Provider.

## Used resources

- `manifold.data.manifold_resource`

## Prerequisites

This example assumes you have an account with [Manifold.co](https://www.manifold.co/) and have access to
AWS.

It will also assume that you have exported the Manifold API Key to the
environment variable, `MANIFOLD_API_TOKEN`.

## Loading credentials

There are 2 examples given to illustrate the posibilities of configuring the
data source.

The first example loads specific credentials for your resource, giving you
control of what will become available. You can name these credentials as you
like and use this name later on as a reference. Within this example, we've
illustrated several ways of defining the credentials. The minimal requirement
is providing a `key`, which will be used to filter out only that key from your
resource credentials.

The second example does not provide a credential filter, meaning we'll load all
available credentials and use the stored KEY as a reference name which you can
use later on.
