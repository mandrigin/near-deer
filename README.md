# Near Deer
A simple utility to compare the node's latest block with the source of truth (checks against https://rpc.<network>.near.org/status)
  
  
It provides an HTTP endpoint `/health` that returns either 200OK is the node is in sync, and 500 if it is not.


Example usage:

all custom params
```
go run ./cmd/deer -network=testnet -node=http://localhost:3030 -port 30303 -host 0.0.0.0 -threshold 10
```

checking mainnet for a local node with rpc on 3030 and hosting the app on 8080 of localhost, at most 3 blocks behind (all defaults)
```
go run ./cmd/deer
```
