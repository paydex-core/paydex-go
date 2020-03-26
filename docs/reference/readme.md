---
title: Overview
---

The Go SDK is a set of packages for interacting with most aspects of the Paydex ecosystem. The primary component is the Horizon SDK, which provides convenient access to Horizon services. There are also packages for other Paydex services such as [TOML support](https://github.com/paydex/paydex-protocol/blob/master/ecosystem/sep-0001.md) and [federation](https://github.com/paydex-core/paydex-go/paydex-go-protocol/blob/master/ecosystem/sep-0002.md).

## Horizon SDK

The Horizon SDK is composed of two complementary libraries: `txnbuild` + `horizonclient`.
The `txnbuild` ([source](https://github.com/paydex-core/paydex-go/tree/master/txnbuild), [docs](https://godoc.org/github.com/paydex-core/paydex-go/txnbuild)) package enables the construction, signing and encoding of Paydex [transactions](https://www.paydex.org/developers/guides/concepts/transactions.html) and [operations](https://www.paydex.org/developers/guides/concepts/list-of-operations.html) in Go. The `horizonclient` ([source](https://github.com/paydex/go/tree/master/clients/horizonclient), [docs](https://godoc.org/github.com/paydex-core/paydex-go/clients/horizonclient)) package provides a web client for interfacing with [Horizon](https://www.paydex.org/developers/guides/get-started/) server REST endpoints to retrieve ledger information, and to submit transactions built with `txnbuild`.

## List of major SDK packages

- `horizonclient` ([source](https://github.com/paydex-core/paydex-go/tree/master/clients/horizonclient), [docs](https://godoc.org/github.com/paydex-core/paydex-go/clients/horizonclient)) - programmatic client access to Horizon
- `txnbuild` ([source](https://github.com/paydex-core/paydex-go/tree/master/txnbuild), [docs](https://godoc.org/github.com/paydex-core/paydex-go/txnbuild)) - construction, signing and encoding of Paydex transactions and operations
- `paydextoml` ([source](https://github.com/paydex-core/paydex-go/tree/master/clients/paydextoml), [docs](https://godoc.org/github.com/paydex-core/paydex-go/clients/paydextoml)) - parse [Paydex.toml](../../guides/concepts/paydex-toml.md) files from the internet
- `federation` ([source](https://godoc.org/github.com/paydex-core/paydex-go/clients/federation)) - resolve federation addresses  into paydex account IDs, suitable for use within a transaction

