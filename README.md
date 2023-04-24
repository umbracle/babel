# Babel

Babel is a standard interface to query the sync status of a blockchain node.

## Usage

Run as a `Grpc` server:

```
$ go run cmd/main.go --plugin ethereum_el server [url=http://localhost:8545]
```

## Supported nodes

The following blockchain/clients have support in Babel:

- Ethereum Execution client (`ethereum_el`): It uses the [Json-RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/) specification. Any `geth` clone should also be compatible with this plugin.
- Ethereum Consensus client (`ethereum_cl`): It uses the [Beacon-API](https://ethereum.github.io/beacon-APIs/#/Node/getSyncingStatus) specification.
