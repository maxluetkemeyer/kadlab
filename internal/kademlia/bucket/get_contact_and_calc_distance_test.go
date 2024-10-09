package bucket_test

import (
	"d7024e_group04/internal/kademlia/bucket"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/mock"
	"testing"
)

func TestBucket_GetContactAndCalcDistance(t *testing.T) {
	client := mock.NewClientMock(nil)
	client.SetPingResult(true)

	contact0 := contact.NewContact(kademliaid.NewKademliaID("0000000000000000000000000000000000000000"), "0")
	contact1 := contact.NewContact(kademliaid.NewKademliaID("0000000000000000000000000000000000000001"), "1")

	distance := []kademliaid.KademliaID{
		kademliaid.NewKademliaID("0000000000000000000000000000000000000011"),
		kademliaid.NewKademliaID("0000000000000000000000000000000000000010"),
	}

	targetContact := contact.NewContact(kademliaid.NewKademliaID("0000000000000000000000000000000000000011"), "2")

	bucket := bucket.NewBucket(2)
	bucket.AddContact(contact0, client)
	bucket.AddContact(contact1, client)

	contactsWithDistance := bucket.GetContactAndCalcDistance(targetContact.ID)

	for i, contact := range contactsWithDistance {
		if !contact.Distance.Equals(distance[i]) {
			t.Fatalf("distance is not correct, wanted %v, got %v", distance[i], contact.Distance)
		}
	}

}
