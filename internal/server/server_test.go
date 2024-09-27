package server

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/store"
	"fmt"
	"reflect"
	"testing"
	"time"

	pb "d7024e_group04/proto"
)

var (
	ServerID      = kademliaid.NewRandomKademliaID()
	ServerAddress = ":50051"

	SenderID      = kademliaid.NewRandomKademliaID()
	SenderAddress = "sender_ip"
)

func initServer() *Server {
	routingTable := routingtable.NewRoutingTable(contact.NewContact(ServerID, ServerAddress))
	return NewServer(routingTable, store.NewMemoryStore())
}

func TestServer_Serve(t *testing.T) {
	server := initServer()
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
	server := initServer()

	t.Run("ping valid node", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		go TimeoutContext(ctx, cancel)

		sender := &pb.Node{
			ID:         SenderID.Bytes(),
			IPWithPort: SenderAddress,
		}

		resp, err := server.Ping(ctx, sender)

		if err != nil {
			t.Error(fmt.Errorf("rpc ping failed: %v", err))
		}

		if !reflect.DeepEqual(resp.ID, ServerID.Bytes()) {
			t.Error(fmt.Errorf("wrong id from responding node, got %v wanted %v", resp.ID, ServerID.Bytes()))
		}

		if resp.IPWithPort != ServerAddress {
			t.Error(fmt.Errorf("wrong address from responding node, got %v wanted %v", resp.IPWithPort, ServerAddress))
		}
	})

	t.Run("ping with invalid node id", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		go TimeoutContext(ctx, cancel)

		sender := &pb.Node{
			ID:         SenderID.Bytes()[:5],
			IPWithPort: SenderAddress,
		}

		if _, err := server.Ping(ctx, sender); err == nil {
			t.Errorf("ping with invalid node id did not fail")
		}
	})
}

func TestServer_FindNode(t *testing.T) {
	t.Run("find node", func(t *testing.T) {
		targetID := kademliaid.NewRandomKademliaID()

		srv := initServer()
		fillRoutingTable(env.BucketSize*2, srv.routingTable, SenderID)

		if candidates := srv.routingTable.FindClosestContacts(SenderID, 1); candidates[0].ID == SenderID {
			t.Fatalf("sender already exists in routing table")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		go TimeoutContext(ctx, cancel)

		resp, err := srv.FindNode(ctx, &pb.FindNodeRequest{TargetID: targetID.Bytes(), RequestingNode: &pb.Node{ID: SenderID.Bytes(), IPWithPort: SenderAddress}})

		if err != nil {
			t.Fatalf("rpc FindNode failed, %v", err)
		}

		nodes := resp.Nodes

		if len(nodes) <= 0 {
			t.Fatalf("empty response")
		}

		if len(nodes) > env.BucketSize {
			t.Fatalf("returned more nodes than allowed k: %v, got %v", env.BucketSize, len(nodes))
		}

		for _, node := range nodes {
			if reflect.DeepEqual(node.ID, ServerID.Bytes()) {
				t.Fatalf("response included server node")
			}

			if reflect.DeepEqual(node.ID, SenderID.Bytes()) {
				t.Fatalf("reponse included sender node")
			}
		}
	})
}

func fillRoutingTable(count int, routingTable *routingtable.RoutingTable, blacklist kademliaid.KademliaID) {
	var kID kademliaid.KademliaID

	for range count {
		for {
			kID = kademliaid.NewRandomKademliaID()
			if kID != blacklist {
				break
			}
		}

		contact0 := contact.NewContact(kID, fmt.Sprintf("node %v", count))

		routingTable.AddContact(contact0)
	}
}

func TimeoutContext(ctx context.Context, cancel context.CancelFunc) {
	<-ctx.Done()
	// timeout test, did not shut down on context cancel
	time.Sleep(30 * time.Second)
	cancel()
	panic("context timed out but test did not finish")
}
