package server

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/store"
	"fmt"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "d7024e_group04/proto"
)

func TestServer_Serve(t *testing.T) {
	routingTable := routingtable.NewRoutingTable(contact.NewContact(TargetID, TargetAddress))
	server := NewServer(routingTable, store.NewMemoryStore())

	t.Run("start and stop", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		go TimeoutContext(ctx, cancel)

		err := server.Start(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestServer_Ping(t *testing.T) {
	InitBufconn()

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(BufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewKademliaClient(conn)

	t.Run("ping valid node", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		go TimeoutContext(ctx, cancel)

		sender := &pb.Node{
			ID:         &pb.KademliaID{Value: ClientID.Bytes()},
			IPWithPort: ClientAddress,
		}

		resp, err := client.Ping(ctx, sender)

		if err != nil {
			t.Error(fmt.Errorf("rpc ping failed: %v", err))
		}

		if !reflect.DeepEqual(resp.ID.Value, TargetID.Bytes()) {
			t.Error(fmt.Errorf("wrong id from responding node, got %v wanted %v", resp.ID.Value, TargetID.Bytes()))
		}

		if resp.IPWithPort != TargetAddress {
			t.Error(fmt.Errorf("wrong address from responding node, got %v wanted %v", resp.IPWithPort, TargetAddress))
		}
	})

	t.Run("ping with invalid node id", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		go TimeoutContext(ctx, cancel)

		sender := &pb.Node{
			ID:         &pb.KademliaID{Value: ClientID.Bytes()[:5]},
			IPWithPort: ClientAddress,
		}

		if _, err := client.Ping(ctx, sender); err == nil {
			t.Errorf("ping with invalid node id did not fail")
		}
	})
}
