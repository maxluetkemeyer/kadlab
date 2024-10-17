package store

import (
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"testing"
)

func TestMemoryStore(t *testing.T) {
	t.Run("Store a value and retrieve it", func(t *testing.T) {
		memStore := NewMemoryStore()

		key := "myKey"
		value := "myValue"
		memStore.SetValue(key, value, nil)

		want := value
		got, err := memStore.GetValue(key)

		if err != nil {
			t.Fatal(err)
		}

		if got.Data != want {
			t.Errorf("got %s, want %s", got, want)
		}

	})

	t.Run("Store a value and retrieve a nonexistent key", func(t *testing.T) {
		memStore := NewMemoryStore()

		key := "myKey"
		value := "myValue"
		memStore.SetValue(key, value, nil)

		_, err := memStore.GetValue("anotherKey")

		if err == nil {
			t.Errorf("got nil, want error, %s", err)
		}
	})
}

func TestMemoryStore_GetOriginalUploader(t *testing.T) {
	memStore := NewMemoryStore()
	uploadingContact := contact.NewContact(kademliaid.NewRandomKademliaID(), "address")

	key := "myKey"
	value := "myValue"
	memStore.SetValue(key, value, uploadingContact)

	t.Run("existing data", func(t *testing.T) {
		contact, err := memStore.GetOriginalUploader(key)

		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if !contact.ID.Equals(uploadingContact.ID) {
			t.Fatalf("mismatching id, expected %v, got %v", uploadingContact.ID, contact.ID)
		}

		if contact.Address != uploadingContact.Address {
			t.Fatalf("mismatching address, expected %v, got %v", uploadingContact.Address, contact.Address)
		}
	})

	t.Run("non-existent data", func(t *testing.T) {
		contact, err := memStore.GetOriginalUploader("non-existent")
		if err == nil {
			t.Fatalf("expected error, found contact: %v", contact)
		}
	})

}
