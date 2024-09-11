package kademlia

import (
	"context"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	"d7024e_group04/kademlia/routingtable"
	pb "d7024e_group04/proto"
	"encoding/hex"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Node struct {
	pb.UnimplementedKademliaServer
	ID           *kademliaid.KademliaID
	Address      string
	routingTable *routingtable.RoutingTable
}

func NewNode(address string) *Node {
	id := kademliaid.NewRandomKademliaID()
	c := contact.NewContact(id, address)
	return &Node{
		ID:           id,
		Address:      address,
		routingTable: routingtable.NewRoutingTable(c),
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
	log.Printf("received ping from \nnode: %v\naddress: %v\n", hex.EncodeToString(sender.ID), sender.Address)

	if len(sender.ID) != 20 {
		return nil, fmt.Errorf("invalid id length %v", len(sender.ID))
	}

	c := contact.NewContact((*kademliaid.KademliaID)(sender.ID), sender.Address)

	n.routingTable.AddContact(c)

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
