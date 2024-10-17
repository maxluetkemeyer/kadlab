package store

import (
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"log"
	"testing"
	"time"
)

func TestSimpleTTLStore(t *testing.T) {
	t.Run("set value with enough time", func(t *testing.T) {
		memStore := NewMemoryStore()
		tStore := NewSimpleTTLStore(memStore)

		want := "value0"

		tStore.SetValue("key0", want, time.Hour, nil)

		got, err := tStore.GetValue("key0")

		if err != nil {
			log.Fatalf("error getting value: %v", err)
		}

		if got.Data != want {
			log.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("set value with zero time", func(t *testing.T) {
		memStore := NewMemoryStore()
		tStore := NewSimpleTTLStore(memStore)

		val := "value0"

		tStore.SetValue("key0", val, time.Second*0, nil)

		got, err := tStore.GetValue("key0")

		if got.Data != "" {
			log.Fatalf("should get empty string, but got %v", got)
		}

		if err == nil {
			log.Fatalf("should error, but got nil")
		}

	})

	t.Run("set value with enough ttl and later it should be expired", func(t *testing.T) {
		memStore := NewMemoryStore()
		tStore := NewSimpleTTLStore(memStore)

		val := "value0"

		tStore.SetValue("key0", val, time.Hour, nil)

		got, err := tStore.GetValue("key0")

		if err != nil {
			log.Fatalf("error getting value: %v", err)
		}

		if got.Data != val {
			log.Fatalf("TTL should be enough: got %q, want %q", got, val)
		}

		tStore.SetTTL("key0", time.Second*0)

		gotAfterExpiration, errAfterExpiration := tStore.GetValue("key0")

		if gotAfterExpiration.Data != "" {
			log.Fatalf("should get empty string, but got %v", gotAfterExpiration)
		}

		if errAfterExpiration == nil {
			log.Fatalf("should error, but got nil")
		}
	})
}

func TestSimpleTTLStore_GetOriginalUploader(t *testing.T) {
	memStore := NewMemoryStore()
	tStore := NewSimpleTTLStore(memStore)
	originalUploader := contact.NewContact(kademliaid.NewRandomKademliaID(), "address")

	val := "value0"
	key := "key0"
	tStore.SetValue(key, val, time.Hour, originalUploader)

	contact, err := tStore.GetOriginalUploader(key)
	if err != nil {
		t.Fatalf("got err: %v", err)
	}

	if !contact.ID.Equals(originalUploader.ID) {
		t.Fatalf("mismatching id: expected %v, got %v", originalUploader.ID, contact.ID)
	}

	if contact.Address != originalUploader.Address {
		t.Fatalf("mismatching address, expected %v, got %v", originalUploader.Address, contact.Address)
	}
}
func TestSimpleTTLStore_StoreLocations(t *testing.T) {
	memStore := NewMemoryStore()
	tStore := NewSimpleTTLStore(memStore)
	originalUploader := contact.NewContact(kademliaid.NewRandomKademliaID(), "address")
	replicateContact := contact.NewContact(kademliaid.NewRandomKademliaID(), "address2")

	val := "value0"
	key := "key0"
	tStore.SetValue(key, val, time.Hour, originalUploader)

	tStore.AddStoreLocation(key, replicateContact)

	storeLocations := tStore.GetStoreLocations(key)

	if len(storeLocations) != 1 {
		t.Fatalf("invalid number of contacts, got %v, expected 1", len(storeLocations))
	}

	if !storeLocations[0].ID.Equals(replicateContact.ID) {
		t.Fatalf("invalid contact: %v, expected %v or %v", storeLocations, originalUploader, replicateContact)

	}

}

func TestSimpleTTLStore_RemoveRefreshContact(t *testing.T) {
	memStore := NewMemoryStore()
	tStore := NewSimpleTTLStore(memStore)
	replicateContact := contact.NewContact(kademliaid.NewRandomKademliaID(), "address2")

	key := "key0"
	tStore.AddStoreLocation(key, replicateContact)

	if len(tStore.GetStoreLocations(key)) != 1 {
		t.Fatalf("invalid size, expected 1 got %v", len(tStore.GetStoreLocations(key)))
	}

	tStore.RemoveRefreshContact(key)

	if len(tStore.GetStoreLocations(key)) != 0 {
		t.Fatalf("invalid size, expected 0 got %v", len(tStore.GetStoreLocations(key)))
	}
}
