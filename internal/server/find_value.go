package server

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/utils"
	pb "d7024e_group04/proto"
	"encoding/hex"
	"fmt"
)

// TODO: Input Validation, Tests, Error handling in string, getvalue, ...
func (s *Server) FindValue(ctx context.Context, request *pb.FindValueRequest) (*pb.FindValueResult, error) {
	key := hex.EncodeToString(request.Hash)

	dataObject, err := s.store.GetValue(key)

	if err == nil {
		// Reset TTL
		s.store.SetTTL(key, env.TTL)

		return &pb.FindValueResult{
			Value: &pb.FindValueResult_DataObject{
				DataObject: &pb.DataObject{Data: dataObject.Data, OriginalUploader: utils.ContactToPbNode(dataObject.Contact)},
			},
		}, nil
	}

	payload := &pb.FindNodeRequest{
		TargetID:       request.Hash,
		RequestingNode: &pb.Node{ID: s.id.Bytes(), IPWithPort: s.address},
	}

	nodes, errFindNodes := s.FindNode(ctx, payload)

	if errFindNodes != nil {
		err = fmt.Errorf("find Node error in Find value with error: %s", errFindNodes)
		return nil, err
	}

	return &pb.FindValueResult{
		Value: &pb.FindValueResult_Nodes{
			Nodes: nodes,
		},
	}, nil
}
