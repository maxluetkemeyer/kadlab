// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/kademlia"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}

	address := host + ":" + strconv.Itoa(env.Port)
	rootCtx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGTERM)

	node := kademlia.NewNode(address)

	err = node.Start(rootCtx)

	cancelCtx()
	log.Printf("node shutdown, reason: %v", err)
}
