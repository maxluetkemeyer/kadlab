package bucket

// least-recently seen (head) (old)
// most-recently seen (tail) (new)

import (
	"container/list"
	"d7024e_group04/kademlia"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
)

//TODO: Insert guard clauses
//TODO: For loops with range
//TODO: Read paper
//TODO: Generic List

// Bucket definition
// contains a List
// TODO: Need to be public for "routingtable.go", make this private
type Bucket struct {
	list *list.List
}

// NewBucket returns a new instance of a Bucket
func NewBucket() *Bucket {
	bucket := &Bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
// TODO: Split up into multiple check functions and test them isolated
func (bucket *Bucket) AddContact(newContact contact.Contact) {

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
	if bucket.list.Len() < kademlia.BucketSize {
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

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *Bucket) GetContactAndCalcDistance(target *kademliaid.KademliaID) []contact.Contact {
	var contacts []contact.Contact

	for element := bucket.list.Front(); element != nil; element = element.Next() {
		contact := element.Value.(contact.Contact) // use generics with list
		contact.CalcDistance(target)
		contacts = append(contacts, contact) // slices are immutable
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *Bucket) Len() int {
	return bucket.list.Len()
}
