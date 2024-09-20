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

func (n *Node) findNode(ctx context.Context, target *contact.Contact, k int) ([]contact.Contact, error) {
	contactedSet := contact.NewContactSet()

	closestContacts := n.RoutingTable.FindClosestContacts(target.ID, k)

	list := NewNodeList(k)
	list.AddNodes(closestContacts, target)
	for list.HasBeenModified() {
		list.ResetModifiedFlag()
		var notContactedList []contact.Contact
		for _, contact := range list.GetClosest() {
			if !contactedSet.Has(contact) {
				notContactedList = append(notContactedList, contact)
			}
		}

		j := 0
		for !list.HasBeenModified() && j < len(notContactedList) {
			var wg sync.WaitGroup
			wg.Add(env.Alpha)
			for i := j; i < env.Alpha+j && i < len(notContactedList); i++ {
				contact := notContactedList[i]
				go func() {
					defer wg.Done()
					contactedSet.Add(contact)
					nodes, err := n.Client.SendFindNode(ctx, &contact)
					if err != nil {
						return
					}

					list.AddNodes(nodes, target)
				}()
			}

			wg.Wait()
			j += env.Alpha
		}
	}

	return list.GetClosest(), nil
}
