package kademlia

import (
	"context"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	pb "d7024e_group04/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Node struct {
	pb.UnimplementedKademliaServer
	ID      kademliaid.KademliaID
	Address string
}

func NewNode(address string) *Node {
	return &Node{
		ID:      *kademliaid.NewRandomKademliaID(),
		Address: address,
	}
}

func (n *Node) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", n.Address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterKademliaServer(grpcServer, n)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
		lis.Close()
	}()

	log.Printf("serving at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

// TODO: I do not expect these functions in a "n.go" file

func (n *Node) Ping(ctx context.Context, sender *pb.Node) (*pb.Node, error) {
	// TODO update bucket with sender info
	log.Printf("received ping\n")
	return &pb.Node{
		ID:      n.ID[:],
		Address: n.Address,
	}, nil
}

func (n *Node) LookupContact(target *contact.Contact) {
	// TODO
}

func (n *Node) LookupData(hash string) {
	// TODO
}

func (n *Node) Store(data []byte) {
	// TODO
}
