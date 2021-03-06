# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## Unreleased

## [0.1.0] - 2020-04-01

### Added

- Support the terraform 0.12 API

## [0.0.3] - 2018-08-07

- Add `manifold_credential` data source for fetching single credential value
- Add `manifold_credential` resource for managing custom credential values
- Update `go-manifold` version

## [0.0.2] - 2017-11-20

### Added

- Added a resource type which allows users to create API tokens.

### Fixed

- Fix a bug upstream in go-manifold where the UA format was malformed.

## [0.0.1] - 2017-10-30

### Added

- The start of the Manifold Terraform Provider
- Added a data source to fetch credentials for a specific resource
- Added a data source to fetch credentials for a specific project
