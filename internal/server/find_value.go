package server

import (
	"context"
	pb "d7024e_group04/proto"
	"fmt"
)

// TODO: Input Validation, Tests, Error handling in string, getvalue, ...
func (s *Server) FindValue(ctx context.Context, request *pb.FindValueRequest) (*pb.NodesOrData, error) {
	key := string(request.Hash)

	value, err := s.store.GetValue(key)

	if err == nil {
		return &pb.NodesOrData{
			Value: &pb.NodesOrData_Data{
				Data: value,
			},
		}, nil
	}

	payload := &pb.FindNodeRequest{
		TargetID: request.Hash,
		SenderID: s.id.Bytes(),
	}

	nodes, errFindNodes := s.FindNode(ctx, payload)

	if errFindNodes != nil {
		err = fmt.Errorf("find Node error in Find value with error: %s", errFindNodes)
		return nil, err
	}

	return &pb.NodesOrData{
		Value: &pb.NodesOrData_Nodes{
			Nodes: nodes,
		},
	}, nil
}
