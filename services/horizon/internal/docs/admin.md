---
title: Horizon Administration Guide
---
## Horizon Administration Guide

Horizon is responsible for providing an HTTP API to data in the Paydex network. It ingests and re-serves the data produced by the paydex network in a form that is easier to consume than the performance-oriented data representations used by paydex-core.

This document describes how to administer a **production** Horizon instance. If you are just starting with Horizon and want to try it out, consider the [Quickstart Guide](quickstart.md) instead. For information about developing on the Horizon codebase, check out the [Development Guide](developing.md).

## Why run Horizon?

The Paydex Development Foundation runs two Horizon servers, one for the public network and one for the test network, free for anyone's use at https://horizon.paydex.org and https://horizon-testnet.paydex.org.  These servers should be fine for development and small scale projects, but it is not recommended that you use them for production services that need strong reliability.  By running Horizon within your own infrastructure provides a number of benefits:

  - Multiple instances can be run for redundancy and scalability.
  - Request rate limiting can be disabled.
  - Full operational control without dependency on the Paydex Development Foundations operations.

## Prerequisites

Horizon is dependent upon a paydex-core server.  Horizon needs access to both the SQL database and the HTTP API that is published by paydex-core. See [the administration guide](https://www.paydex.org/developers/paydex-core/learn/admin.html
) to learn how to set up and administer a paydex-core server.  Secondly, Horizon is dependent upon a postgres server, which it uses to store processed core data for ease of use. Horizon requires postgres version >= 9.5.

In addition to the two prerequisites above, you may optionally install a redis server to be used for rate limiting requests.

## Installing

To install Horizon, you have a choice: either downloading a [prebuilt release for your target architecture](https://github.com/paydex-core/paydex-go/releases) and operation system, or [building Horizon yourself](#Building).  When either approach is complete, you will find yourself with a directory containing a file named `horizon`.  This file is a native binary.

After building or unpacking Horizon, you simply need to copy the native binary into a directory that is part of your PATH.  Most unix-like systems have `/usr/local/bin` in PATH by default, so unless you have a preference or know better, we recommend you copy the binary there.

To test the installation, simply run `horizon --help` from a terminal.  If the help for Horizon is displayed, your installation was successful. Note: some shells, such as zsh, cache PATH lookups.  You may need to clear your cache  (by using `rehash` in zsh, for example) before trying to run `horizon --help`.


## Building

Should you decide not to use one of our prebuilt releases, you may instead build Horizon from source.  To do so, you need to install some developer tools:

- A unix-like operating system with the common core commands (cp, tar, mkdir, bash, etc.)
- A compatible distribution of Go (Go 1.12 or later)
- [git](https://git-scm.com/)
- [mercurial](https://www.mercurial-scm.org/)

1. See the details in [README.md](../../../../README.md#dependencies) for installing dependencies.
2. Compile the Horizon binary: `go install github.com/paydex-core/paydex-go/services/horizon`. You should see the `horizon` binary in `$GOPATH/bin`.
3. Add Go binaries to your PATH in your `bashrc` or equivalent, for easy access: `export PATH=${GOPATH//://bin:}/bin:$PATH`

Open a new terminal. Confirm everything worked by running `horizon --help` successfully.

Note:  Building directly on windows is not supported.


## Configuring

Horizon is configured using command line flags or environment variables.  To see the list of command line flags that are available (and their default values) for your version of Horizon, run:

`horizon --help`

As you will see if you run the command above, Horizon defines a large number of flags, however only three are required:

| flag                    | envvar                      | example                              |
|-------------------------|-----------------------------|--------------------------------------|
| `--db-url`              | `DATABASE_URL`              | postgres://localhost/horizon_testnet |
| `--paydex-core-db-url`  | `PAYDEX_CORE_DATABASE_URL` | postgres://localhost/core_testnet    |
| `--paydex-core-url`     | `PAYDEX_CORE_URL`          | http://localhost:11626               |

`--db-url` specifies the Horizon database, and its value should be a valid [PostgreSQL Connection URI](http://www.postgresql.org/docs/9.2/static/libpq-connect.html#AEN38419).  `--paydex-core-db-url` specifies a paydex-core database which will be used to load data about the paydex ledger.  Finally, `--paydex-core-url` specifies the HTTP control port for an instance of paydex-core.  This URL should be associated with the paydex-core that is writing to the database at `--paydex-core-db-url`.

Specifying command line flags every time you invoke Horizon can be cumbersome, and so we recommend using environment variables.  There are many tools you can use to manage environment variables:  we recommend either [direnv](http://direnv.net/) or [dotenv](https://github.com/bkeepers/dotenv).  A template configuration that is compatible with dotenv can be found in the [Horizon git repo](https://github.com/paydex-core/paydex-go/blob/master/services/horizon/.env.template).



## Preparing the database

Before the Horizon server can be run, we must first prepare the Horizon database.  This database will be used for all of the information produced by Horizon, notably historical information about successful transactions that have occurred on the paydex network.

To prepare a database for Horizon's use, first you must ensure the database is blank.  It's easiest to simply create a new database on your postgres server specifically for Horizon's use.  Next you must install the schema by running `horizon db init`.  Remember to use the appropriate command line flags or environment variables to configure Horizon as explained in [Configuring ](#Configuring).  This command will log any errors that occur.

### Postgres configuration

It is recommended to set `random_page_cost=1` in Postgres configuration if you are using SSD storage. With this setting Query Planner will make a better use of indexes, expecially for `JOIN` queries. We have noticed a huge speed improvement for some queries.

## Running

Once your Horizon database is configured, you're ready to run Horizon.  To run Horizon you simply run `horizon` or `horizon serve`, both of which start the HTTP server and start logging to standard out.  When run, you should see some output that similar to:

```
INFO[0000] Starting horizon on :8000                     pid=29013
```

The log line above announces that Horizon is ready to serve client requests. Note: the numbers shown above may be different for your installation.  Next we can confirm that Horizon is responding correctly by loading the root resource.  In the example above, that URL would be [http://127.0.0.1:8000/] and simply running `curl http://127.0.0.1:8000/` shows you that the root resource can be loaded correctly.

If you didn't set up a paydex-core yet, you may see an error like this:
```
ERRO[2019-05-06T16:21:14.126+08:00] Error getting core latest ledger err="get failed: pq: relation \"ledgerheaders\" does not exist"
```
Horizon requires a functional paydex-core. Go back and set up paydex-core as described in the admin guide. In particular, you need to initialise the database as [described here](https://www.paydex.org/developers/paydex-core/software/admin.html#database-and-local-state).

## Ingesting live paydex-core data

Horizon provides most of its utility through ingested data.  Your Horizon server can be configured
to listen for and ingest transaction results from the connected paydex-core.  We recommend that
within your infrastructure you run one (and only one) Horizon process that is configured in this
way. While running multiple ingestion processes will not corrupt the Horizon database, your error
logs will quickly fill up as the two instances race to ingest the data from paydex-core. A notable
exception to this is when you are reingesting data, which we recommend using multiple processes for
speed (more on this below).

To enable ingestion, you must either pass `--ingest=true` on the command line or set the `INGEST`
environment variable to "true".

### Ingesting historical data

To enable ingestion of historical data from paydex-core you need to run `horizon db backfill NUM_LEDGERS`. If you're running a full validator with published history archive, for example, you might want to ingest all of history. In this case your `NUM_LEDGERS` should be slightly higher than the current ledger id on the network. You can run this process in the background while your Horizon server is up. This continuously decrements the `history.elder_ledger` in your /metrics endpoint until `NUM_LEDGERS` is reached and the backfill is complete.

### Reingesting Ledgers
A notable exception to running only a single Horizon process is when you are reingesting ledgers,
which we recommend you run multiple processes for in order to dramatically speed up re-ingestion
time. This is done through the `horizon db range [START_LEDGER] [END_LEDGER]` command, which could
be run as follows:

```
horizon1> horizon db reingest range 1 10000
horizon2> horizon db reingest range 10001 20000
horizon3> horizon db reingest range 20001 30000
# ... etc.
```

This allows reingestion to be split up and done in parallel by multiple Horizon processes, and is
available as of Horizon [0.17.4](https://github.com/paydex-core/paydex-go/releases/tag/horizon-v0.17.4).

### Managing storage for historical data

Over time, the recorded network history will grow unbounded, increasing storage used by the database. Horizon expands the data ingested from paydex-core and needs sufficient disk space. Unless you need to maintain a history archive you may configure Horizon to only retain a certain number of ledgers in the database. This is done using the `--history-retention-count` flag or the `HISTORY_RETENTION_COUNT` environment variable. Set the value to the number of recent ledgers you wish to keep around, and every hour the Horizon subsystem will reap expired data.  Alternatively, you may execute the command `horizon db reap` to force a collection.

### Surviving paydex-core downtime

Horizon tries to maintain a gap-free window into the history of the paydex-network.  This reduces the number of edge cases that Horizon-dependent software must deal with, aiming to make the integration process simpler.  To maintain a gap-free history, Horizon needs access to all of the metadata produced by paydex-core in the process of closing a ledger, and there are instances when this metadata can be lost.  Usually, this loss of metadata occurs because the paydex-core node went offline and performed a catchup operation when restarted.

To ensure that the metadata required by Horizon is maintained, you have several options: You may either set the `CATCHUP_COMPLETE` paydex-core configuration option to `true` or configure `CATCHUP_RECENT` to determine the amount of time your paydex-core can be offline without having to rebuild your Horizon database.

Unless your node is a full validator and archive publisher we _do not_ recommend using the `CATCHUP_COMPLETE` method, as this will force paydex-core to apply every transaction from the beginning of the ledger, which will take an ever increasing amount of time. Instead, we recommend you set the `CATCHUP_RECENT` config value. To do this, determine how long of a downtime you would like to survive (expressed in seconds) and divide by ten.  This roughly equates to the number of ledgers that occur within your desired grace period (ledgers roughly close at a rate of one every ten seconds).  With this value set, paydex-core will replay transactions for ledgers that are recent enough, ensuring that the metadata needed by Horizon is present.

### Correcting gaps in historical data

We recommend you configure the HISTORY_RETENTION_COUNT in Horizon to a value less than or equal to the configured value for CATCHUP_RECENT in paydex-core.  Given this situation any downtime that would cause a ledger gap will require a downtime greater than the amount of historical data retained by Horizon.  To re-establish continuity:

1.  Stop Horizon.
2.  Run `horizon db reap` to clear the historical database.
3.  Clear the cursor for Horizon by running `paydex-core -c "dropcursor?id=HORIZON"` (ensure capitilization is maintained).
4.  Clear ledger metadata from before the gap by running `paydex-core -c "maintenance?queue=true"`.
5.  Restart Horizon.

## Managing Stale Historical Data

Horizon ingests ledger data from a connected instance of paydex-core.  In the event that paydex-core stops running (or if Horizon stops ingesting data for any other reason), the view provided by Horizon will start to lag behind reality.  For simpler applications, this may be fine, but in many cases this lag is unacceptable and the application should not continue operating until the lag is resolved.

To help applications that cannot tolerate lag, Horizon provides a configurable "staleness" threshold.  Given that enough lag has accumulated to surpass this threshold (expressed in number of ledgers), Horizon will only respond with an error: [`stale_history`](./errors/stale-history.md).  To configure this option, use either the `--history-stale-threshold` command line flag or the `HISTORY_STALE_THRESHOLD` environment variable.  NOTE:  non-historical requests (such as submitting transactions or finding payment paths) will not error out when the staleness threshold is surpassed.

## Monitoring

To ensure that your instance of Horizon is performing correctly we encourage you to monitor it, and provide both logs and metrics to do so.

Horizon will output logs to standard out.  Information about what requests are coming in will be reported, but more importantly, warnings or errors will also be emitted by default.  A correctly running Horizon instance will not output any warning or error log entries.

Metrics are collected while a Horizon process is running and they are exposed at the `/metrics` path.  You can see an example at (https://horizon-testnet.paydex.org/metrics).

Below we present a few standard log entries with associated fields. You can use them to build metrics and alerts. We present below some examples. Please note that this represents Horizon app metrics only. You should also monitor your hardware metrics like CPU or RAM Utilization.

### Starting HTTP request

| Key              | Value                                                                                          |
|------------------|------------------------------------------------------------------------------------------------|
| **`msg`**        | **`Starting request`**                                                                         |
| `client_name`    | Value of `X-Client-Name` HTTP header representing client name                                  |
| `client_version` | Value of `X-Client-Version` HTTP header representing client version                            |
| `app_name`       | Value of `X-App-Name` HTTP header representing app name                                        |
| `app_version`    | Value of `X-App-Version` HTTP header representing app version                                  |
| `forwarded_ip`   | First value of `X-Forwarded-For` header                                                        |
| `host`           | Value of `Host` header                                                                         |
| `ip`             | IP of a client sending HTTP request                                                            |
| `ip_port`        | IP and port of a client sending HTTP request                                                   |
| `method`         | HTTP method (`GET`, `POST`, ...)                                                               |
| `path`           | Full request path, including query string (ex. `/transactions?order=desc`)                     |
| `streaming`      | Boolean, `true` if request is a streaming request                                              |
| `referer`        | Value of `Referer` header                                                                      |
| `req`            | Random value that uniquely identifies a request, attached to all logs within this HTTP request |

### Finished HTTP request

| Key              | Value                                                                                          |
|------------------|------------------------------------------------------------------------------------------------|
| **`msg`**        | **`Finished request`**                                                                         |
| `bytes`          | Number of response bytes sent                                                                  |
| `client_name`    | Value of `X-Client-Name` HTTP header representing client name                                  |
| `client_version` | Value of `X-Client-Version` HTTP header representing client version                            |
| `app_name`       | Value of `X-App-Name` HTTP header representing app name                                        |
| `app_version`    | Value of `X-App-Version` HTTP header representing app version                                  |
| `duration`       | Duration of request in seconds                                                                 |
| `forwarded_ip`   | First value of `X-Forwarded-For` header                                                        |
| `host`           | Value of `Host` header                                                                         |
| `ip`             | IP of a client sending HTTP request                                                            |
| `ip_port`        | IP and port of a client sending HTTP request                                                   |
| `method`         | HTTP method (`GET`, `POST`, ...)                                                               |
| `path`           | Full request path, including query string (ex. `/transactions?order=desc`)                     |
| `route`          | Route pattern without query string (ex. `/accounts/{id}`)                                      |
| `status`         | HTTP status code (ex. `200`)                                                                   |
| `streaming`      | Boolean, `true` if request is a streaming request                                              |
| `referer`        | Value of `Referer` header                                                                      |
| `req`            | Random value that uniquely identifies a request, attached to all logs within this HTTP request |

### Processing (ingesting) a new ledger

| Key       | Value                    |
|-----------|--------------------------|
| **`msg`** | **`Reading new ledger`** |
| `ledger`  | Ledger sequence          |

### Finished processing (ingesting) a new ledger

| Key            | Value                                |
|----------------|--------------------------------------|
| **`msg`**      | **`Finished processing ledger`**     |
| `ledger`       | Ledger sequence                      |
| `duration`     | Duration in seconds                  |
| `transactions` | Number of transactions in the ledger |


### Metrics

Using the entries above you can build metrics that will help understand performance of a given Horizon node, some examples below:
* Number of requests per minute.
* Number of requests per route (the most popular routes).
* Average response time per route.
* Maximum response time for non-streaming requests.
* Number of streaming vs. non-streaming requests.
* Number of rate-limited requests.
* List of rate-limited IPs.
* Unique IPs.
* The most popular SDKs/apps sending requests to a given Horizon node.
* Average ingestion time of a ledger.
* Average ingestion time of a transaction.

### Alerts

Below we present example alerts with potential cause and solution. Feel free to add more alerts using your metrics.

Alert | Cause | Solution
-|-|-
Spike in number of requests | Potential DoS attack | Lower rate-limiting threshold
Large number of rate-limited requests | Rate-limiting threshold too low | Increase rate-limiting threshold
Ingestion is slow | Horizon server spec too low | Increase hardware spec
Spike in average response time of a single route | Possible bug in a code responsible for rendering a route | Report an issue in Horizon repository.

## I'm Stuck! Help!

If any of the above steps don't work or you are otherwise prevented from correctly setting up
Horizon, please come to our community and tell us. Either
[post a question at our Stack Exchange](https://paydex.stackexchange.com/) or
[chat with us on Keybase in #dev_discussion](https://keybase.io/team/paydex.public) to ask for
help.