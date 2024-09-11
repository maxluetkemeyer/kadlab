package bucket

import (
	"d7024e_group04/kademlia"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	"testing"
)

func TestAddContact(t *testing.T) {

	cases := []struct {
		Name        string
		Bucket      *Bucket
		Contacts    []contact.Contact
		BucketSize  int
		ExpectedLen int
	}{
		{"One single contact",
			NewBucket(),
			[]contact.Contact{
				contact.NewContact(kademliaid.NewRandomKademliaID(), ""),
			},
			kademlia.BucketSize,
			1,
		},
		{"Two different contacts",
			NewBucket(),
			[]contact.Contact{
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), ""),
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000001"), ""),
			},
			kademlia.BucketSize,
			2,
		},
		{"Two similar contacts",
			NewBucket(),
			[]contact.Contact{
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), ""),
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), ""),
			},
			kademlia.BucketSize,
			1,
		},
		{"Bucket is full, two different contacts",
			NewBucket(),
			[]contact.Contact{
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000000"), ""),
				contact.NewContact(kademliaid.NewKademliaID("1111111100000000000000000000000000000002"), ""),
			},
			1,
			1,
		},
		// TODO: Add test cases for non-responding nodes
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {

			for _, myContact := range test.Contacts {
				test.Bucket.addContactCustom(myContact, test.BucketSize)
			}

			got := test.Bucket.Len()
			want := test.ExpectedLen

			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
	}

	// Dont delete, this tests AddContact, the other ones test AddContactCustom
	t.Run("Bucket length should have increased after insertion a new unknown contact", func(t *testing.T) {
		bucket := NewBucket()
		contact0 := contact.NewContact(kademliaid.NewRandomKademliaID(), "")

		bucket.AddContact(contact0)

		want := 1
		got := bucket.list.Len()

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
