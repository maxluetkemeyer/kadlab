package server

import (
	"context"
	pb "d7024e_group04/proto"
	"log"
)

// TODO: Input validation, Tests
func (s *Server) Store(ctx context.Context, content *pb.Content) (*pb.StoreResult, error) {
	log.Print("received store rpc")

	key := content.Key
	value := content.Value

	//TODO: Is string in proto better or bytes?
	s.store.SetValue(key.String(), []byte(value))

	return &pb.StoreResult{
		Success: true,
	}, nil
}
