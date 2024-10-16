package store

import (
	"d7024e_group04/internal/kademlia/contact"
	"fmt"
	"sync"
	"time"
)

// SimpleTTLStore TODO: Thread safety!
type SimpleTTLStore struct {
	sync.RWMutex
	store               Store
	ttlOfKey            map[string]time.Time
	DataRefreshContacts map[string][]*contact.Contact
}

func NewSimpleTTLStore(store Store) *SimpleTTLStore {
	return &SimpleTTLStore{
		store:               store,
		ttlOfKey:            make(map[string]time.Time),
		DataRefreshContacts: make(map[string][]*contact.Contact),
	}
}

func (s *SimpleTTLStore) SetValue(key string, value string, ttl time.Duration, uploader *contact.Contact) {
	s.store.SetValue(key, value, uploader)
	s.SetTTL(key, ttl)
}

// Does not reset the TTL!
func (s *SimpleTTLStore) GetValue(key string) (value string, error error) {
	if !s.isValid(key) {
		return "", fmt.Errorf("invalid key (too old or not stored)")
	}

	return s.store.GetValue(key)
}

func (s *SimpleTTLStore) SetTTL(key string, ttl time.Duration) {
	s.ttlOfKey[key] = time.Now().Add(ttl)
}

func (s *SimpleTTLStore) GetTTL(key string) time.Duration {
	// If key is not present, the "zero" time will be returned
	endTime := s.ttlOfKey[key]

	// If the result exceeds the maximum (or minimum) value that can be stored in a Duration,
	// the maximum (or minimum) duration will be returned.
	remainingTime := time.Until(endTime)
	return remainingTime
}

func (s *SimpleTTLStore) GetOriginalUploader(key string) (*contact.Contact, error) {
	return s.store.GetOriginalUploader(key)
}

func (s *SimpleTTLStore) GetStoreLocations(key string) []*contact.Contact {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.DataRefreshContacts[key]
}

func (s *SimpleTTLStore) AddStoreLocation(key string, contactToAdd *contact.Contact) {
	s.RWMutex.Lock()
	if _, found := s.DataRefreshContacts[key]; !found {
		s.DataRefreshContacts[key] = []*contact.Contact{contactToAdd}
	} else {
		s.DataRefreshContacts[key] = append(s.DataRefreshContacts[key], contactToAdd)
	}

	s.RWMutex.Unlock()
}

func (s *SimpleTTLStore) isValid(key string) bool {
	remainingTime := s.GetTTL(key)

	return remainingTime > 0
}
