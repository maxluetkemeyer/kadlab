package kademlia

import (
	"context"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	"d7024e_group04/kademlia/network"
	"d7024e_group04/kademlia/routingtable"
)

type Node struct {
	id           *kademliaid.KademliaID
	address      string
	store        map[string][]byte
	routingTable *routingtable.RoutingTable
	server       *network.Server
}

func NewNode(address string) *Node {
	id := kademliaid.NewRandomKademliaID()
	c := contact.NewContact(id, address)

	routingTable := routingtable.NewRoutingTable(c)
	server := network.NewServer(address, id, routingTable)

	return &Node{
		id:           id,
		address:      address,
		store:        make(map[string][]byte),
		routingTable: routingTable,
		server:       server,
	}
}

func (n *Node) Start(ctx context.Context) error {
	// TODO eventually start more stuff here that should be running during runtime
	return n.server.Start(ctx)
}
