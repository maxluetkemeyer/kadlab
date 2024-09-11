package bucket

import (
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
)

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *Bucket) GetContactAndCalcDistance(targetId *kademliaid.KademliaID) []contact.Contact {
	var contactsWithDistances []contact.Contact

	for contactElement := bucket.list.Front(); contactElement != nil; contactElement = contactElement.Next() {
		contactFromList := contactElement.Value.(contact.Contact) // use generics with list
		contactFromList.CalcDistance(targetId)

		contactsWithDistances = append(contactsWithDistances, contactFromList) // slices are immutable
	}

	return contactsWithDistances
}
