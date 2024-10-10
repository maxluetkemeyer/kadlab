package node

import (
	"context"
	"log"
	"sync"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
)

// TODO maybe reconsider how we init these for tests. this is assigned during init time
var (
	alpha = env.Alpha
	k     = env.BucketSize
)

func (n *Node) findNode(rootCtx context.Context, contactWeAreSearchingFor *contact.Contact) []*contact.Contact {
	visitedSet := contact.NewContactSet()
	kClosest := &kClosestList{}

	// TODO these are reassigned during runtime
	alpha = env.Alpha
	k = env.BucketSize

	kClosest.list = n.RoutingTable.FindClosestContacts(contactWeAreSearchingFor.ID, k)

	for {
		kClosest.updated = false

		n.findNodeIteration(rootCtx, contactWeAreSearchingFor, visitedSet, kClosest)

		// can we terminate?
		if !kClosest.updated && kClosest.isSubset(visitedSet) {
			return kClosest.list
		}

	}
}

func (n *Node) findNodeIteration(
	rootCtx context.Context,
	contactWeAreSearchingFor *contact.Contact,
	visitedSet *contact.ContactSet,
	kClosest *kClosestList) {

	// get kClosest that are unvisited
	// candidates should already be sorted
	candidates := getCandidates(kClosest, visitedSet)

	// run findNode in strict parallelism and wait for responses
	ctx, cancel := context.WithTimeout(rootCtx, env.RPCTimeout)
	wg := new(sync.WaitGroup)
	responseContactChannel := n.runParallelFindNodeRequest(ctx, candidates, wg, visitedSet, contactWeAreSearchingFor, kClosest)
	wg.Wait()
	close(responseContactChannel)
	cancel() // prevent context leaks, need to be called

	// loop through responses and add to kClosest
	for contacts := range responseContactChannel {
		for _, contact := range contacts {
			// don't add nodes that are already visited
			if !visitedSet.Has(contact) {
				kClosest.addContact(contact, contactWeAreSearchingFor)
			}
		}
	}
}

func (n *Node) runParallelFindNodeRequest(
	ctx context.Context,
	candidates []*contact.Contact,
	wg *sync.WaitGroup,
	visitedSet *contact.ContactSet,
	contactWeAreSearchingFor *contact.Contact,
	kClosest *kClosestList) chan []*contact.Contact {
	// Goroutines, strict parallelism
	responseContactChannel := make(chan []*contact.Contact, alpha)

	for i := 0; i < alpha && i < len(candidates); i++ {
		wg.Add(1)
		visitedSet.Add(candidates[i])

		go func() {
			defer wg.Done()
			contacts, err := n.Client.SendFindNode(ctx, candidates[i], contactWeAreSearchingFor)

			if err != nil {
				kClosest.remove(candidates[i])
				log.Printf("WARNING: client findNode error: %v", err)
				return
			}

			responseContactChannel <- contacts
		}()
	}
	return responseContactChannel
}

func getCandidates(kClosest *kClosestList, visitedSet *contact.ContactSet) (candidates []*contact.Contact) {
	for _, closeContact := range kClosest.list {
		if !visitedSet.Has(closeContact) {
			candidates = append(candidates, closeContact)
		}
	}
	return candidates
}
