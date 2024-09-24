package api

import (
	"context"
)

type KademliaAPI interface {
	ListenAndServe(ctx context.Context) error
}
