package node

import (
	"context"
	"sync"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
)

func (n *Node) findNode(ctx context.Context, target *contact.Contact) ([]contact.Contact, error) {
	alpha := env.Alpha
	k := env.BucketSize
	visitedSet := contact.NewContactSet()
	me := n.RoutingTable.Me()
	var kClosets []contact.Contact
	updated := false

	responseContactChannel := make(chan []contact.Contact, k*alpha)

	kClosets = n.RoutingTable.FindClosestContacts(target.ID, k)

	for {
		updated = false
		var candidates contact.ContactCandidates
		for _, closeContact := range kClosets {
			if !visitedSet.Has(closeContact) {
				contactSlice := []contact.Contact{closeContact}
				candidates.Append(contactSlice)
			}
		}
		candidates.Sort()
	}

	// for list.HasBeenModified() {
	// 	list.ResetModifiedFlag()
	// 	var notContactedList []contact.Contact
	// 	for _, contact := range list.GetClosest() {
	// 		if !visitedSet.Has(contact) {
	// 			notContactedList = append(notContactedList, contact)
	// 		}
	// 	}
	//
	// 	j := 0
	// 	for !list.HasBeenModified() && j < len(notContactedList) {
	// 		var wg sync.WaitGroup
	// 		wg.Add(alpha)
	// 		for i := j; i < alpha+j && i < len(notContactedList); i++ {
	// 			contact := notContactedList[i]
	// 			go func() {
	// 				defer wg.Done()
	// 				visitedSet.Add(contact)
	// 				nodes, err := n.Client.SendFindNode(ctx, &me, &contact)
	// 				if err != nil {
	// 					return
	// 				}
	//
	// 				list.AddNodes(nodes, target)
	// 			}()
	// 		}
	//
	// 		wg.Wait()
	// 		j += alpha
	// 	}
	// }

	return list.GetClosest(), nil
}
