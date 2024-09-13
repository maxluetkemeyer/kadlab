package network

import (
	"context"
	"d7024e_group04/kademlia/contact"
	pb "d7024e_group04/proto"
)

type ClientRPC interface {
	SendPing(ctx context.Context, contact *contact.Contact) error
	SendFindNode(ctx context.Context, contact *contact.Contact)
	SendFindValue(ctx context.Context, hash string)
	SendStore(ctx context.Context, data []byte)
}

type ServerRPC interface {
	Start(ctx context.Context) error
	Ping(ctx context.Context, sender *pb.Node) (*pb.Node, error)
	FindValue(ctx context.Context, kademliaID *pb.KademliaID) (*pb.NodesOrData, error)
	FindNode(ctx context.Context, kademliaID *pb.KademliaID) (*pb.Nodes, error)
	Store(ctx context.Context, context *pb.Content) (*pb.StoreResult, error)
}
