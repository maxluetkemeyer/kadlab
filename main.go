// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"context"
	"d7024e_group04/kademlia"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const BucketSize = "20"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Setenv("BUCKET_SIZE", BucketSize) // TODO move this to config file?

	address := host + ":" + port
	rootCtx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGTERM)

	node := kademlia.NewNode(address)

	err = node.Start(rootCtx)

	cancelCtx()
	log.Printf("node shutdown, reason: %v", err)
}
