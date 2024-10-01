package contact

import (
	"d7024e_group04/internal/kademlia/kademliaid"
	"fmt"
	"slices"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
// 3-tuple mentioned in the paper
// TODO: It is not like the paper <ip, udp port, node id>, but <id, address(ip+port), distance(maybe cached)>
type Contact struct {
	ID       kademliaid.KademliaID
	Address  string
	Distance kademliaid.KademliaID
}

// NewContact returns a new instance of a Contact
func NewContact(id kademliaid.KademliaID, address string) *Contact {
	return &Contact{id, address, kademliaid.KademliaID{}}
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
// TODO: target should be a contact
func (contact *Contact) CalcDistance(target kademliaid.KademliaID) {
	contact.Distance = contact.ID.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
// TODO: Don't know if we use it in this way, lets see
// We just compare distances here
// TODO: It implements the comparable interface or smth like this, for sorting
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.Distance.Less(otherContact.Distance)
}

// String returns a simple string representation of a Contact
// TODO: It should implement the standard toString interface, check this
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

// SortContacts sorts a slice of contacts
func SortContacts(contacts []*Contact) {
	slices.SortStableFunc(contacts, func(a, b *Contact) int {
		if a.Less(b) {
			return -1
		} else {
			return 1
		}
	})
}

// RemoveID filters a slice of contact and removes any instance of a contact with matching id.
func RemoveID(contacts []*Contact, id kademliaid.KademliaID) (contactsWithoutId []*Contact) {
	for idx, contact := range contacts {
		if contact.ID.Equals(id) {
			return append(contacts[:idx], contacts[idx+1:]...)
		}
	}
	return contacts
}
