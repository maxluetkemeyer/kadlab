package node

import (
	"container/list"
	"context"
	"fmt"
	"sync"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/bucket"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/network"
	"d7024e_group04/internal/store"
)

type Node struct {
	Client       network.ClientRPC
	RoutingTable *routingtable.RoutingTable
	Store        store.Store
	kNet		 network.Network
}

func New(client network.ClientRPC, routingTable *routingtable.RoutingTable, store store.Store, kNet network.Network) *Node {
	return &Node{
		Client:       client,
		RoutingTable: routingTable,
		Store:        store,
		kNet:		  kNet,
	}
}

/*
	* A node joins the network as follows:
      1. if it does not already have a nodeID n, it generates one
      2. it inserts the value of some known node c into the appropriate bucket as its first contact
      3. it does an iterativeFindNode for n
      4. it refreshes all buckets further away than its closest neighbor, which will be in the occupied bucket with the lowest index.
*/
func (n *Node) Bootstrap(ctx context.Context) error {
	addresses, err := n.kNet.ResolveDNS(ctx, env.KnownDomain)
	if err != nil {
		return err
	}

	me := n.RoutingTable.Me()
	contact, err := n.pingContacts(ctx, &me, addresses)
	if err != nil {
		return err
	}

	n.RoutingTable.AddContact(*contact)
	// TODO: iterative search
	// neighbors, err := n.Client.SendFindNode(ctx, &me)
	// if err != nil {
	// 	return err
	// }

	// 4. Automatically done via AddContact()

	panic("unimplemented")
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
// TODO: add concurrency
func (n *Node) pingContacts(ctx context.Context, me *contact.Contact, targetIps []string) (*contact.Contact, error) {
	for _, targetIp := range targetIps {
		contact, err := n.Client.SendPing(ctx, me, targetIp)
		if err == nil {
			return &contact, nil
		}
	}

	return nil, fmt.Errorf("unable to ping any contacts")
}
