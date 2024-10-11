package contact

import (
	"d7024e_group04/internal/kademlia/kademliaid"
	"sync"
)

// A mathematical set of contacts with a mutex for thread safety
type KademliaIdSet struct {
	m map[kademliaid.KademliaID]bool
	sync.RWMutex
}

func NewKademliaIdSet() *KademliaIdSet {
	return &KademliaIdSet{
		m: make(map[kademliaid.KademliaID]bool),
	}
}

// Add adds a contact to the set
func (s *KademliaIdSet) Add(item *Contact) {
	s.Lock()
	defer s.Unlock()
	s.m[item.ID] = true
}

// Add adds a slice of contacts to the set
func (s *KademliaIdSet) Adds(items []*Contact) {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		s.m[item.ID] = true
	}
}

// Remove deletes the contact from the set
func (s *KademliaIdSet) Remove(item *Contact) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item.ID)
}

// Has looks for the existence of the contact
func (s *KademliaIdSet) Has(item *Contact) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item.ID]
	return ok
}

// Len returns the number of contacts in the set.
func (s *KademliaIdSet) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.m)
}

// Clear removes all contacts from the set
func (s *KademliaIdSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[kademliaid.KademliaID]bool)
}

// IsEmpty checks for emptiness
func (s *KademliaIdSet) IsEmpty() bool {
	return s.Len() == 0
}
