// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"d7024e_group04/api"
	"d7024e_group04/cli"
	"d7024e_group04/env"
	"d7024e_group04/internal/client"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/network"
	"d7024e_group04/internal/node"
	"d7024e_group04/internal/server"
	"d7024e_group04/internal/store"

	"golang.org/x/sync/errgroup"
)

func main() {
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}

	ip, err := net.LookupIP(host)
	if err != nil {
		panic("bad ip")
	}

	address := ip[0].String() + ":" + strconv.Itoa(env.Port)

	rootCtx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	errGroup, errCtx := errgroup.WithContext(rootCtx)

	// Assumption: this will give unique IDs for the whole network
	// TODO: Can we generate equally distributed IDs based of IPs?
	id := kademliaid.NewRandomKademliaID()
	c := contact.NewContact(id, address)

	routingTable := routingtable.NewRoutingTable(c)

	memoryStore := store.NewMemoryStore()

	client := client.NewClient()

	node := node.New(client, routingTable, memoryStore, &network.PublicNetwork{})

	log.Println("STARTING SERVER")

	server := server.NewServer(routingTable, memoryStore)
	errGroup.Go(func() error {
		return server.Start(errCtx)
	})

	log.Println("STARTING BOOTSTRAP")
	errGroup.Go(func() error {
		time.Sleep(5 * time.Second)

		err := node.Bootstrap(errCtx)
		log.Printf("bootstrap err: %v\n", err)
		return err
	})

	if err = errGroup.Wait(); err != nil {
		log.Fatalf("bootstrap failed, %v", err)
	}

	log.Println("STARTING API")

	// REST API
	handler := api.NewHandler(node)
	errGroup.Go(func() error {
		return handler.ListenAndServe(errCtx)
	})

	log.Println("STARTING CLI")

	// CLI loop
	errGroup.Go(func() error {
		return cli.InputLoop(errCtx, cancelCtx, os.Stdin, node)
	})

	err = errGroup.Wait()
	cancelCtx()
	log.Printf("node shutdown, reason: %v", err)
}
