# near-deer
simple utility to compare the node's latest block with the source of truth (checks against https://rpc.<network>.near.org/status)


Example usage:

all custom params
```
go run ./cmd/deer -network=testnet -node=http://localhost:3030 -port 30303 -host 0.0.0.0 -threshold 10
```

checking mainnet for a local node with rpc on 3030 and hosting the app on 8080 of localhost, at most 3 blocks behind (all defaults)
```
go run ./cmd/deer
```
