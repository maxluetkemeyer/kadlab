package node

import (
	"context"
	"log"
	"slices"
	"sync"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
)

type kClosestList struct {
	mut     sync.RWMutex
	list    []contact.Contact
	updated bool
}

func (n *Node) findNode(rootCtx context.Context, target *contact.Contact) []contact.Contact {
	alpha := env.Alpha
	k := env.BucketSize
	visitedSet := contact.NewContactSet()
	me := n.RoutingTable.Me()
	kClosets := kClosestList{}

	wg := new(sync.WaitGroup)

	responseContactChannel := make(chan []contact.Contact, k*alpha) // TODO look at size

	kClosets.list = n.RoutingTable.FindClosestContacts(target.ID, k)

	for {
		kClosets.updated = false

		// get kClosest that are unvisited
		// TODO fix, tried to use contact.candidates but did not work. find some nicer way
		var candidates []contact.Contact
		for _, closeContact := range kClosets.list {
			if !visitedSet.Has(closeContact) {
				candidates = append(candidates, closeContact)
			}
		}
		// candidates should already be sorted

		// Goroutines, strict parallelism
		ctx, cancel := context.WithTimeout(rootCtx, env.RPCTimeout)
		for i := 0; i < alpha && i < len(candidates); i++ {
			wg.Add(1)
			visitedSet.Add(candidates[i])

			go func() {
				defer wg.Done()
				contacts, err := n.Client.SendFindNode(ctx, &me, &candidates[i])

				if err != nil {
					kClosets.remove(candidates[i])
					// TODO fix logging solution
					log.Printf("WARNING: findNode error: %v", err)
					return
				}

				responseContactChannel <- contacts
			}()
		}

		wg.Wait()
		cancel() //TODO maybe move

		for contacts := range responseContactChannel {
			for _, contact := range contacts {
				contact.CalcDistance(target.ID)

				// TODO refactor ifs
				if len(kClosets.list) < k {
					kClosets.list = append(kClosets.list, contact)
					kClosets.sort()
					kClosets.updated = true
				} else {
					if contact.Less(&kClosets.list[k-1]) {
						kClosets.list[k-1] = contact
						kClosets.sort()
						kClosets.updated = true
					}
				}
			}
		}

		// can we terminate?
		if !kClosets.updated && kClosets.isSubset(visitedSet) {
			return kClosets.list
		}
	}
}

func (kClosestList *kClosestList) isSubset(set *contact.ContactSet) bool {
	for _, contact := range kClosestList.list {
		if !set.Has(contact) {
			return false
		}
	}
	return true
}

func (kClosestList *kClosestList) sort() {
	kClosestList.mut.Lock()
	defer kClosestList.mut.Unlock()
	slices.SortStableFunc(kClosestList.list, func(a, b contact.Contact) int {
		if a.Less(&b) {
			return -1
		} else {
			return 1
		}
	})
}

// TODO make good
func (kClosestList *kClosestList) remove(target contact.Contact) {
	kClosestList.mut.Lock()
	defer kClosestList.mut.Unlock()

	var contactList []contact.Contact

	for _, contact := range kClosestList.list {
		if !contact.ID.Equals(target.ID) {
			contactList = append(contactList, contact)
		}
	}

	kClosestList.list = contactList
}
