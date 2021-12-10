package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagNetwork     = flag.String("network", "mainnet", "select network: mainnet or testnet")
	flagNodeAddress = flag.String("node", "localhost:3030", "node to check against")
	flagThreshold   = flag.Int("threshold", 3, "health threshold, how many blocks behind is it okay to be for a local node")
)

func main() {
	flag.Parse()
	fmt.Println("hello", *flagNetwork, *flagNodeAddress)
	sourceOfTruthBlock, err := getLatestBlockNumber(
		getSourceOfTruthAddressForNetwork(*flagNetwork),
	)
	if err != nil {
		fmt.Printf("error while getting source of truth block number: %v", err)
		os.Exit(1)
	}
	localNodeBlock, err := getLatestBlockNumber(
		*flagNodeAddress,
	)
	if err != nil {
		fmt.Printf("error while getting local node block number: %v", err)
		os.Exit(2)
	}

	if sourceOfTruthBlock-localNodeBlock > *flagThreshold { //negative is fine
		fmt.Printf("error, the local node is too far away from the source of truth")
		os.Exit(3)
	}
}

func getLatestBlockNumber(address string) (int, error) {
	return 0, nil
}

func getSourceOfTruthAddressForNetwork(network string) string {
	if len(network) == 0 {
		return ""
	}

	return fmt.Sprintf("https://rpc.%s.near.org/status", network)
}
