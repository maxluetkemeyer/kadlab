package node

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/model"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/network"
	"d7024e_group04/internal/store"
)

type NodeHandler interface {
	Me() *contact.Contact
	Bootstrap(rootCtx context.Context) error
	PutObject(ctx context.Context, data string) (hashAsHex string, err error)
	GetObject(rootCtx context.Context, hash string) (valueObject *model.ValueObject, candidates []*contact.Contact, err error)
}

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

func (n *Node) Me() *contact.Contact {
	return n.RoutingTable.Me()
}

/*
- A node joins the network as follows:
 1. if it does not already have a nodeID n, it generates one
 2. it inserts the value of some known node c into the appropriate bucket as its first contact
 3. it does an iterativeFindNode for n
 4. it refreshes all buckets further away than its closest neighbor, which will be in the occupied bucket with the lowest index.
*/
func (n *Node) Bootstrap(rootCtx context.Context) error {
	var recipient *contact.Contact

	logger := slog.Default().With("me", n.RoutingTable.Me().Address)
	ctx, cancel := context.WithTimeout(rootCtx, env.BootstrapTimeout)
	for {
		if ctx.Err() != nil {
			defer cancel()
			return ctx.Err()
		}

		ips, err := n.kNet.ResolveDNS(env.NodesProxyDomain)
		if err != nil {
			logger.Warn("Unable to resolve DNS", slog.Any("domain", env.NodesProxyDomain))
			time.Sleep(100 * time.Millisecond)
			continue
		}

		ourAddress := n.RoutingTable.Me().Address
		ips = removeAddress(ips, ourAddress)

		logger := logger.With(slog.Any("ips", ips))

		logger.Debug("pinging")

		recipient, err = n.pingIPsAndGetContact(ctx, ips)
		if err == nil {
			break
		}

		logger.Warn("Unable to ping any ip", slog.Any("err", err))
		slog.Warn("Retrying bootstrap in 100 milis")
		time.Sleep(100 * time.Millisecond)
	}

	cancel()

	n.RoutingTable.AddContact(recipient)
	// TODO: iterative search, should we update once for the list or for each node visited in findNode?
	me := n.RoutingTable.Me()
	closestContacts := n.findNode(rootCtx, me)
	for _, contact := range closestContacts {
		n.RoutingTable.AddContact(contact)
	}

	logger.Debug("FOUND NODES", slog.Any("nodes", closestContacts))

	return nil
}

// Put takes content of the file and outputs the hash of the object
// This is the Kademlia store operation. The initiating node does an iterativeFindNode, collecting a set of k closest contacts, and then sends a primitive STORE RPC to each.
func (n *Node) PutObject(ctx context.Context, data string) (hashAsHex string, err error) {
	hash := kademliaid.NewKademliaIDFromData(data)
	storeContact := contact.NewContact(hash, "")

	storeCandidates := n.findNode(ctx, storeContact)

	wg := new(sync.WaitGroup)
	errChan := make(chan error, env.BucketSize)

	for _, contact := range storeCandidates {
		wg.Add(1)

		go func() {
			defer wg.Done()
			err := n.Client.SendStore(ctx, contact, data)
			if err != nil {
				errChan <- err
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		slog.Warn("SendStore error", slog.Any("err", err))
	}

	return hash.String(), nil
}

// Get takes hash and outputs the contents of the object and the node it was retrieved
func (n *Node) GetObject(rootCtx context.Context, hash string) (valueObject *model.ValueObject, candidates []*contact.Contact, err error) {
	// check for value in our own store first
	value, err := n.Store.GetValue(hash)

	if err == nil {
		return &model.ValueObject{DataValue: value, NodeWithValue: n.RoutingTable.Me()}, nil, nil
	}

	// search the network
	ctx, cancel := context.WithCancel(rootCtx)

	valueChan := make(chan model.ValueObject, alpha)
	candidateChan := make(chan []*contact.Contact, 1)

	hashAsKademliaID := kademliaid.NewKademliaID(hash)
	hashAsContact := contact.NewContact(hashAsKademliaID, "")

	visitedSet := contact.NewContactSet()
	kClosest := &kClosestList{}
	kClosest.list = n.RoutingTable.FindClosestContacts(hashAsKademliaID, k)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return

			default:
				responseContactChan := make(chan []*contact.Contact, alpha)
				kClosest.updated = false
				rpcCtx, cancel := context.WithTimeout(ctx, env.RPCTimeout)

				// blocking call
				n.runParallelFindValueRequest(rpcCtx, kClosest, visitedSet, hash, responseContactChan, valueChan)
				cancel()

				for contacts := range responseContactChan {
					for _, contact := range contacts {
						// don't add nodes that are already visited
						if !visitedSet.Has(contact) {
							kClosest.addContact(contact, hashAsContact)
						}
					}
				}

				if !kClosest.updated && kClosest.isSubset(visitedSet) {
					candidateChan <- kClosest.list
					return
				}
			}
		}
	}()

	for {
		select {
		case <-rootCtx.Done():
			cancel()
			return nil, nil, rootCtx.Err()
		case candidates := <-candidateChan:
			cancel()
			return nil, candidates, nil
		case valueObject := <-valueChan:
			cancel()
			wg.Wait()
			n.storeAtClosestNode(rootCtx, kClosest, visitedSet, hash, valueObject.DataValue)
			return &valueObject, nil, nil
		}
	}

}

// storeAtClosestNode stores the data in the closest node seen which does not have the data
func (n *Node) storeAtClosestNode(rootCtx context.Context, kClosest *kClosestList, visitedSet *contact.ContactSet, hash, data string) {
	candidates := kClosest.List()
	for _, contact := range candidates {
		if visitedSet.Has(contact) {
			ctx, cancel := context.WithTimeout(rootCtx, env.RPCTimeout)
			candidates, _, err := n.Client.SendFindValue(ctx, contact, hash)
			cancel()

			// data exists or node is down
			if candidates == nil || err != nil {
				continue
			}

			ctx, cancel = context.WithTimeout(rootCtx, env.RPCTimeout)
			err = n.Client.SendStore(ctx, contact, data)
			cancel()

			if err == nil {
				return
			}

		}
	}
}

func (n *Node) runParallelFindValueRequest(
	ctx context.Context,
	kClosest *kClosestList,
	visitedSet *contact.ContactSet,
	hash string,
	responseContactChannel chan []*contact.Contact,
	valueChan chan model.ValueObject) {

	wg := new(sync.WaitGroup)
	candidates := getCandidates(kClosest, visitedSet)

	for i := 0; i < alpha && i < len(candidates); i++ {
		wg.Add(1)
		visitedSet.Add(candidates[i])

		go func() {
			defer wg.Done()
			contacts, data, err := n.Client.SendFindValue(ctx, candidates[i], hash)

			if err != nil {
				kClosest.remove(candidates[i])
				log.Printf("WARNING: client findValue error: %v", err)
				return
			}

			if len(data) > 0 {
				valueChan <- model.ValueObject{DataValue: data, NodeWithValue: candidates[i]}
				return
			}

			responseContactChannel <- contacts
		}()
	}
	wg.Wait()
	close(responseContactChannel)
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

func removeAddress(ips []string, addressToRemove string) []string {
	portLength := 5 + /* colon */ 1
	ipToRemove := addressToRemove[:len(addressToRemove)-portLength]

	for idx, ip := range ips {
		if ip == ipToRemove {
			// Remove the ip
			ips = append(ips[:idx], ips[idx+1:]...)
			break
		}
	}

	return ips
}
