package bucket_test

import (
	"d7024e_group04/internal/kademlia/bucket"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/mock"
	"testing"
)

func TestAddContact(t *testing.T) {

	cases := []struct {
		Name        string
		Bucket      *bucket.Bucket
		Contacts    []*contact.Contact
		ExpectedLen int
	}{
		{"One single contact",
			bucket.NewBucket(20),
			[]*contact.Contact{
				contact.NewContact(kademliaid.NewRandomKademliaID(), ""),
			},
			1,
		},
		{"Two different contacts",
			bucket.NewBucket(20),
			[]*contact.Contact{
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), ""),
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000001"), ""),
			},
			2,
		},
		{"Two similar contacts",
			bucket.NewBucket(20),
			[]*contact.Contact{
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), ""),
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), ""),
			},
			1,
		},
		{"Bucket is full, two different contacts",
			bucket.NewBucket(1),
			[]*contact.Contact{
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), ""),
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000002"), ""),
			},
			1,
		},
		// TODO: Add test cases for non-responding nodes
	}

	client := mock.NewClientMock(nil)
	client.SetPingResult(true)

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {

			for _, myContact := range test.Contacts {
				test.Bucket.AddContact(myContact, client)
			}

			got := test.Bucket.Len()
			want := test.ExpectedLen

			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
	}

	// Don't delete, this tests AddContact, the other ones test AddContactCustom
	t.Run("Bucket length should have increased after insertion a new unknown contact", func(t *testing.T) {
		bucket := bucket.NewBucket(20)
		contact0 := contact.NewContact(kademliaid.NewRandomKademliaID(), "")

		bucket.AddContact(contact0, client)

		want := 1
		got := getList(bucket).Len()

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("ping fails", func(t *testing.T) {
		bucket := bucket.NewBucket(2)

		contact0 := contact.NewContact(kademliaid.NewRandomKademliaID(), "0")
		contact1 := contact.NewContact(kademliaid.NewRandomKademliaID(), "1")
		contact2 := contact.NewContact(kademliaid.NewRandomKademliaID(), "2")

		bucket.AddContact(contact0, client)
		bucket.AddContact(contact1, client)

		front := getList(bucket).Front()

		client.SetPingResult(false)

		bucket.AddContact(contact2, client)

		frontAfterFailedPing := getList(bucket).Front()

		if front == frontAfterFailedPing {
			t.Fatal("did not remove head when ping failed")
		}
	})
}
