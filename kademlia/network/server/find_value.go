package server

import (
	"context"
	pb "d7024e_group04/proto"
	"fmt"
)

// TODO: Input Validation, Tests, Error handling in string, getvalue, ...
func (s *Server) FindValue(ctx context.Context, kademliaID *pb.KademliaID) (*pb.NodesOrData, error) {
	key := string(kademliaID.Value)

	value, err := s.store.GetValue(key)

	if err == nil {
		return &pb.NodesOrData{
			Value: &pb.NodesOrData_Data{
				Data: value,
			},
		}, nil
	}

	nodes, errFindNodes := s.FindNode(ctx, kademliaID)

	if errFindNodes != nil {
		err = fmt.Errorf("Find Node error in Find value with error: %s", errFindNodes)
		return nil, err
	}

	return &pb.NodesOrData{
		Value: &pb.NodesOrData_Nodes{
			Nodes: nodes,
		},
	}, nil
}
