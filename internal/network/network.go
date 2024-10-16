package network

import (
	"context"

	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/model"
	pb "d7024e_group04/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ClientRPC interface {
	SendPing(ctx context.Context, targetIpWithPort string) (*contact.Contact, error)
	SendFindNode(ctx context.Context, contactWeRequest, contactWeAreSearchingFor *contact.Contact) ([]*contact.Contact, error)
	SendFindValue(ctx context.Context, contactWeRequest *contact.Contact, hash string) (candidates []*contact.Contact, dataObject model.FindValueSuccessfulResponse, err error)
	SendStore(ctx context.Context, contactWeRequest *contact.Contact, data string, originalUploader *contact.Contact) error
	SendRefreshTTL(ctx context.Context, key string, contactWeRequest *contact.Contact) error
	SendNewStoredLocation(ctx context.Context, key string, originalUploader, newContactStoringData *contact.Contact) error
}

type ServerRPC interface {
	Start(ctx context.Context) error
	Ping(ctx context.Context, sender *pb.Node) (*pb.Node, error)
	FindValue(ctx context.Context, request *pb.FindValueRequest) (*pb.FindValueResult, error)
	FindNode(ctx context.Context, request *pb.FindNodeRequest) (*pb.FindNodeResult, error)
	Store(ctx context.Context, content *pb.StoreRequest) (*pb.StoreResult, error)
	RefreshTTL(ctx context.Context, request *pb.RefreshTTLRequest) (*emptypb.Empty, error)
	NewStoreLocation(ctx context.Context, request *pb.NewStoreLocationRequest) (*emptypb.Empty, error)
}

type Network interface {
	ResolveDNS(domain string) ([]string, error)
}
