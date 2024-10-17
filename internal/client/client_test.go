package client

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/mock"
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
	mock.StartMockGrpcServer(serverID, serverAddress)

	clientContact := contact.NewContact(clientID, clientAddress)
	client := NewClient(clientContact, grpc.WithContextDialer(mock.BufDialer))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	go TimeoutContext(ctx, cancel)

	resp, err := client.SendPing(ctx, mock.MockServerAddress)

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
	server := mock.StartMockGrpcServer(serverID, serverAddress)
	serverContact := contact.NewContact(serverID, mock.MockServerAddress)

	clientContact := contact.NewContact(clientID, clientAddress)
	client := NewClient(clientContact, grpc.WithContextDialer(mock.BufDialer))

	contacts := server.FillRoutingTable(env.BucketSize)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	go TimeoutContext(ctx, cancel)

	candidates, err := client.SendFindNode(ctx, serverContact, clientContact)
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
	server := mock.StartMockGrpcServer(serverID, serverAddress)

	clientContact := contact.NewContact(clientID, clientAddress)
	client := NewClient(clientContact, grpc.WithContextDialer(mock.BufDialer))

	t.Run("Data exists on node", func(t *testing.T) {
		value := "some_value"
		hash := kademliaid.NewKademliaIDFromData(value)
		server.TTLStore.SetValue(hash.String(), value, time.Hour, clientContact)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		go TimeoutContext(ctx, cancel)

		targetNode := contact.NewContact(serverID, mock.MockServerAddress)
		candidates, data, err := client.SendFindValue(ctx, targetNode, hash.String())
		if err != nil {
			t.Fatalf("failed to send find value request, %v", err)
		}

		if candidates != nil {
			t.Fatalf("expected no candidates to be returned, got %v", len(candidates))
		}

		if data.DataValue != value {
			t.Fatalf("invalid data returned, got %v, expected %v", data, value)
		}
	})

	t.Run("Data does not exist on node", func(t *testing.T) {
		value := "non-existent"
		hash := kademliaid.NewKademliaIDFromData(value)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		go TimeoutContext(ctx, cancel)

		targetNode := contact.NewContact(serverID, mock.MockServerAddress)
		candidates, data, err := client.SendFindValue(ctx, targetNode, hash.String())
		if err != nil {
			t.Fatalf("failed to send find value request, %v", err)
		}

		if len(data.DataValue) != 0 {
			t.Fatalf("expected no data to be found, got %v", data)
		}

		if candidates == nil {
			t.Fatalf("no candidates returned in response")
		}
	})
}

func TestClient_Store(t *testing.T) {
	server := mock.StartMockGrpcServer(serverID, serverAddress)
	targetNode := contact.NewContact(serverID, mock.MockServerAddress)

	clientContact := contact.NewContact(clientID, clientAddress)
	client := NewClient(clientContact, grpc.WithContextDialer(mock.BufDialer))

	value := "some_value"
	hash := kademliaid.NewKademliaIDFromData(value)

	t.Run("store data", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		go TimeoutContext(ctx, cancel)

		err := client.SendStore(ctx, targetNode, value, clientContact)

		if err != nil {
			t.Fatalf("failed to send store request")
		}

		serverValue, err := server.TTLStore.GetValue(hash.String())

		if err != nil {
			t.Fatalf("err: %v", err)
		}

		if serverValue.Data != value {
			t.Fatalf("data stored in server is not same as sent data, expected %v, got %v", value, serverValue)
		}

	})
}

func TestClient_SendRefreshTTL(t *testing.T) {
	server := mock.StartMockGrpcServer(serverID, serverAddress)
	targetNode := contact.NewContact(serverID, mock.MockServerAddress)

	clientContact := contact.NewContact(clientID, clientAddress)
	client := NewClient(clientContact, grpc.WithContextDialer(mock.BufDialer))

	key := "some_key"
	t.Run("send refreshTTL", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		go TimeoutContext(ctx, cancel)

		err := client.SendRefreshTTL(ctx, key, targetNode)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		ttl := server.TTLStore.GetTTL(key)
		if ttl.Seconds() <= 0 {
			t.Fatalf("invalid ttl: %v", ttl)
		}
	})
}

func TestClient_SendNewStoredLocation(t *testing.T) {
	server := mock.StartMockGrpcServer(serverID, serverAddress)
	targetNode := contact.NewContact(serverID, mock.MockServerAddress)

	clientContact := contact.NewContact(clientID, clientAddress)
	client := NewClient(clientContact, grpc.WithContextDialer(mock.BufDialer))

	key := "some_key"

	t.Run("add new stored location", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		go TimeoutContext(ctx, cancel)

		err := client.SendNewStoredLocation(ctx, key, targetNode, clientContact)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		locationContacts := server.TTLStore.GetStoreLocations(key)

		if len(locationContacts) != 1 {
			t.Fatalf("invalid number of locations: %v", locationContacts)
		}

		if !locationContacts[0].ID.Equals(clientContact.ID) {
			t.Fatalf("invalid id, expected %v, got %v", clientContact.ID, locationContacts[0].ID)
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
