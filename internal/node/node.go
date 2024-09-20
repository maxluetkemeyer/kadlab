package node

import (
	"context"
	"d7024e_group04/internal/kademlia/kademliaid"
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

// Put takes content of the file and outputs the hash of the object
func (n *Node) PutObject() {
	panic("TODO")
}

// Get takes hash and outputs the contents of the object and the node it was retrieved
func (n *Node) GetObject(ctx context.Context, hash *kademliaid.KademliaID) (string, error) {
	// check local store, if it was not found do clientRPC call
	if value, err := n.Store.GetValue(hash.String()); err == nil {
		return value, nil
	}

	return n.Client.SendFindValue(ctx, hash)
}
