package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	resty "github.com/go-resty/resty/v2"
)

var (
	flagNetwork     = flag.String("network", "mainnet", "select network: mainnet or testnet")
	flagNodeAddress = flag.String("node", "http://localhost:3030", "node to check against")
	flagThreshold   = flag.Int("threshold", 3, "health threshold, how many blocks behind is it okay to be for a local node")
	flagPort        = flag.Int("port", 8080, "port on which to serve")
	flagHost        = flag.String("host", "localhost", "host on which to serve (default=localhost)")
)

func main() {
	flag.Parse()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		fmt.Println("DEBUG: initiating healthcheck")
		err := checkNodeHealth()

		if err == nil {
			fmt.Println("INFO: node is healthy")
			c.JSON(200, gin.H{})
		} else {
			fmt.Println("ERR: while checking health:", err)
			c.JSON(500, gin.H{"error": err.Error()})
		}
	})

	r.Run(fmt.Sprintf("%s:%v", *flagHost, *flagPort))
}

func checkNodeHealth() error {

	sourceOfTruthBlock, err := getLatestBlockNumber(
		getSourceOfTruthAddressForNetwork(*flagNetwork),
	)
	if err != nil {
		return fmt.Errorf("can't get source of truth block number. err=%w", err)
	}
	localNodeBlock, err := getLatestBlockNumber(
		*flagNodeAddress,
	)
	if err != nil {
		return fmt.Errorf("can't get local node block number. err=%w", err)
	}

	fmt.Println("DEBUG", "sot", sourceOfTruthBlock, "local", localNodeBlock)

	diff := sourceOfTruthBlock - localNodeBlock

	if diff > *flagThreshold { //negative is fine
		return fmt.Errorf("the local node is too far away from the source of truth. diff=%v", diff)
	}

	fmt.Println("node is healthy", "difference", diff)
	return nil
}

type StatusResult struct {
	SyncInfo struct {
		LatestBlockHeight int `json:"latest_block_height"`
	} `json:"sync_info"`
}

func getLatestBlockNumber(address string) (int, error) {
	fmt.Println("DEBUG: checking block for address", address)
	client := resty.New()

	if !strings.HasSuffix(address, "/") {
		address += "/"
	}

	resp, err := client.R().SetResult(&StatusResult{}).Get(fmt.Sprintf("%sstatus", address))
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
