package client

import (
	"context"
	"d7024e_group04/internal/kademlia/kademliaid"
	pb "d7024e_group04/proto"

	"google.golang.org/grpc"
)

type mockGrpcClient struct {
	id      kademliaid.KademliaID
	address string
	store   map[string]string
}

func (m *mockGrpcClient) Ping(ctx context.Context, in *pb.Node, opts ...grpc.CallOption) (*pb.Node, error) {
	return &pb.Node{ID: &pb.KademliaID{Value: m.id.Bytes()}, IPWithPort: m.address}, nil
}

func (m *mockGrpcClient) Find_Node(ctx context.Context, in *pb.KademliaID, opts ...grpc.CallOption) (*pb.Nodes, error) {
	panic("TODO")
}

func (m *mockGrpcClient) Find_Value(ctx context.Context, in *pb.KademliaID, opts ...grpc.CallOption) (*pb.NodesOrData, error) {
	panic("TODO")
}

func (m *mockGrpcClient) Store(ctx context.Context, in *pb.Content, opts ...grpc.CallOption) (*pb.StoreResult, error) {
	m.store[in.Key.String()] = in.Value
	return &pb.StoreResult{Success: true}, nil
}

func newMockGrpcClient(id kademliaid.KademliaID, address string) *mockGrpcClient {
	return &mockGrpcClient{
		id:      id,
		address: address,
		store:   make(map[string]string),
	}
}
