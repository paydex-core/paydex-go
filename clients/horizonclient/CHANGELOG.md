# Changelog

All notable changes to this project will be documented in this
file.  This project adheres to [Semantic Versioning](http://semver.org/).

## Unreleased

- Dropped support for Go 1.10, 1.11.

### Added

- `Client.Root()` method for querying the root endpoint of a horizon server.

### Changes

- `Client.Fund()` now returns `TransactionSuccess` instead of a http response pointer.

- Querying the effects endpoint now supports returning the concrete effect type for each effect. This is also supported in streaming mode. See the [docs](https://godoc.org/github.com/paydex-core/paydex-go/clients/horizonclient#Client.Effects) for examples.

## [v1.0.0](https://github.com/paydex-core/paydex-go/releases/tag/horizonclient-v0.1.1) - 2020-03-26

 * Initial release
