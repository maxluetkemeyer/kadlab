package node

import (
	"context"
	"fmt"
	"testing"
	"time"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/server"
	"d7024e_group04/internal/store"
)

/*
	1  0x0400000000000000000000000000000000000000
	4  0x1000000000000000000000000000000000000000
	5  0x1400000000000000000000000000000000000000
	12 0x3000000000000000000000000000000000000000
	13 0x3400000000000000000000000000000000000000
	15 0x3c00000000000000000000000000000000000000
	18 0x4800000000000000000000000000000000000000
*/

type TestNode struct {
	server       *server.Server
	contact      *contact.Contact
	routingTable *routingtable.RoutingTable
}

func newTestNode(me *contact.Contact, contacts []*contact.Contact) *TestNode {
	routingTable := routingtable.NewRoutingTable(me)

	for _, c := range contacts {
		routingTable.AddContact(*c)
	}

	return &TestNode{
		server:       server.NewServer(routingTable, store.NewMemoryStore()),
		contact:      me,
		routingTable: routingTable,
	}
}

var (
	one      = contact.NewContact(kademliaid.NewKademliaID("0400000000000000000000000000000000000000"), ":1")
	four     = contact.NewContact(kademliaid.NewKademliaID("1000000000000000000000000000000000000000"), ":4")
	five     = contact.NewContact(kademliaid.NewKademliaID("1400000000000000000000000000000000000000"), ":5")
	twelve   = contact.NewContact(kademliaid.NewKademliaID("3000000000000000000000000000000000000000"), ":12")
	thirteen = contact.NewContact(kademliaid.NewKademliaID("3400000000000000000000000000000000000000"), ":13")
	fifteen  = contact.NewContact(kademliaid.NewKademliaID("3c00000000000000000000000000000000000000"), ":15")
	eighteen = contact.NewContact(kademliaid.NewKademliaID("4800000000000000000000000000000000000000"), ":18")
)

func populateTestNodes() map[string]*TestNode {
	testNodes := make(map[string]*TestNode, 7)

	testNodes[one.Address] = newTestNode(one, []*contact.Contact{
		twelve, thirteen, fifteen, five, four,
	})

	testNodes[four.Address] = newTestNode(four, []*contact.Contact{
		five, twelve, thirteen, fifteen,
	})

	testNodes[five.Address] = newTestNode(five, []*contact.Contact{
		four, twelve, thirteen, fifteen,
	})

	testNodes[twelve.Address] = newTestNode(twelve, []*contact.Contact{
		one, four, five, twelve, thirteen,
	})

	testNodes[thirteen.Address] = newTestNode(one, []*contact.Contact{
		one, four, five, twelve, fifteen,
	})

	testNodes[fifteen.Address] = newTestNode(one, []*contact.Contact{
		one, four, five, twelve, thirteen,
	})

	testNodes[eighteen.Address] = newTestNode(one, []*contact.Contact{
		one, four, five,
	})

	return testNodes
}

type ClientMock struct{
	testNodes map[string]*TestNode
}

func newClientMock(testNodes map[string]*TestNode) *ClientMock {
	return &ClientMock{
		testNodes: testNodes,
	}
}

func (c *ClientMock) SendPing(ctx context.Context, me *contact.Contact, target string) (*contact.Contact, error) {
	return nil, fmt.Errorf("should not be used")
}

func (c *ClientMock) SendFindNode(ctx context.Context, me *contact.Contact, candidate string, targetId kademliaid.KademliaID) ([]contact.Contact, error) {
	candidateNode := c.testNodes[candidate]
	return candidateNode.routingTable.FindClosestContacts(targetId, env.BucketSize, me.ID), nil
}

func (c *ClientMock) SendFindValue(ctx context.Context, me, target contact.Contact, hash string) ([]contact.Contact, string, error) {
	return nil, "", fmt.Errorf("should not be used")
}

func (c *ClientMock) SendStore(ctx context.Context, data string) error {
	return fmt.Errorf("should not be used")
}

func TestFindNode(t *testing.T) {
	testNodes := populateTestNodes()

	t.Run("findNode", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// We are 18
		node := Node{
			RoutingTable: testNodes[":18"].routingTable,
			Client: newClientMock(testNodes),
		}

		// Trying to find 13
		nodesFound := node.findNode(ctx, thirteen)


		// Expecting 5,12,13,15
		t.Logf("nodes found = %v", nodesFound)
	})

}
