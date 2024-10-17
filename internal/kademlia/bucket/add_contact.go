package bucket

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/network"
)

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed

// TODO: Split up into multiple check functions and test them isolated
func (bucket *Bucket) AddContact(newContact *contact.Contact, client network.ClientRPC) {
	bucket.Refresh()

	// Paper: "If the sending node already exists in the recipient’s k-bucket,
	// the recipient moves it to the tail of the list."
	for listContact := bucket.list.Front(); listContact != nil; listContact = listContact.Next() {
		listContactId := listContact.Value.(*contact.Contact).ID

		if newContact.ID.Equals(listContactId) {
			bucket.list.MoveToBack(listContact)
			return
		}
	}

	// Is there space in the bucket?
	// Paper: "If the node is not already in the appropriate k-bucket and the bucket has fewer than k entries,
	// then the recipient just inserts the new sender at the tail of the list.
	if bucket.list.Len() < bucket.size {
		bucket.list.PushBack(newContact)
		return
	}

	// Ping least-recently (head)
	head := bucket.list.Front()

	ctx, cancel := context.WithTimeout(context.Background(), env.RPCTimeout)
	_, err := client.SendPing(ctx, head.Value.(*contact.Contact).Address)
	cancel() // prevent context leaks

	if err == nil {
		bucket.list.MoveToBack(head)
		// discard new contact
		return
	}

	// Failed to respond
	bucket.list.Remove(head)
	bucket.list.PushBack(newContact)
}
