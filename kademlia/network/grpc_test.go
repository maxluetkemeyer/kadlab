package network

import (
	"context"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	server2 "d7024e_group04/kademlia/network/server"
	"d7024e_group04/kademlia/network/store"
	"d7024e_group04/kademlia/routingtable"
	pb "d7024e_group04/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis *bufconn.Listener

	targetID      = kademliaid.NewRandomKademliaID()
	targetAddress = ":50051"

	clientID      = kademliaid.NewRandomKademliaID()
	clientAddress = "sender_ip"
)

func initBufconn() {
	c := contact.NewContact(targetID, targetAddress)
	routingTable := routingtable.NewRoutingTable(c)

	server := server2.NewServer(targetAddress, targetID, routingTable, store.NewMemoryStore())
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

func timeoutContext(ctx context.Context, cancel context.CancelFunc) {
	<-ctx.Done()
	// timeout test, did not shutdown on context cancel
	time.Sleep(30 * time.Second)
	cancel()
	panic("context timedout but test did not finish")
}
