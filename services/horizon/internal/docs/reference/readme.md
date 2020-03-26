---
title: Overview
---

Horizon is an API server for the Paydex ecosystem.  It acts as the interface between [paydex-core](https://github.com/paydex/paydex-core) and applications that want to access the Paydex network. It allows you to submit transactions to the network, check the status of accounts, subscribe to event streams, etc. See [an overview of the Paydex ecosystem](https://www.paydex.org/developers/guides/) for details of where Horizon fits in. You can also watch a [talk on Horizon](https://www.youtube.com/watch?v=AtJ-f6Ih4A4) by Paydex.org developer Scott Fleckenstein:

[![Horizon: API webserver for the Paydex network](https://img.youtube.com/vi/AtJ-f6Ih4A4/sddefault.jpg "Horizon: API webserver for the Paydex network")](https://www.youtube.com/watch?v=AtJ-f6Ih4A4)

Horizon provides a RESTful API to allow client applications to interact with the Paydex network. You can communicate with Horizon using cURL or just your web browser. However, if you're building a client application, you'll likely want to use a Paydex SDK in the language of your client.
SDF provides a [JavaScript SDK](https://www.paydex.org/developers/js-payde-sdk/learn/index.html) for clients to use to interact with Horizon.

SDF runs a instance of Horizon that is connected to the test net: [https://horizon-testnet.paydex.org/](https://horizon-testnet.paydex.org/) and one that is connected to the public Paydex network:
[https://horizon.paydex.org/](https://horizon.paydex.org/).

## Libraries

SDF maintained libraries:<br />
- [JavaScript](https://github.com/paydex-core/js-paydex-sdk)
- [Go](https://github.com/paydex-core/paydex-go/tree/master/clients/horizonclient)