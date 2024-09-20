package node

import (
	"context"
	"d7024e_group04/env"
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
func (n *Node) GetObject(rootCtx context.Context, hash string) (data string, err error) {
	me := n.RoutingTable.Me()
	// check local store, if it was not found do clientRPC call
	if value, err := n.Store.GetValue(hash); err == nil {
		return value, nil
	}

	ctx, cancelCtx := context.WithCancel(rootCtx)
	ch := make(chan struct{}, env.Alpha)
	dataChan := make(chan string)
	kademliaHash := kademliaid.NewKademliaID(hash)

	shortlist := n.RoutingTable.FindClosestContacts(kademliaHash, me.ID, env.BucketSize)

	for {
		select {
		case ch <- struct{}{}:
			go func() {
				candidates, data, err := n.Client.SendFindValue(ctx, shortlist[0], me, hash)
				if err != nil {
					// TODO mark as unavailable, check if ctx canceled?
				}

				if data != "" {
					dataChan <- data

				} else {
					// TODO update shortlist
					panic(candidates)
				}
			}()

		case <-rootCtx.Done():
			cancelCtx()
			return "", rootCtx.Err()

		case data := <-dataChan:
			// cancel ctx as soon as value is found
			cancelCtx()

			// TODO store value in closest node if it did not have the value

			return data, nil
		}
	}

}
