package bucket

import (
	"d7024e_group04/kademlia"
	"d7024e_group04/kademlia/contact"
)

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *Bucket) AddContact(newContact contact.Contact) {
	bucket.addContactCustom(newContact, kademlia.BucketSize)
}

// TODO: Split up into multiple check functions and test them isolated
func (bucket *Bucket) addContactCustom(newContact contact.Contact, bucketSize int) {

	// Is the new contact already stored in the list?
	// Paper: "If the sending node already exists in the recipient’s k-bucket,
	// the recipient moves it to the tail of the list."
	for listContact := bucket.list.Front(); listContact != nil; listContact = listContact.Next() {
		listContactId := listContact.Value.(contact.Contact).ID

		if newContact.ID.Equals(listContactId) {
			bucket.list.MoveToBack(listContact)
			return
		}
	}

	// Is there space in the bucket?
	// Paper: "If the node is not already in the appropriate k-bucket and the bucket has fewer than k entries,
	// then the recipient just inserts the new sender at the tail of the list.
	if bucket.list.Len() < bucketSize {
		bucket.list.PushBack(newContact)
		return
	}

	// Ping least-recently (head)
	head := bucket.list.Front()
	// TODO: ping

	// Responded
	if true {
		// TODO: Check if list is not empty, should not be possible with a BucketSize > 0
		bucket.list.MoveToBack(head)
		// discard new contact
		return
	}

	// Failed to respond
	bucket.list.Remove(head)
	bucket.list.PushBack(newContact)
}
