package kademlia

import (
	"context"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	"d7024e_group04/kademlia/network"
	"d7024e_group04/kademlia/routingtable"

	"golang.org/x/sync/errgroup"
)

type Node struct {
	store  map[string][]byte
	Client *network.Client
	server *network.Server
}

func NewNode(address string) *Node {
	id := kademliaid.NewRandomKademliaID()
	c := contact.NewContact(id, address)

	routingTable := routingtable.NewRoutingTable(c)
	server := network.NewServer(address, id, routingTable)

	client := network.NewClient(address, id, routingTable)

	return &Node{
		store:  make(map[string][]byte),
		Client: client,
		server: server,
	}
}

func (n *Node) Start(ctx context.Context) error {
	// TODO eventually start more stuff here that should be running during runtime
	errGroup, errCtx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return n.server.Start(errCtx)
	})

	return errGroup.Wait()
}
