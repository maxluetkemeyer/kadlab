package network

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	pb "d7024e_group04/proto"
)

type ClientRPC interface {
	SendPing(ctx context.Context, grpc pb.KademliaClient, me, target *contact.Contact) (contact.Contact, error)
	SendFindNode(ctx context.Context, contact *contact.Contact) ([]contact.Contact, error)
	SendFindValue(ctx context.Context, hash *kademliaid.KademliaID) (string, error)
	SendStore(ctx context.Context, data string) error
}

type ServerRPC interface {
	Start(ctx context.Context) error
	Ping(ctx context.Context, sender *pb.Node) (*pb.Node, error)
	FindValue(ctx context.Context, kademliaID *pb.KademliaID) (*pb.NodesOrData, error)
	FindNode(ctx context.Context, kademliaID *pb.KademliaID) (*pb.Nodes, error)
	Store(ctx context.Context, context *pb.Content) (*pb.StoreResult, error)
}

type Network interface {
	ResolveDNS(ctx context.Context, domain string) []string
}
