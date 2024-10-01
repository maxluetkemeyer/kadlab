package server

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	pb "d7024e_group04/proto"
	"fmt"
	"log"
)

// TODO: Input validation
func (s *Server) Store(ctx context.Context, content *pb.StoreRequest) (*pb.StoreResult, error) {
	log.Print("received store rpc")

	senderID := kademliaid.NewKademliaIDFromBytes([env.IDLength]byte(content.RequestingNode.ID))
	senderContact := contact.NewContact(senderID, content.RequestingNode.IPWithPort)
	s.routingTable.AddContact(senderContact)

	key := content.Key
	value := content.Value

	if len(key) != env.IDLength || len(value) == 0 {
		return &pb.StoreResult{
				Success: false,
			},
			fmt.Errorf("invalid content: %v", content)
	}

	s.store.SetValue(string(key), value)

	return &pb.StoreResult{
		Success: true,
	}, nil
}
