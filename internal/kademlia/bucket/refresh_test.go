package bucket

import (
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/mock"
	"testing"
	"time"
)

func TestRefreshBucket(t *testing.T) {
	client := mock.NewClientMock(nil)
	client.SetPingResult(true)

	t.Run("Bucket shouldn't need a refresh if TRefresh has not yet happened", func(t *testing.T) {
		env.TRefresh = 50*time.Millisecond
		bucket := NewBucket(20)

		got := bucket.NeedsRefresh()
		want := false

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Bucket should need a refresh if TRefresh has happened", func(t *testing.T) {
		env.TRefresh = 50*time.Millisecond
		bucket := NewBucket(20)

		time.Sleep(env.TRefresh)

		got := bucket.NeedsRefresh()
		want := true

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Adding a contact should refresh the bucket", func(t *testing.T) {
		env.TRefresh = 50*time.Millisecond
		bucket := NewBucket(20)

		time.Sleep(env.TRefresh)

		contact0 := contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), "")
		bucket.AddContact(contact0, client)

		got := bucket.NeedsRefresh()
		want := false

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
