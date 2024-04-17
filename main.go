package main

import (
	"fmt"
	"getBlock/getBlock"
	_ "net/http/pprof"
	"os"
	"time"
)

func main() {
	os.Getenv("API_KEY")

	nodeEndpoint := "https://go.getblock.io/e3c21d46fb4545d6bff3df29c42fd4a0"

	//nodeEndpoint := "https://api.securerpc.com/v1"
	_ = nodeEndpoint

	start := time.Now()
	service := getBlock.New(nodeEndpoint)
	topFive := service.GetTop5Addresses()
	end := time.Now()

	fmt.Printf("bench: %f seconds\n", end.Sub(start).Seconds())

	for _, val := range topFive {
		fmt.Printf("address: %s activity: %d\n", val.Address, val.Count)
	}
}
