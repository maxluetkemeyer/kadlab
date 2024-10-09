package cli

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/mock"
	"testing"
	"time"
)

var (
	me = contact.NewContact(kademliaid.NewRandomKademliaID(), "client_address")
)

func TestCLI_getCommand(t *testing.T) {
	node := mock.NewNodeMock(me)

	data := "some_data"
	hash := kademliaid.NewKademliaIDFromData(data).String()

	t.Run("value does not exists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		str, err := getCommand(ctx, node, hash)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if len(str) <= 0 {
			t.Fatalf("got empty string as result")
		}
	})

	t.Run("value does exist", func(t *testing.T) {
		node.Store[hash] = data
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		str, err := getCommand(ctx, node, hash)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if len(str) <= 0 {
			t.Fatalf("got empty string as result")
		}
	})
}
