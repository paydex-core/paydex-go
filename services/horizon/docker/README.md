# Running Paydex with Docker Compose

## Dependencies

The only dependency you will need to install is [Docker](https://www.docker.com/products/docker-desktop).

## Start script

[start.sh](./start.sh) will setup the env file and run docker-compose to start the Paydex docker containers. Feel free to use this script, otherwise continue with the next two steps.

## Set up a .env file

Mac OS X and Windows users should create an [`.env`](https://docs.docker.com/compose/environment-variables/#the-env_file-configuration-option) file which consists of:

`NETWORK_MODE=bridge`

Linux users should also create an `.env` file. However, the contents of the file should look like:

`NETWORK_MODE=host`

Additionally, you will need to add `127.0.0.1 host.docker.internal` to the `/etc/hosts` file on your linux machine.

If https://github.com/docker/for-linux/issues/264 is ever fixed then it won't be necessary to alias `host.docker.internal` to localhost and there won't be any differences between the Linux and Mac OS X / Windows configurations.


## Run docker-compose

Run the following command to start all the Paydex docker containers:

```
docker-compose up -d --build
```

Horizon will be exposed on port 8000. Paydex Core will be exposed on port 11626. The Paydex Core postgres instance will be exposed on port 5641.
The Horizon postgres instance will be exposed on port 5432.

## Swapping in a local service

If you're developing a service locally you may want to run that service locally while also being able to interact with the other Paydex components running in Docker. You can do that by stopping the container corresponding to the service you're developing.

For example, to run Horizon locally from source, you would perform the following steps:

```
# stop horizon in docker-compose
docker-compose stop horizon
```

Now you can run horizon locally in vscode using the following configuration:
```
    {
        "name": "Launch",
        "type": "go",
        "request": "launch",
        "mode": "debug",
        "remotePath": "",
        "port": 2345,
        "host": "127.0.0.1",
        "program": "${workspaceRoot}/services/horizon/main.go",
        "env": {
            "DATABASE_URL": "postgres://postgres@localhost:5432/horizon?sslmode=disable",
            "PAYDEX_CORE_DATABASE_URL": "postgres://postgres:mysecretpassword@localhost:5641/paydex?sslmode=disable",
            "NETWORK_PASSPHRASE": "Test SDF Network ; September 2015",
            "PAYDEX_CORE_URL": "http://localhost:11626",
            "INGEST": "true",
        },
        "args": []
    }
```

Similarly, to run Paydex core locally from source and have it interact with Horizon in docker, all you need to do is run `docker-compose stop core` before running Paydex core from source.

## Connecting to the Paydex Public Network

By default, the Docker Compose file configures Paydex Core to connect to the Paydex test network. If you would like to run the docker containers against the
Paydex public network, set the `core` container's env_file to `./paydex-core-pubnet.env` instead of `./paydex-core-testnet.env`. You will also need to
change the `NETWORK_PASSPHRASE` variable in horizon to `Public Global Paydex Network ; September 2015`.

When you switch between the Paydex test network and the Paydex public network, or vice versa, you will need to clear the Paydex Core and Paydex Horizon
databases. You can wipe out the databases by running `docker-compose down -v`.

## Using a specific version of Paydex Core

By default the Docker Compose file is configured to use the latest version of Paydex Core. To use a specific version, you can edit [docker-compose.yml](./docker-compose.yml) and set the appropriate [tag](https://hub.docker.com/r/paydex/paydex-core/tags) on the Paydex Core docker image