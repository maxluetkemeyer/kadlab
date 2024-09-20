package node

import (
	"context"
	"fmt"
	"log"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/network"
	"d7024e_group04/internal/store"
)

type Node struct {
	Client       network.ClientRPC
	RoutingTable *routingtable.RoutingTable
	Store        store.Store
}

func New(client network.ClientRPC, routingTable *routingtable.RoutingTable, store store.Store) *Node {
	return &Node{
		Client:       client,
		RoutingTable: routingTable,
		Store:        store,
	}
}

func (n *Node) Bootstrap(ctx context.Context) error {
	addresses, err := network.ResolveDNS(ctx, env.KnownDomain)
	if err != nil {
		return err;
	}

	me := n.RoutingTable.Me()
	contact, err := n.pingContacts(ctx, &me, addresses)
	if err != nil {
		return err
	}

	n.RoutingTable.AddContact(*contact)
	// TODO: iterative search
	n.Client.SendFindNode(ctx, &me)

	return nil
}

// Put takes content of the file and outputs the hash of the object
func (n *Node) PutObject() {
	panic("TODO")
}

// Get takes hash and outputs the contents of the object and the node it was retrieved
func (n *Node) GetObject() {
	// check store, if not found then do clientRPC call
	panic("TODO")
}



// Ping each contact in <contacts> until one responeses and returns it.
func (n *Node) pingContacts(ctx context.Context, me *contact.Contact, targets []string) (*contact.Contact, error) {
	for _, target := range targets {
		contact, err := n.Client.SendPing(ctx, me, target)
		if err == nil {
			return &contact, nil
		}
	}

	return nil, fmt.Errorf("unable to ping any contacts")
}
