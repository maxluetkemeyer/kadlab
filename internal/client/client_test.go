package client

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
)

var (
	clientID      = kademliaid.NewRandomKademliaID()
	clientAddress = "client address"

	serverID      = kademliaid.NewRandomKademliaID()
	serverAddress = "server address"
)

func TestClient_SendPing(t *testing.T) {
	startMockGrpcServer(serverID, serverAddress)
	client := NewClient(grpc.WithContextDialer(bufDialer))

	clientContact := contact.NewContact(clientID, clientAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	go TimeoutContext(ctx, cancel)

	resp, err := client.SendPing(ctx, *clientContact, mockServerAddress)

	if err != nil {
		t.Fatalf("failed to ping, %v", err)
	}

	if !reflect.DeepEqual(resp.ID.Bytes(), serverID.Bytes()) {
		t.Fatalf("wrong node id in response, wanted %v got %v", serverID, resp.ID)
	}

	if resp.Address != serverAddress {
		t.Fatalf("wrong node address in response, wanted %v got %v", serverAddress, resp.Address)
	}
}

func TestClient_SendFindNode(t *testing.T) {
	server := startMockGrpcServer(serverID, serverAddress)

	client := NewClient(grpc.WithContextDialer(bufDialer))
	clientContact := contact.NewContact(clientID, clientAddress)

	contacts := server.fillRoutingTable(env.BucketSize)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	go TimeoutContext(ctx, cancel)

	candidates, err := client.SendFindNode(ctx, *clientContact, mockServerAddress, clientContact.ID)
	if err != nil {
		t.Fatalf("failed to send find node request, %v", err)
	}

	if len(candidates) != len(contacts) {
		t.Fatalf("wrong number of contacts in response, got %v, expected %v", len(candidates), env.BucketSize)
	}

	if !reflect.DeepEqual(candidates, contacts) {
		t.Fatalf("wrong contacts returned")
	}
}

func TestClient_SendFindValue(t *testing.T) {
	server := startMockGrpcServer(serverID, serverAddress)

	client := NewClient(grpc.WithContextDialer(bufDialer))
	clientContact := contact.NewContact(clientID, clientAddress)

	t.Run("Data exists on node", func(t *testing.T) {
		value := "some_value"
		hash := kademliaid.NewKademliaIDFromData(value)
		server.DataStore[string(hash.Bytes())] = value

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		go TimeoutContext(ctx, cancel)

		targetNode := contact.NewContact(serverID, mockServerAddress)
		candidates, data, err := client.SendFindValue(ctx, *clientContact, *targetNode, value)
		if err != nil {
			t.Fatalf("failed to send find value request, %v", err)
		}

		if candidates != nil {
			t.Fatalf("expected no candidates to be returned, got %v", len(candidates))
		}

		if data != value {
			t.Fatalf("invalid data returned, got %v, expected %v", data, value)
		}
	})

	t.Run("Data does not exist on node", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		go TimeoutContext(ctx, cancel)

		targetNode := contact.NewContact(serverID, mockServerAddress)
		candidates, data, err := client.SendFindValue(ctx, *clientContact, *targetNode, "non-existent")
		if err != nil {
			t.Fatalf("failed to send find value request, %v", err)
		}

		if len(data) != 0 {
			t.Fatalf("expected no data to be found, got %v", data)
		}

		if candidates == nil {
			t.Fatalf("no candidates returned in response")
		}
	})

}

func TimeoutContext(ctx context.Context, cancel context.CancelFunc) {
	<-ctx.Done()
	// timeout test, did not shut down on context cancel
	time.Sleep(30 * time.Second)
	cancel()
	panic("context timed out but test did not finish")
}
