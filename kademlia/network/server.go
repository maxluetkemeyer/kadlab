package network

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

// Server represents the node grpc server
type Server struct {
	pb.UnimplementedKademliaServer
	id           *kademliaid.KademliaID
	address      string
	routingTable *routingtable.RoutingTable
}

// NewServer returns a new instance of Server
func NewServer(address string, id *kademliaid.KademliaID, routingTable *routingtable.RoutingTable) *Server {
	return &Server{
		id:           id,
		address:      address,
		routingTable: routingTable,
	}
}

// Start listens on the grpc port and serves rpc calls
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterKademliaServer(grpcServer, s)

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

// Ping serves PING rpc calls, saves the senders contact info and replies with it's own contact info
func (s *Server) Ping(ctx context.Context, sender *pb.Node) (*pb.Node, error) {
	log.Printf("received ping from\nNode: %v\nAddress: %v\n", hex.EncodeToString(sender.ID), sender.Address)

	if len(sender.ID) != 20 {
		return nil, fmt.Errorf("invalid id length %v", len(sender.ID))
	}

	c := contact.NewContact((*kademliaid.KademliaID)(sender.ID), sender.Address)

	s.routingTable.AddContact(c)

	return &pb.Node{
		ID:      s.id[:],
		Address: s.address,
	}, nil
}

func (s *Server) LookupContact(target *contact.Contact) {
	// TODO
}

func (s *Server) LookupData(hash string) {
	// TODO
}

func (s *Server) Store(data []byte) {
	// TODO
}
