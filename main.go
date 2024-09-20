// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"context"
	"d7024e_group04/api"
	"d7024e_group04/cli"
	"d7024e_group04/env"
	"d7024e_group04/internal/client"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/node"
	"d7024e_group04/internal/server"
	"d7024e_group04/internal/store"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}

	address := host + ":" + strconv.Itoa(env.Port)
	rootCtx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	errGroup, errCtx := errgroup.WithContext(rootCtx)

	id := kademliaid.NewRandomKademliaID()
	c := contact.NewContact(id, address)

	routingTable := routingtable.NewRoutingTable(c)

	memoryStore := store.NewMemoryStore()

	client := client.NewClient()

	node := node.New(client, routingTable, memoryStore)

	errGroup.Go(func() error {
		// TODO node bootstrap stuff
		return err
	})

	if err = errGroup.Wait(); err != nil {
		log.Fatalf("bootstrap failed, %v", err)
	}

	server := server.NewServer(routingTable, memoryStore)
	errGroup.Go(func() error {
		return server.Start(errCtx)
	})

	// REST API
	handler := api.NewHandler(node)
	errGroup.Go(func() error {
		return handler.ListenAndServe(errCtx)
	})

	// CLI loop
	errGroup.Go(func() error {
		return cli.InputLoop(errCtx, cancelCtx, os.Stdin, node)
	})

	err = errGroup.Wait()
	cancelCtx()
	log.Printf("node shutdown, reason: %v", err)
}
