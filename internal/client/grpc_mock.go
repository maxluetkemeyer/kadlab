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
	return &pb.Node{ID: m.id.Bytes(), IPWithPort: m.address}, nil
}

func (m *mockGrpcClient) FindNode(ctx context.Context, in *pb.FindNodeRequest, opts ...grpc.CallOption) (*pb.Nodes, error) {
	panic("TODO")
}

func (m *mockGrpcClient) FindValue(ctx context.Context, in *pb.FindValueRequest, opts ...grpc.CallOption) (*pb.NodesOrData, error) {
	panic("TODO")
}

func (m *mockGrpcClient) Store(ctx context.Context, in *pb.Content, opts ...grpc.CallOption) (*pb.StoreResult, error) {
	m.store[string(in.Key)] = in.Value
	return &pb.StoreResult{Success: true}, nil
}

func newMockGrpcClient(id kademliaid.KademliaID, address string) *mockGrpcClient {
	return &mockGrpcClient{
		id:      id,
		address: address,
		store:   make(map[string]string),
	}
}
