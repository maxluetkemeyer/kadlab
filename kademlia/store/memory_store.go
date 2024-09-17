package store

import (
	"fmt"
	"log"
	"sync"
)

type MemoryStore struct {
	mut   *sync.RWMutex
	myMap map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		mut:   new(sync.RWMutex),
		myMap: make(map[string][]byte),
	}
}

func (store *MemoryStore) SetValue(key string, value []byte) {
	log.Printf("stored value for key: %s", key)
	store.mut.Lock()
	defer store.mut.Unlock()

	store.myMap[key] = value
}

func (store *MemoryStore) GetValue(key string) (value []byte, err error) {
	store.mut.RLock()
	defer store.mut.RUnlock()

	value, found := store.myMap[key]

	if !found {
		err = fmt.Errorf("value not found for key: %s", key)
		return nil, err
	}

	return value, nil
}
