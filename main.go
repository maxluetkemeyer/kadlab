// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/kademlia"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	"d7024e_group04/kademlia/network"
	"d7024e_group04/kademlia/routingtable"
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

	id := kademliaid.NewRandomKademliaID()
	c := contact.NewContact(id, address)

	routingTable := routingtable.NewRoutingTable(c)
	server := network.NewServer(address, id, routingTable)

	client := network.NewClient(address, id, routingTable)

	node := kademlia.NewNode(address, client, server)

	err = node.Start(rootCtx)

	cancelCtx()
	log.Printf("node shutdown, reason: %v", err)
}
