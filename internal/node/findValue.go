package node

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/model"
	"log"
	"sync"
)

// Get takes hash and outputs the contents of the object and the node it was retrieved
func (n *Node) GetObject(rootCtx context.Context, hash string) (FindValueSuccessfulResponse *model.FindValueSuccessfulResponse, candidates []*contact.Contact, err error) {
	// check for value in our own store first
	dataObject, err := n.Store.GetValue(hash)

	if err == nil {
		n.Store.SetTTL(hash, env.TTL)
		return &model.FindValueSuccessfulResponse{DataValue: dataObject.Data, NodeWithValue: n.RoutingTable.Me()}, nil, nil
	}

	// search the network
	ctx, cancel := context.WithCancel(rootCtx)

	valueChan := make(chan model.FindValueSuccessfulResponse, alpha)
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

				// if we don't find value at all, and have visited all the kClosest
				if !kClosest.updated && kClosest.isSubset(visitedSet) {
					candidateChan <- kClosest.list
					return
				}
			}
		}
	}()

	for {
		select {
		// main context is canceled
		case <-rootCtx.Done():
			cancel()
			return nil, nil, rootCtx.Err()
			// if we don't find the value, return kclosest
		case candidates := <-candidateChan:
			cancel()
			return nil, candidates, nil

			// if we find value, cancel rest of calls early
		case FindValueSuccessfulResponse := <-valueChan:
			cancel()
			wg.Wait()
			n.storeAtClosestNode(rootCtx, kClosest, visitedSet, hash, FindValueSuccessfulResponse)
			return &FindValueSuccessfulResponse, nil, nil
		}
	}

}

func (n *Node) runParallelFindValueRequest(
	ctx context.Context,
	kClosest *kClosestList,
	visitedSet *contact.ContactSet,
	hash string,
	responseContactChannel chan []*contact.Contact,
	valueChan chan model.FindValueSuccessfulResponse) {

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

			if len(data.DataValue) > 0 {
				valueChan <- model.FindValueSuccessfulResponse{DataValue: data.DataValue, NodeWithValue: candidates[i], OriginalUploader: data.OriginalUploader}
				return
			}

			responseContactChannel <- contacts
		}()
	}
	wg.Wait()
	close(responseContactChannel)
}

// storeAtClosestNode stores the data in the closest node seen which does not have the data
func (n *Node) storeAtClosestNode(rootCtx context.Context, kClosest *kClosestList, visitedSet *contact.ContactSet, hash string, dataObject model.FindValueSuccessfulResponse) {
	candidates := kClosest.List()
	for _, contact := range candidates {
		if !visitedSet.Has(contact) {
			continue
		}

		ctx, cancel := context.WithTimeout(rootCtx, env.RPCTimeout)
		_, dataFromNode, err := n.Client.SendFindValue(ctx, contact, hash)
		cancel()

		// data exists or node is down
		if len(dataFromNode.DataValue) > 0 || err != nil {
			continue
		}

		ctx, cancel = context.WithTimeout(rootCtx, env.RPCTimeout)
		err = n.Client.SendStore(ctx, contact, dataObject.DataValue, dataObject.OriginalUploader)
		cancel()

		if err == nil {
			ctx, cancel = context.WithTimeout(rootCtx, env.RPCTimeout)
			n.Client.SendNewStoredLocation(ctx, hash, dataObject.OriginalUploader, contact)
			cancel()
			return
		}

	}
}
