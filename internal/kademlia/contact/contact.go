package contact

import (
	"d7024e_group04/internal/kademlia/kademliaid"
	"fmt"
	"sort"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
// 3-tuple mentioned in the paper
// TODO: It is not like the paper <ip, udp port, node id>, but <id, address(ip+port), distance(maybe cached)>
type Contact struct {
	ID       *kademliaid.KademliaID
	Address  string
	distance *kademliaid.KademliaID
}

// NewContact returns a new instance of a Contact
func NewContact(id *kademliaid.KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *kademliaid.KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
// TODO: Don't know if we use it in this way, lets see
// We just compare distances here
// TODO: It implements the comparable interface or smth like this, for sorting
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

// String returns a simple string representation of a Contact
// TODO: It should implement the standard toString interface, check this
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

// ContactCandidates definition
// stores a slice of Contacts
// TODO: Don't know why they implemented that, we can just use a slice of contacts
type ContactCandidates struct {
	contacts []Contact
}

// Append a slice of Contacts to the ContactCandidates
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

// GetContacts returns the first 'count' number of Contacts
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.contacts[:count]
}

// Sort the Contacts in ContactCandidates
// TODO: Maybe use slices.sortfunc
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the ContactCandidates
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

// Swap the position of the Contacts at i and j
// TODO: WARNING does not check if either i or j is within range
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}

func (candidates *ContactCandidates) RemoveID(id *kademliaid.KademliaID) {
	for idx, contact := range candidates.contacts {
		if contact.ID == id {
			candidates.contacts = append(candidates.contacts[:idx], candidates.contacts[idx+1:]...)
		}
	}
}
