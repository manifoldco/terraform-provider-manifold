# Manifold Terraform Provider

This is a Terraform Provider to help you read the data of your provisioned
resources on the [Manifold.co](https://manifold.co) platform.

Currently, we only support data sources, not resources.

## Configuration

To use the Manifold Provider, you'll need an API Key. You can either provide
this in the provider configuration with the `api_key` field, or use an ENV,
`MANIFOLD_API_KEY`.

If you want to specify the team you want to use, you can do this by either
setting the field `team`, or by using the ENV `MANIFOLD_TEAM`.

## Getting an API Token

To retrieve an API token, use [our CLI tool](http://github.com/manifoldco/manifold-cli) and run the following:

```
$ manifold tokens create
```

## Examples

We've included a set of examples to get you started and to understand what you
can do with our provider.

### Setup

The Manifold setup for our examples is as follows:

- *Project:* manifold-terraform
    - *Resource:* custom-resource1
        - *Credential*: TOKEN_ID
        - *Credential*: TOKEN_SECRET
    - *Resource:* custom-resource2
        - *Credential*: USERNAME
        - *Credential*: PASSWORD

### Examples

- [Load data for an entire project](_examples/manifold-project/README.md)
- [Load data for a specific resource](_examples/manifold-resource/README.md)
