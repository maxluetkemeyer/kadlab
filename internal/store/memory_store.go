package store

import (
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/model"
	"fmt"
	"log"
	"sync"
)

type MemoryStore struct {
	mut   sync.RWMutex
	myMap map[string]model.DataWithOriginalUploader
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		myMap: make(map[string]model.DataWithOriginalUploader),
	}
}

func (store *MemoryStore) SetValue(key string, value string, uploader *contact.Contact) {
	log.Printf("stored value for key: %s", key)
	store.mut.Lock()
	defer store.mut.Unlock()

	store.myMap[key] = model.DataWithOriginalUploader{Data: value, Contact: uploader}
}

func (store *MemoryStore) GetValue(key string) (string, error) {
	store.mut.RLock()
	defer store.mut.RUnlock()

	object, found := store.myMap[key]

	if !found {
		err := fmt.Errorf("value not found for key: %s", key)
		return "", err
	}

	return object.Data, nil
}

func (store *MemoryStore) GetOriginalUploader(key string) (*contact.Contact, error) {
	store.mut.RLock()
	defer store.mut.RUnlock()

	object, found := store.myMap[key]

	if !found {
		err := fmt.Errorf("value not found for key: %s", key)
		return nil, err
	}

	return object.Contact, nil
}
