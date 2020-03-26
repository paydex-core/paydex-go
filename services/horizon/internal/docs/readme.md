---
title: Horizon
---

Horizon is the server for the client facing API for the Paydex ecosystem.  It acts as the interface between [paydex-core](https://www.paydex.org/developers/learn/paydex-core) and applications that want to access the Paydex network. It allows you to submit transactions to the network, check the status of accounts, subscribe to event streams, etc. See [an overview of the Paydex ecosystem](https://www.paydex.org/developers/guides/) for more details.

You can interact directly with horizon via curl or a web browser but SDF provides a [JavaScript SDK](https://www.paydex.org/developers/js-paydex-sdk/learn/) for clients to use to interact with Horizon.

SDF runs a instance of Horizon that is connected to the test net [https://horizon-testnet.paydex.org/](https://horizon-testnet.paydex.org/).

## Libraries

SDF maintained libraries:<br />
- [JavaScript](https://github.com/paydex/js-paydex-sdk)
- [Go](https://github.com/paydex-core/paydex-go/tree/master/clients/horizonclient)
- [Java](https://github.com/paydex-core/java-paydex-sdk)

Community maintained libraries (in various states of completeness) for interacting with Horizon in other languages:<br>
- [Python](https://github.com/PaydexCN/py-paydex-base)
- [C# .NET Core 2.x](https://github.com/elucidsoft/dotnetcore-paydex-sdk)
- [Ruby](https://github.com/bloom-solutions/ruby-paydex-sdk)
- [iOS and macOS](https://github.com/Soneso/paydex-ios-mac-sdk)
- [Scala SDK](https://github.com/synesso/scala-paydex-sdk)
- [C++ SDK](https://github.com/bnogalm/PaydexQtSDK)
