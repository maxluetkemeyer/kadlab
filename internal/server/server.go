package server

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/store"
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
	store        store.Store
}

// NewServer returns a new instance of Server
func NewServer(routingTable *routingtable.RoutingTable, store store.Store) *Server {
	return &Server{
		id:           routingTable.Me().ID,
		address:      routingTable.Me().Address,
		routingTable: routingTable,
		store:        store,
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
	log.Printf("received ping from\nNode: %v\nAddress: %v\n", hex.EncodeToString(sender.ID.Value), sender.IPWithPort)

	if len(sender.ID.Value) != env.IDLength {
		return nil, fmt.Errorf("invalid id length %v", len(sender.ID.Value))
	}

	c := contact.NewContact((*kademliaid.KademliaID)(sender.ID.Value), sender.IPWithPort)

	s.routingTable.AddContact(c)

	return &pb.Node{
		ID:         &pb.KademliaID{Value: s.id.Bytes()},
		IPWithPort: s.address,
	}, nil
}

func (s *Server) FindNode(ctx context.Context, request *pb.FindNodeRequest) (*pb.Nodes, error) {
	// TODO this needs to make sure that sender node and itself is not included in response
	candidates := s.routingTable.FindClosestContacts((*kademliaid.KademliaID)(request.Target.Value), (*kademliaid.KademliaID)(request.Sender.Value), env.BucketSize)

	nodes := make([]*pb.Node, 0, len(candidates))

	for _, contact := range candidates {
		nodes = append(nodes, &pb.Node{ID: &pb.KademliaID{Value: contact.ID.Bytes()}, IPWithPort: contact.Address})
	}

	return &pb.Nodes{Node: nodes}, nil
}
