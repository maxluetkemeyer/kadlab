package bucket_test

import (
	"testing"

	"d7024e_group04/internal/kademlia/bucket"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/mock"
)

var getList = bucket.ExportGetList

func TestBucket(t *testing.T) {
	client := mock.NewClientMock(nil)
	client.SetPingResult(true)

	t.Run("New Bucket should be empty", func(t *testing.T) {
		bucket := bucket.NewBucket(20)

		want := 0
		got := getList(bucket).Len()

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Bucket length", func(t *testing.T) {
		bucket := bucket.NewBucket(20)

		contact0 := contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), "")
		contact1 := contact.NewContact(kademliaid.NewKademliaID("1111111200000000000000000000000000000000"), "")
		contact2 := contact.NewContact(kademliaid.NewKademliaID("1111111300000000000000000000000000000000"), "")

		bucket.AddContact(contact0, client)
		bucket.AddContact(contact1, client)
		bucket.AddContact(contact2, client)

		want := 3
		got := bucket.Len()

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
