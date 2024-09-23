package node

import (
	"context"
	"fmt"
	"log"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
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
	addresses, err := n.kNet.ResolveDNS(ctx, env.NodesProxyDomain)
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

	me := n.RoutingTable.Me()
	contact, err := n.pingContacts(ctx, &me, addresses)
	if err != nil {
		return err
	}

	n.RoutingTable.AddContact(*contact)
	// TODO: iterative search, should we update once for the list or for each node visited in findNode?
	closestContacts := n.findNode(ctx, &me)
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
	me := n.RoutingTable.Me()
	// check local store, if it was not found do clientRPC call
	if value, err := n.Store.GetValue(hash); err == nil {
		return value, nil
	}

	ctx, cancelCtx := context.WithCancel(rootCtx)
	ch := make(chan struct{}, env.Alpha)
	dataChan := make(chan string)
	kademliaHash := kademliaid.NewKademliaID(hash)

	shortlist := n.RoutingTable.FindClosestContacts(kademliaHash, env.BucketSize, me.ID)

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

// Ping each contact in <contacts> until one responeses and returns it.
// TODO: add concurrency
func (n *Node) pingContacts(ctx context.Context, me *contact.Contact, targetIps []string) (*contact.Contact, error) {
	for _, targetIp := range targetIps {
		contact, err := n.Client.SendPing(ctx, me, targetIp+":50051") // TODO fix this port
		if err == nil {
			return contact, nil
		}
	}

	return nil, fmt.Errorf("unable to ping any contacts")
}
