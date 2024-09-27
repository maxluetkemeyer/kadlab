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
	kNet         network.Network
}

func New(client network.ClientRPC, routingTable *routingtable.RoutingTable, store store.Store, kNet network.Network) *Node {
	return &Node{
		Client:       client,
		RoutingTable: routingTable,
		Store:        store,
		kNet:         kNet,
	}
}

/*
- A node joins the network as follows:
 1. if it does not already have a nodeID n, it generates one
 2. it inserts the value of some known node c into the appropriate bucket as its first contact
 3. it does an iterativeFindNode for n
 4. it refreshes all buckets further away than its closest neighbor, which will be in the occupied bucket with the lowest index.
*/
func (n *Node) Bootstrap(ctx context.Context) error {
	addresses, err := n.kNet.ResolveDNS(env.NodesProxyDomain)
	if err != nil {
		return err
	}

	ourIP := n.RoutingTable.Me().Address

	for idx, address := range addresses {
		if address == ourIP[:len(ourIP)-6] {
			addresses = append(addresses[:idx], addresses[idx+1:]...)
			break
		}
	}

	log.Printf("addresses:%v\n", addresses)

	contact, err := n.pingIPsAndGetContact(ctx, addresses)
	if err != nil {
		return err
	}

	n.RoutingTable.AddContact(contact)
	// TODO: iterative search, should we update once for the list or for each node visited in findNode?
	me := n.RoutingTable.Me()
	closestContacts := n.findNode(ctx, me)
	for _, contact := range closestContacts {
		n.RoutingTable.AddContact(contact)
	}

	log.Println("FOUND NODES")
	log.Print(closestContacts)

	return nil
}

// Put takes content of the file and outputs the hash of the object
func (n *Node) PutObject() {
	panic("TODO")
}

// Get takes hash and outputs the contents of the object and the node it was retrieved
func (n *Node) GetObject(rootCtx context.Context, hash string) (data string, err error) {
	panic("TODO")
}

// Ping each contact in <contacts> until one responeses and returns it.
func (n *Node) pingIPsAndGetContact(ctx context.Context, targetIPs []string) (*contact.Contact, error) {
	for _, targetIP := range targetIPs {
		contact, err := n.Client.SendPing(ctx, fmt.Sprintf("%v:%v", targetIP, env.Port))
		if err == nil {
			return contact, nil
		}
	}

	return nil, fmt.Errorf("unable to ping any contacts")
}
