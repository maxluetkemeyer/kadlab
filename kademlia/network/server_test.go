package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	"d7024e_group04/kademlia/routingtable"
	pb "d7024e_group04/proto"
)

const bufSize = 1024 * 1024

var (
	id      *kademliaid.KademliaID
	lis     *bufconn.Listener
	address = ":50051"
)

func init() {
	os.Setenv("BUCKET_SIZE", "20")
	id = kademliaid.NewRandomKademliaID()
}

func initBufconn() {
	c := contact.NewContact(id, address)
	routingTable := routingtable.NewRoutingTable(c)

	server := NewServer(address, id, routingTable)
	lis = bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	pb.RegisterKademliaServer(grpcServer, server)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestServer_Serve(t *testing.T) {
	routingTable := routingtable.NewRoutingTable(contact.NewContact(id, address))
	server := NewServer(address, id, routingTable)

	t.Run("start and stop", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)

		go func() {
			// timeout, did not shutdown on context cancel
			time.Sleep(250 * time.Millisecond)
			cancel()
			panic("did not shutdown grpc server on context cancel")
		}()

		err := server.Start(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestServer_Ping(t *testing.T) {
	initBufconn()

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewKademliaClient(conn)

	t.Run("ping valid node", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		sender := &pb.Node{
			ID:      kademliaid.NewRandomKademliaID().Bytes(),
			Address: "sender ip",
		}

		resp, err := client.Ping(ctx, sender)

		if err != nil {
			t.Error(fmt.Errorf("rpc ping failed: %v", err))
		}

		if !reflect.DeepEqual(resp.ID, id.Bytes()) {
			t.Error(fmt.Errorf("wrong id from responding node, got %v wanted %v", resp.ID, id.Bytes()))
		}

		if resp.Address != address {
			t.Error(fmt.Errorf("wrong address from responding node, got %v wanted %v", resp.Address, address))
		}
	})

	t.Run("ping with invalid node id", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		sender := &pb.Node{
			ID:      kademliaid.NewRandomKademliaID().Bytes()[5:],
			Address: "sender ip",
		}

		if _, err := client.Ping(ctx, sender); err == nil {
			t.Errorf("ping with invalid node id did not fail")
		}
	})
}
