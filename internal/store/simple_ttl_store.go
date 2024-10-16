package store

import (
	"fmt"
	"time"
)

// SimpleTTLStore TODO: Thread safety!
type SimpleTTLStore struct {
	store    Store
	ttlOfKey map[string]time.Time
}

func NewSimpleTTLStore(store Store) *SimpleTTLStore {
	return &SimpleTTLStore{
		store:    store,
		ttlOfKey: make(map[string]time.Time),
	}
}

func (s *SimpleTTLStore) SetValue(key string, value string, ttl time.Duration) {
	s.store.SetValue(key, value)
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

func (s *SimpleTTLStore) isValid(key string) bool {
	remainingTime := s.GetTTL(key)

	return remainingTime > 0
}
