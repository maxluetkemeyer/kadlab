package contact

import "sync"

// A mathematical set of contacts with a mutex for thread safety
type ContactSet struct {
	m map[Contact]struct{}
	sync.RWMutex
}

func NewContactSet() *ContactSet {
	return &ContactSet{
		m: make(map[Contact]struct{}),
	}
}

// Add adds a contact to the set
func (s *ContactSet) Add(item Contact) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = struct{}{}
}

// Add adds a slice of contacts to the set
func (s *ContactSet) Adds(items []Contact) {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		s.m[item] = struct{}{}
	}
}

// Remove deletes the contact from the set
func (s *ContactSet) Remove(item Contact) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

// Has looks for the existence of the contact
func (s *ContactSet) Has(item Contact) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

// Len returns the number of contacts in the set.
func (s *ContactSet) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.m)
}

// Clear removes all contacts from the set
func (s *ContactSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[Contact]struct{})
}

// IsEmpty checks for emptiness
func (s *ContactSet) IsEmpty() bool {
	return s.Len() == 0
}
