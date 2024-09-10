package kademlia

import (
	"container/list"
)

//TODO: Insert guard clauses
//TODO: For loops with range
//TODO: Read paper
//TODO: Generic List

// bucket definition
// contains a List
type bucket struct {
	list *list.List
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	// nil pointer
	var element *list.Element

	// full list scan in worst case
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	// if element is not in the list
	if element == nil {
		// if bucket has space left
		if bucket.list.Len() < bucketSize {
			// shouldn't it be at the end?
			// it was: bucket.list.PushFront(contact)
			bucket.list.PushBack(contact)
		} else {
			//TODO: Read paper
		}
	} else {
		// it should be moved to the back
		// it was MoveToFront
		bucket.list.MoveToBack(element)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for element := bucket.list.Front(); element != nil; element = element.Next() {
		contact := element.Value.(Contact) // use generics with list
		contact.CalcDistance(target)
		contacts = append(contacts, contact) // slices are immutable
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
