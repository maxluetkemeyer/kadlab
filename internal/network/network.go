package network

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	pb "d7024e_group04/proto"
)

type ClientRPC interface {
	SendPing(ctx context.Context, me *contact.Contact, target string) (*contact.Contact, error)
	SendFindNode(ctx context.Context, me, target *contact.Contact) ([]contact.Contact, error)
	SendFindValue(ctx context.Context, me, target contact.Contact, hash string) ([]contact.Contact, string, error)
	SendStore(ctx context.Context, data string) error
}

type ServerRPC interface {
	Start(ctx context.Context) error
	Ping(ctx context.Context, sender *pb.Node) (*pb.Node, error)
	FindValue(ctx context.Context, request *pb.FindValueRequest) (*pb.FindValueResult, error)
	FindNode(ctx context.Context, request *pb.FindNodeRequest) (*pb.FindNodeResult, error)
	Store(ctx context.Context, context *pb.StoreRequest) (*pb.StoreResult, error)
}

type Network interface {
	ResolveDNS(domain string) ([]string, error)
}
