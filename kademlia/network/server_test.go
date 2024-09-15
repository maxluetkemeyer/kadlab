package server

import (
	"context"
	"d7024e_group04/kademlia/network"
	"fmt"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/routingtable"
	pb "d7024e_group04/proto"
)

func TestServer_Serve(t *testing.T) {
	routingTable := routingtable.NewRoutingTable(contact.NewContact(network.targetID, network.targetAddress))
	server := NewServer(network.targetAddress, network.targetID, routingTable)

	t.Run("start and stop", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		go network.timeoutContext(ctx, cancel)

		err := server.Start(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestServer_Ping(t *testing.T) {
	network.initBufconn()

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(network.bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewKademliaClient(conn)

	t.Run("ping valid node", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		go network.timeoutContext(ctx, cancel)

		sender := &pb.Node{
			ID:         &pb.KademliaID{Value: network.clientID.Bytes()},
			IPWithPort: network.clientAddress,
		}

		resp, err := client.Ping(ctx, sender)

		if err != nil {
			t.Error(fmt.Errorf("rpc ping failed: %v", err))
		}

		if !reflect.DeepEqual(resp.ID.Value, network.targetID.Bytes()) {
			t.Error(fmt.Errorf("wrong id from responding node, got %v wanted %v", resp.ID.Value, network.targetID.Bytes()))
		}

		if resp.IPWithPort != network.targetAddress {
			t.Error(fmt.Errorf("wrong address from responding node, got %v wanted %v", resp.IPWithPort, network.targetAddress))
		}
	})

	t.Run("ping with invalid node id", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		go network.timeoutContext(ctx, cancel)

		sender := &pb.Node{
			ID:         &pb.KademliaID{Value: network.clientID.Bytes()[:5]},
			IPWithPort: network.clientAddress,
		}

		if _, err := client.Ping(ctx, sender); err == nil {
			t.Errorf("ping with invalid node id did not fail")
		}
	})
}
