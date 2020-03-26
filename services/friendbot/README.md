# Friendbot Service for the Paydex Test Network

Friendbot helps users of the Paydex testnet by exposing a REST endpoint that creates & funds new accounts.

Horizon needs to be started with the following command line param: --friendbot-url="http://localhost:8004/"
This will forward any query params received against /friendbot to the friendbot instance.
The ideal setup for horizon is to proxy all requests to the /friendbot url to the friendbot service
