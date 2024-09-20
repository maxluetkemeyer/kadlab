package contact

import "sync"

// A mathematical set of contacts with a mutex for thread safety
type ContactSet struct {
	m map[Contact]bool
	sync.RWMutex
}

func NewContactSet() *ContactSet {
	return &ContactSet{
		m: make(map[Contact]bool),
	}
}

// func main() {
// 	// Initialize our ContactSet
// 	s := New()
//
// 	// Add example items
// 	s.Add("item1")
// 	s.Add("item1") // duplicate item
// 	s.Add("item2")
// 	fmt.Printf("%d items\n", s.Len())
//
// 	// Clear all items
// 	s.Clear()
// 	if s.IsEmpty() {
// 		fmt.Printf("0 items\n")
// 	}
//
// 	s.Add("item2")
// 	s.Add("item3")
// 	s.Add("item4")
//
// 	// Check for existence
// 	if s.Has("item2") {
// 		fmt.Println("item2 does exist")
// 	}
//
// 	// Remove some of our items
// 	s.Remove("item2")
// 	s.Remove("item4")
// 	fmt.Println("list of all items:", s.List())
// }

// Add add a contact to the set
func (s *ContactSet) Add(item Contact) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

// Add add a slice of contacts to the set
func (s *ContactSet) Adds(items []Contact) {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		s.m[item] = true
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

func (s *ContactSet) GetClosests(n int) []Contact {

	panic("unim")
}

// Len returns the number of contacts in the set.
func (s *ContactSet) Len() int {
	return len(s.List())
}

// Clear removes all contacts from the set
func (s *ContactSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[Contact]bool)
}

// IsEmpty checks for emptiness
func (s *ContactSet) IsEmpty() bool {
	return s.Len() == 0
}

// ContactSet returns a slice of all items
func (s *ContactSet) List() []Contact {
	s.RLock()
	defer s.RUnlock()
	list := make([]Contact, 0)
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
