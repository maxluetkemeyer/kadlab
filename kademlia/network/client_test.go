package network

import (
	"context"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/routingtable"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestClient_Ping(t *testing.T) {
	initBufconn()

	routingTable := routingtable.NewRoutingTable(contact.NewContact(clientID, clientAddress))

	client := NewClient(clientAddress, clientID, routingTable, grpc.WithContextDialer(bufDialer))

	t.Run("ping valid node", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		go timeoutContext(ctx, cancel)

		target := contact.NewContact(targetID, "passthrough://bufnet")

		candidates := routingTable.FindClosestContacts(targetID, 1)

		if len(candidates) != 0 {
			t.Fatalf("target exists in routing table")
		}

		if err := client.SendPing(ctx, &target); err != nil {
			t.Fatalf("failed to ping, err: %v", err)
		}

		candidates = routingTable.FindClosestContacts(targetID, 1)

		if len(candidates) != 1 {
			t.Fatal("target was not added to routing table")
		}

		if !candidates[0].ID.Equals(targetID) {
			t.Fatalf("candidate is not target, expected: %v got: %v", targetID, candidates[0].ID)
		}

	})
}
