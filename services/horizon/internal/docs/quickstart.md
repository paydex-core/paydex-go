---
title: Horizon Quickstart
---
## Horizon Quickstart
This document describes how to quickly set up a **test** Paydex Core + Horizon node, that you can play around with to get a feel for how a paydex node operates. **This configuration is not secure!** It is **not** intended as a guide for production administration.

For detailed information about running Horizon and Paydex Core safely in production see the [Horizon Administration Guide](admin.md) and the [Paydex Core Administration Guide](https://www.paydex.org/developers/paydex-core/software/admin.html).

If you're ready to roll up your sleeves and dig into the code, check out the [Developer Guide](developing.md).

### Install and run the Quickstart Docker Image
The fastest way to get up and running is using the [Paydex Quickstart Docker Image](https://github.com/paydex/docker-paydex-core-horizon). This is a Docker container that provides both `paydex-core` and `horizon`, pre-configured for testing.

1. Install [Docker](https://www.docker.com/get-started).
2. Verify your Docker installation works: `docker run hello-world`
3. Create a local directory that the container can use to record state. This is helpful because it can take a few minutes to sync a new `paydex-core` with enough data for testing, and because it allows you to inspect and modify the configuration if needed. Here, we create a directory called `paydex` to use as the persistent volume:
`cd $HOME; mkdir paydex`
4. Download and run the Paydex Quickstart container, replacing `USER` with your username:

```bash
docker run --rm -it -p "8000:8000" -p "11626:11626" -p "11625:11625" -p"8002:5432" -v $HOME/paydex:/opt/paydex --name paydex paydex/quickstart --testnet
```

You can check out Paydex Core status by browsing to http://localhost:11626.

You can check out your Horizon instance by browsing to http://localhost:8000.

You can tail logs within the container to see what's going on behind the scenes:
```bash
docker exec -it paydex /bin/bash
supervisorctl tail -f paydex-core
supervisorctl tail -f horizon stderr
```

On a modern laptop this test setup takes about 15 minutes to synchronise with the last couple of days of testnet ledgers. At that point Horizon will be available for querying. 

See the [Quickstart Docker Image](https://github.com/paydex-core/docker-paydex-core-horizon) documentation for more details, and alternative ways to run the container. 

You can test your Horizon instance with a query like: http://localhost:8001/transactions?cursor=&limit=10&order=asc. Use the [Paydex Laboratory](https://www.paydex.org/laboratory/) to craft other queries to try out,
and read about the available endpoints and see examples in the [Horizon API reference](https://www.paydex.org/developers/horizon/reference/).

