package mock

import (
	"context"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	server2 "d7024e_group04/kademlia/network/server"
	"d7024e_group04/kademlia/routingtable"
	"d7024e_group04/kademlia/store"
	pb "d7024e_group04/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	Lis *bufconn.Listener

	TargetID      = kademliaid.NewRandomKademliaID()
	TargetAddress = ":50051"

	ClientID      = kademliaid.NewRandomKademliaID()
	ClientAddress = "sender_ip"
)

func InitBufconn() {
	c := contact.NewContact(TargetID, TargetAddress)
	routingTable := routingtable.NewRoutingTable(c)

	server := server2.NewServer(TargetAddress, TargetID, routingTable, store.NewMemoryStore())
	Lis = bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	pb.RegisterKademliaServer(grpcServer, server)
	go func() {
		if err := grpcServer.Serve(Lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func BufDialer(context.Context, string) (net.Conn, error) {
	return Lis.Dial()
}

func TimeoutContext(ctx context.Context, cancel context.CancelFunc) {
	<-ctx.Done()
	// timeout test, did not shutdown on context cancel
	time.Sleep(30 * time.Second)
	cancel()
	panic("context timedout but test did not finish")
}
