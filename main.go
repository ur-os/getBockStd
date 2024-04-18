package main

import (
	"getBlock/getBlock"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	service := getBlock.NewService()

	http.HandleFunc("/getTopFiveUserActivity", service.GetTopFiveUserActivity) // Update this line of code

	port := os.Getenv("GET_BLOCK_PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
