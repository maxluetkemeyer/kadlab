package mock

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	pb "d7024e_group04/proto"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const BufSize = 1024 * 1024
const MockServerAddress = "passthrough://bufnet"

var Lis *bufconn.Listener

type mockGrpcServer struct {
	pb.UnimplementedKademliaServer
	ServerContact *contact.Contact
	RoutingTable  []*contact.Contact
	DataStore     map[string]string
}

func BufDialer(context.Context, string) (net.Conn, error) {
	return Lis.Dial()
}

func StartMockGrpcServer(id kademliaid.KademliaID, address string) *mockGrpcServer {
	server := &mockGrpcServer{
		ServerContact: contact.NewContact(id, address),
		DataStore:     make(map[string]string),
	}

	Lis = bufconn.Listen(BufSize)
	s := grpc.NewServer()
	pb.RegisterKademliaServer(s, server)

	go func() {
		if err := s.Serve(Lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return server
}

func (m *mockGrpcServer) Ping(ctx context.Context, in *pb.Node) (*pb.Node, error) {
	return &pb.Node{ID: m.ServerContact.ID.Bytes(), IPWithPort: m.ServerContact.Address}, nil
}

func (m *mockGrpcServer) FindNode(ctx context.Context, in *pb.FindNodeRequest) (*pb.FindNodeResult, error) {
	nodes := make([]*pb.Node, 0, len(m.RoutingTable))

	for _, contact := range m.RoutingTable {
		nodes = append(nodes, &pb.Node{ID: contact.ID.Bytes(), IPWithPort: contact.Address})
	}

	return &pb.FindNodeResult{Nodes: nodes}, nil
}

func (m *mockGrpcServer) FindValue(ctx context.Context, in *pb.FindValueRequest) (*pb.FindValueResult, error) {
	value, found := m.DataStore[string(in.Hash)]
	if found {
		return &pb.FindValueResult{Value: &pb.FindValueResult_Data{Data: value}}, nil
	}

	return &pb.FindValueResult{Value: &pb.FindValueResult_Nodes{Nodes: &pb.FindNodeResult{}}}, nil
}

func (m *mockGrpcServer) Store(ctx context.Context, in *pb.StoreRequest) (*pb.StoreResult, error) {
	panic("TODO")
}

func (m *mockGrpcServer) FillRoutingTable(count int) (contacts []*contact.Contact) {
	for i := range count {
		id := kademliaid.NewRandomKademliaID()
		address := fmt.Sprintf("server %v", i)

		contacts = append(contacts, contact.NewContact(id, address))
	}
	m.RoutingTable = contacts
	return
}
