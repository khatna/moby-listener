# Moby Listener Service

The **Moby Listener Service** is a gRPC service that listens to all broadcasted transactions on the Bitcoin network, and stream large transactions via a **gRPC** service.

## How it Works
It is possible to publish all raw transactions using ZMQ on Bitcoin Core. Furthermore, raw transactions can be decoded into JSON via the [Bitcoin Core RPC API](https://developer.bitcoin.org/reference/rpc/).

`moby-listener` listens to these raw transactions, uses the RPC API to decode them and streams them to connected clients using a basic gRPC server (defined in `proto/transaction.proto`).

## Setup
The following environment variables should be set:
```
BTC_RPC_HOST
BTC_RPC_USER
BTC_RPC_PASS
BTC_TCP_ENDPOINT

```

Access to a bitcoin (pruned) node is required. Furthermore, the following configurations must be specified in the node:
```
zmqpubrawtx=$BTC_TCP_ENDPOINT
rpcuser=$BTC_RPC_USER
rpcpassword=$BTC_RPC_PASS

```

Finally, [libzmq](https://zeromq.org/download/) should be installed.