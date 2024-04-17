package main

import (
	"fmt"
	"icescream/getBlockStd/getBlock"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	nodeEndpoint := "https://go.getblock.io/bd9745e549b34652a20fb8cbf5f10bbf"
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
