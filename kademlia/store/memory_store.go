package store

import (
	"fmt"
	"log"
	"sync"
)

type MemoryStore struct {
	mut   *sync.RWMutex
	myMap map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		mut:   new(sync.RWMutex),
		myMap: make(map[string]string),
	}
}

func (store *MemoryStore) SetValue(key string, value string) {
	log.Printf("stored value for key: %s", key)
	store.mut.Lock()
	defer store.mut.Unlock()

	store.myMap[key] = value
}

func (store *MemoryStore) GetValue(key string) (value string, err error) {
	store.mut.RLock()
	defer store.mut.RUnlock()

	value, found := store.myMap[key]

	if !found {
		err = fmt.Errorf("value not found for key: %s", key)
		return "", err
	}

	return value, nil
}
