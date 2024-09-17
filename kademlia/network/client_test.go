package network

import (
	"context"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/mock"
	"d7024e_group04/kademlia/routingtable"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestClient_Ping(t *testing.T) {
	mock.InitBufconn()

	routingTable := routingtable.NewRoutingTable(contact.NewContact(mock.ClientID, mock.ClientAddress))

	client := NewClient(mock.ClientAddress, mock.ClientID, routingTable, grpc.WithContextDialer(mock.BufDialer))

	t.Run("ping valid node", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		go mock.TimeoutContext(ctx, cancel)

		target := contact.NewContact(mock.TargetID, "passthrough://bufnet")

		candidates := routingTable.FindClosestContacts(mock.TargetID, 1)

		if len(candidates) != 0 {
			t.Fatalf("target exists in routing table")
		}

		if err := client.SendPing(ctx, &target); err != nil {
			t.Fatalf("failed to ping, err: %v", err)
		}

		candidates = routingTable.FindClosestContacts(mock.TargetID, 1)

		if len(candidates) != 1 {
			t.Fatal("target was not added to routing table")
		}

		if !candidates[0].ID.Equals(mock.TargetID) {
			t.Fatalf("candidate is not target, expected: %v got: %v", mock.TargetID, candidates[0].ID)
		}

	})
}
