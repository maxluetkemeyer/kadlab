package store

import (
	"errors"
	"fmt"
	"log"
)

type MemoryStore struct {
	myMap map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		myMap: make(map[string][]byte),
	}
}

func (store *MemoryStore) SetValue(key string, value []byte) {
	log.Printf("stored value for key: %s", key)
	store.myMap[key] = value
}

func (store *MemoryStore) GetValue(key string) (value []byte, error error) {
	val, found := store.myMap[key]

	if !found {
		err := fmt.Errorf("value not found for key: %s", key)
		log.Printf(err.Error())
		return nil, errors.New(err.Error())
	}

	return val, nil
}
