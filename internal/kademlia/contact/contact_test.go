package contact

import (
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/kademliaid"
	"reflect"
	"strconv"
	"testing"

	"golang.org/x/exp/rand"
)

func TestContact_CalcDistance(t *testing.T) {
	id0 := kademliaid.NewKademliaID("1111111100000000000000000000000000000000")
	id1 := kademliaid.NewKademliaID("1111111100000000000000000000000000000011")
	distance := "0000000000000000000000000000000000000011"

	contact0 := NewContact(id0, "addr 0")

	contact0.CalcDistance(id1)

	if contact0.distance.String() != distance {
		t.Fatalf("wrong distance, got %v, expected %v", contact0.distance.String(), distance)
	}
}

func TestContact_Less(t *testing.T) {
	distance0 := kademliaid.NewKademliaID("0000000000000000000000000000000000000000")
	distance1 := kademliaid.NewKademliaID("0000000000000000000000000000000000000001")
	contact0 := &Contact{distance: distance0}
	contact1 := &Contact{distance: distance1}

	if !contact0.Less(contact1) {
		t.Fatalf("expected contact0 to be less than contact1")
	}

	if contact1.Less(contact0) {
		t.Fatalf("expected contact1 to not be less than contact0")
	}
}

func TestContact_String(t *testing.T) {
	contact := NewContact(kademliaid.NewRandomKademliaID(), "address")

	str := contact.String()

	if str == "" {
		t.Fatalf("string is empty")
	}
}

func TestContact_SortContacts(t *testing.T) {
	contactsSorted := createContactSlice(10)

	contacts := make([]Contact, len(contactsSorted))
	_ = copy(contacts, contactsSorted)

	if len(contacts) != len(contactsSorted) || !reflect.DeepEqual(contacts, contactsSorted) {
		t.Fatalf("invalid copy of slice")
	}

	for reflect.DeepEqual(contacts, contactsSorted) {
		shuffle(contacts)
	}

	SortContacts(contacts)

	if !reflect.DeepEqual(contacts, contactsSorted) {
		t.Fatalf("did not sort correctly\ngot %v\nexpected %v", contacts[0].distance.Bytes(), contactsSorted[0].distance.Bytes())
	}
}

func TestContact_RemoveID(t *testing.T) {
	contacts := createContactSlice(3)

	removedContact := contacts[1]

	RemoveID(contacts, removedContact.ID)

	for _, contact := range contacts {
		if contact.ID.Equals(removedContact.ID) {
			t.Fatalf("failed to remove contact from slice")
		}
	}
}

func createContactSlice(count int) (contacts []Contact) {
	for i := range count {
		contact := Contact{
			ID:       kademliaid.NewKademliaIDFromData(strconv.Itoa(i)),
			distance: [env.IDLength]byte{byte(i)},
		}
		contacts = append(contacts, contact)
	}
	return
}

func shuffle(contacts []Contact) {
	rand.Shuffle(len(contacts), func(i, j int) { contacts[i], contacts[j] = contacts[j], contacts[i] })
}
