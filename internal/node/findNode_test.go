package node_test

import (
	"context"
	"testing"
	"time"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/node"
	"d7024e_group04/mock"
)

var findNode = node.ExportFindNode

/*
	1  0x0400000000000000000000000000000000000000
	4  0x1000000000000000000000000000000000000000
	5  0x1400000000000000000000000000000000000000
	12 0x3000000000000000000000000000000000000000
	13 0x3400000000000000000000000000000000000000
	15 0x3c00000000000000000000000000000000000000
	18 0x4800000000000000000000000000000000000000
*/

var (
	one      = contact.NewContact(kademliaid.NewKademliaID("0400000000000000000000000000000000000000"), ":1")
	four     = contact.NewContact(kademliaid.NewKademliaID("1000000000000000000000000000000000000000"), ":4")
	five     = contact.NewContact(kademliaid.NewKademliaID("1400000000000000000000000000000000000000"), ":5")
	twelve   = contact.NewContact(kademliaid.NewKademliaID("3000000000000000000000000000000000000000"), ":12")
	thirteen = contact.NewContact(kademliaid.NewKademliaID("3400000000000000000000000000000000000000"), ":13")
	fifteen  = contact.NewContact(kademliaid.NewKademliaID("3c00000000000000000000000000000000000000"), ":15")
	eighteen = contact.NewContact(kademliaid.NewKademliaID("4800000000000000000000000000000000000000"), ":18")
)

func populateTestNodes() map[string]*mock.MockNode {
	testNodes := make(map[string]*mock.MockNode, 7)

	testNodes[one.Address] = mock.NewNodeMock(one, []*contact.Contact{
		twelve, thirteen, fifteen, five, four,
	})

	testNodes[four.Address] = mock.NewNodeMock(four, []*contact.Contact{
		five, twelve, thirteen, fifteen,
	})

	testNodes[five.Address] = mock.NewNodeMock(five, []*contact.Contact{
		four, twelve, thirteen, fifteen,
	})

	testNodes[twelve.Address] = mock.NewNodeMock(twelve, []*contact.Contact{
		one, four, five, twelve, thirteen,
	})

	testNodes[thirteen.Address] = mock.NewNodeMock(one, []*contact.Contact{
		one, four, five, twelve, fifteen,
	})

	testNodes[fifteen.Address] = mock.NewNodeMock(one, []*contact.Contact{
		one, four, five, twelve, thirteen,
	})

	testNodes[eighteen.Address] = mock.NewNodeMock(one, []*contact.Contact{
		one, four, five,
	})

	return testNodes
}

func TestFindNode(t *testing.T) {
	// Simplify testing
	env.BucketSize = 4
	testNodes := populateTestNodes()

	t.Run("findNode with working network", func(t *testing.T) {
		expectedNodes := []*contact.Contact{thirteen, twelve, fifteen, five}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// We are 18
		client, err := mock.NewClientMockWithNodes(eighteen.Address, testNodes)
		if err != nil {
			t.Fatal(err)
		}

		me, err := client.GetNode(eighteen.Address)
		if err != nil {
			t.Fatal(err)
		}

		client.SetFindNodeSuccesfulCount(0)

		// Trying to find 13
		nodesFound := findNode(me, ctx, thirteen)

		// Expecting 13,12,15,5
		if len(nodesFound) != len(expectedNodes) {
			t.Fatalf("wrong number of nodes, expected %v, got %v", expectedNodes, nodesFound)
		}

		for i, node := range nodesFound {
			if !node.ID.Equals(expectedNodes[i].ID) {
				t.Fatalf("wrong nodes, expected %v, got %v", expectedNodes, nodesFound)
			}
		}
	})

	t.Run("findNode with bad network", func(t *testing.T) {
		expectedNodes := []*contact.Contact{}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// We are 18
		client, err := mock.NewClientMockWithNodes(eighteen.Address, testNodes)
		if err != nil {
			t.Fatal(err)
		}

		me, err := client.GetNode(eighteen.Address)
		if err != nil {
			t.Fatal(err)
		}

		client.SetFindNodeSuccesfulCount(1)

		// Trying to find 13
		nodesFound := findNode(me, ctx, thirteen)

		if len(nodesFound) != len(expectedNodes) {
			t.Fatalf("wrong number of nodes, expected %v, got %v", expectedNodes, nodesFound)
		}

		for i, node := range nodesFound {
			if !node.ID.Equals(expectedNodes[i].ID) {
				t.Fatalf("wrong nodes, expected %v, got %v", expectedNodes, nodesFound)
			}
		}
	})

	t.Run("findNode with faulty network", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// We are 18
		client, err := mock.NewClientMockWithNodes(eighteen.Address, testNodes)
		if err != nil {
			t.Fatal(err)
		}

		me, err := client.GetNode(eighteen.Address)
		if err != nil {
			t.Fatal(err)
		}

		client.SetFindNodeSuccesfulCount(3)

		// Trying to find 13
		nodesFound := findNode(me, ctx, thirteen)

		if len(nodesFound) < 3 {
			t.Fatalf("wrong number of nodes, expected at least 3, got %v", nodesFound)
		}

	})

}
