# Manifold Terraform Provider

[Homepage](https://manifold.co) |
[Twitter](https://twitter.com/manifoldco) |
[Code of Conduct](./.github/CODE_OF_CONDUCT.md) |
[Contribution Guidelines](./.github/CONTRIBUTING.md)

[![GitHub release](https://img.shields.io/github/tag/manifoldco/terraform-provider-manifold.svg?label=latest)](https://github.com/manifoldco/terraform-provider-manifold/releases)
[![Build Status](https://travis-ci.org/manifoldco/terraform-provider-manifold.svg?branch=master)](https://travis-ci.org/manifoldco/terraform-provider-manifold)
[![Go Report Card](https://goreportcard.com/badge/github.com/manifoldco/terraform-provider-manifold)](https://goreportcard.com/report/github.com/manifoldco/terraform-provider-manifold)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](./LICENSE.md)

![Terraform Manifold](./banner.png)

Manifold gives you a single account to purchase and manage cloud services from multiple providers, giving you managed logging, email, MySQL, Postgres, Memcache, Redis, and more. Manifold also lets you register configurations for your services external to Manifold's marketplace, giving you a single location to organize and manage the credentials for your environments.

This is a Terraform Provider to help you read the data of your provisioned
resources on the [Manifold.co](https://manifold.co) platform.

Currently, we only support data sources, not resources.

## Configuration

To use the Manifold Provider, you'll need an API Key. You can either provide
this in the provider configuration with the `api_token` field, or use an ENV,
`MANIFOLD_API_TOKEN`.

If you want to specify the team you want to use, you can do this by either
setting the field `team`, or by using the ENV `MANIFOLD_TEAM`.

## Getting an API Token

To retrieve an API token, use [our CLI tool](http://github.com/manifoldco/manifold-cli) and run the following:

```
$ manifold tokens create
```

## Installation

Bare zip archives per release version are available on [https://releases.manifold.co](https://releases.manifold.co/terraform-provider-manifold/).

Terraform currently doesn't allow custom providers to be fetched automatically,
so to use this plugin, you'll have to put the compiled binary in your terraform
plugin folder, which is located at `$HOME/.terraform.d/plugins/`.

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
