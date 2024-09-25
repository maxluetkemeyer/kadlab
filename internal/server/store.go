package server

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	pb "d7024e_group04/proto"
	"log"
)

// TODO: Input validation, Tests
func (s *Server) Store(ctx context.Context, content *pb.StoreRequest) (*pb.StoreResult, error) {
	log.Print("received store rpc")

	// TODO kademliaidfrombytes
	senderID := (kademliaid.KademliaID)(content.RequestingNode.ID)
	senderContact := *contact.NewContact(senderID, content.RequestingNode.IPWithPort)
	s.routingTable.AddContact(senderContact)

	key := content.Key
	value := content.Value

	s.store.SetValue(string(key), value)

	return &pb.StoreResult{
		Success: true,
	}, nil
}
