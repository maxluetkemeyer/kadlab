package bucket

import (
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	"testing"
)

func TestBucket(t *testing.T) {
	t.Run("Bucket length", func(t *testing.T) {
		bucket := NewBucket()

		contact0 := contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), "")
		contact1 := contact.NewContact(kademliaid.NewKademliaID("1111111200000000000000000000000000000000"), "")
		contact2 := contact.NewContact(kademliaid.NewKademliaID("1111111300000000000000000000000000000000"), "")

		bucket.AddContact(contact0)
		bucket.AddContact(contact1)
		bucket.AddContact(contact2)

		want := 3
		got := bucket.Len()

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
