package main

import (
	"flag"
	"fmt"
	"os"

	resty "github.com/go-resty/resty/v2"
)

var (
	flagNetwork     = flag.String("network", "mainnet", "select network: mainnet or testnet")
	flagNodeAddress = flag.String("node", "http://localhost:3030", "node to check against")
	flagThreshold   = flag.Int("threshold", 3, "health threshold, how many blocks behind is it okay to be for a local node")
)

func main() {
	flag.Parse()

	fmt.Println("checking health", *flagNetwork, *flagNodeAddress)

	sourceOfTruthBlock, err := getLatestBlockNumber(
		getSourceOfTruthAddressForNetwork(*flagNetwork),
	)
	if err != nil {
		fmt.Printf("ERR: can't get source of truth block number: %v", err)
		os.Exit(1)
	}
	localNodeBlock, err := getLatestBlockNumber(
		*flagNodeAddress,
	)
	if err != nil {
		fmt.Printf("ERR: can't get local node block number: %v", err)
		os.Exit(2)
	}

	diff := sourceOfTruthBlock - localNodeBlock

	if diff > *flagThreshold { //negative is fine
		fmt.Println("ERR: The local node is too far away from the source of truth.", "difference", diff)
		os.Exit(3)
	}

	fmt.Println("node is healthy", "difference", diff)
}

type StatusResult struct {
	SyncInfo struct {
		LatestBlockHeight int `json:"latest_block_height"`
	} `json:"sync_info"`
}

func getLatestBlockNumber(address string) (int, error) {
	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().SetResult(&StatusResult{}).Get(fmt.Sprintf("%s/status", address))
	if err != nil {
		return 0, err
	}

	result := resp.Result().(*StatusResult)
	if result == nil {
		return 0, nil
	}

	return result.SyncInfo.LatestBlockHeight, nil
}

func getSourceOfTruthAddressForNetwork(network string) string {
	if len(network) == 0 {
		return ""
	}

	return fmt.Sprintf("https://rpc.%s.near.org", network)
}
