package node

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/kademliaid"
	"testing"
	"time"
)

func initBucketRefreshTest() (context.Context, Node) {
	env.TRefresh = 5 * time.Second
	testNodes := populateTestNodes()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// We are 18
	node := Node{
		RoutingTable: testNodes[eighteen.Address].routingTable,
		Client:       newClientMock(testNodes, eighteen),
	}

	return ctx, node
}

func TestBucketRefresh(t *testing.T) {
	// Simplify testing
	env.BucketSize = 4

	t.Run("bucket refresh", func(t *testing.T) {
		ctx, node := initBucketRefreshTest()

		// Wait refresh time
		time.Sleep(env.TRefresh + time.Millisecond)
	})

}
