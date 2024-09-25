package contact

import (
	"d7024e_group04/internal/kademlia/kademliaid"
	"fmt"
	"testing"
)

func TestSet_Add(t *testing.T) {
	set := NewContactSet()

	contact := NewContact(kademliaid.NewRandomKademliaID(), "address")

	set.Add(*contact)

	if len(set.m) != 1 {
		t.Fatalf("wrong number of items in set, found %v, expected 1", len(set.m))
	}

	_, found := set.m[*contact]

	if !found {
		t.Fatalf("contact is not in set")
	}
}
func TestSet_Adds(t *testing.T) {
	set := NewContactSet()

	contact1 := NewContact(kademliaid.NewRandomKademliaID(), "address 1")
	contact2 := NewContact(kademliaid.NewRandomKademliaID(), "address 2")

	contacts := []Contact{*contact1, *contact2}

	set.Adds(contacts)

	if len(set.m) != 2 {
		t.Fatalf("wrong number of items in set, found %v, expected %v", len(set.m), len(contacts))
	}

	for _, contact := range contacts {
		_, found := set.m[contact]
		if !found {
			t.Fatalf("contact is not in set")
		}
	}

}

func TestSet_Remove(t *testing.T) {
	set := NewContactSet()

	contacts := fillSet(set, 2)

	if len(set.m) != len(contacts) {
		t.Fatalf("invalid number of items in set, got %v, expected %v", len(set.m), len(contacts))
	}

	_, found := set.m[contacts[0]]

	if !found {
		t.Fatalf("contact is not in set")
	}

	set.Remove(contacts[0])

	_, found = set.m[contacts[0]]

	if found {
		t.Fatalf("contact is still in set")
	}

}

func TestSet_Has(t *testing.T) {
	set := NewContactSet()

	contacts := fillSet(set, 2)

	_, found := set.m[contacts[0]]

	if !found {
		t.Fatalf("contact is not in set")
	}

	if !set.Has(contacts[0]) {
		t.Fatalf("contact exists in set but not found in Has()")
	}
}

func TestSet_Len(t *testing.T) {
	set := NewContactSet()

	fillSet(set, 2)

	if len(set.m) != set.Len() {
		t.Fatalf("mismatching length, got %v, expected %v", set.Len(), len(set.m))
	}
}

func TestSet_Clear(t *testing.T) {
	set := NewContactSet()

	contacts := fillSet(set, 2)

	if len(set.m) != len(contacts) {
		t.Fatalf("mismatching length, got %v, expected %v", set.Len(), len(set.m))
	}

	set.Clear()

	if len(set.m) != 0 {
		t.Fatalf("set was not cleared, len of set is: %v", len(set.m))
	}

}

func TestSet_IsEmpty(t *testing.T) {
	set := NewContactSet()

	if len(set.m) != 0 {
		t.Fatalf("set is not empty, len of set is %v", len(set.m))
	}

	if !set.IsEmpty() {
		t.Fatalf("set is empty but IsEmpty() returned false")
	}
}

func fillSet(set *ContactSet, count int) (contacts []Contact) {
	for i := range count {
		contact := NewContact(kademliaid.NewRandomKademliaID(), fmt.Sprintf("address %v", i))
		contacts = append(contacts, *contact)
		set.m[*contact] = struct{}{}
	}
	return
}
