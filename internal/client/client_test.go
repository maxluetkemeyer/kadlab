package client

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"reflect"
	"testing"
	"time"
)

var (
	clientID      = kademliaid.NewRandomKademliaID()
	clientAddress = "client address"

	serverID      = kademliaid.NewRandomKademliaID()
	serverAddress = "server address"
)

func TestClient_Ping(t *testing.T) {
	client := NewClient()

	clientContact := contact.NewContact(clientID, clientAddress)
	serverContact := contact.NewContact(serverID, serverAddress)

	t.Run("ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		go TimeoutContext(ctx, cancel)

		resp, err := client.SendPing(ctx, &clientContact, serverContact.Address)

		if err != nil {
			t.Fatalf("failed to ping, %v", err)
		}

		if !reflect.DeepEqual(resp.ID.Bytes(), serverID.Bytes()) {
			t.Fatalf("wrong node id in response, wanted %v got %v", serverID, resp.ID)
		}

		if resp.Address != serverAddress {
			t.Fatalf("wrong node address in response, wanted %v got %v", serverAddress, resp.Address)
		}
	})
}

func TimeoutContext(ctx context.Context, cancel context.CancelFunc) {
	<-ctx.Done()
	// timeout test, did not shutdown on context cancel
	time.Sleep(30 * time.Second)
	cancel()
	panic("context timed out but test did not finish")
}
