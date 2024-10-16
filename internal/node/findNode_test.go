package node

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/model"
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
	store        store.TTLStore
	contact      *contact.Contact
	routingTable *routingtable.RoutingTable
}

func newTestNode(me *contact.Contact, contacts []*contact.Contact) *TestNode {
	routingTable := routingtable.NewRoutingTable(me)

	for _, c := range contacts {
		routingTable.AddContact(c)
	}

	memoryStore := store.NewMemoryStore()
	simpleTtlStore := store.NewSimpleTTLStore(memoryStore)

	return &TestNode{
		server:       server.NewServer(routingTable, simpleTtlStore),
		store:        simpleTtlStore,
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

	testNodes[thirteen.Address] = newTestNode(thirteen, []*contact.Contact{
		one, four, five, twelve, fifteen,
	})

	testNodes[fifteen.Address] = newTestNode(fifteen, []*contact.Contact{
		one, four, five, twelve, thirteen,
	})

	testNodes[eighteen.Address] = newTestNode(eighteen, []*contact.Contact{
		one, four, five,
	})

	return testNodes
}

type ClientMock struct {
	me                     *contact.Contact
	testNodes              map[string]*TestNode
	findNodeCountUntilFail int
	findNodeSuccesfulCount atomic.Uint32
}

func newClientMock(testNodes map[string]*TestNode, me *contact.Contact) *ClientMock {
	return &ClientMock{
		me:        me,
		testNodes: testNodes,
	}
}

// Set the number of requests that will success until one fails
// Eg. Setting it to 1 will make all requests fail, while
//
//	setting it to 3 will make 2 requests work and 1 fail.
//
// Setting it to 0 will disable this feature
func (c *ClientMock) SetFindNodeSuccesfulCount(count int) {
	c.findNodeCountUntilFail = count
}

func (c *ClientMock) SendPing(ctx context.Context, targetIpWithPort string) (*contact.Contact, error) {
	return nil, fmt.Errorf("should not be used")
}

func (c *ClientMock) SendFindNode(ctx context.Context, contactWeRequest, contactWeAreSearchingFor *contact.Contact) ([]*contact.Contact, error) {
	if c.findNodeCountUntilFail != 0 && int(c.findNodeSuccesfulCount.Load()) >= c.findNodeCountUntilFail-1 {
		c.findNodeSuccesfulCount.Store(0)
		return nil, fmt.Errorf("bad network (not a real error)")
	} else {
		c.findNodeSuccesfulCount.Add(1)
	}

	candidateNode := c.testNodes[contactWeRequest.Address]
	return candidateNode.routingTable.FindClosestContacts(contactWeAreSearchingFor.ID, env.BucketSize), nil
}

func (c *ClientMock) SendFindValue(ctx context.Context, contactWeRequest *contact.Contact, hash string) (candidates []*contact.Contact, dataObject model.FindValueSuccessfulResponse, err error) {
	candidateNode := c.testNodes[contactWeRequest.Address]

	value, err := candidateNode.store.GetValue(hash)
	if err != nil {
		hashKademliaID := kademliaid.NewKademliaID(hash)
		return candidateNode.routingTable.FindClosestContacts(hashKademliaID, env.BucketSize), dataObject, nil
	}
	originalUploader, _ := candidateNode.store.GetOriginalUploader(hash)
	dataObject = model.FindValueSuccessfulResponse{
		DataValue:        value,
		NodeWithValue:    candidateNode.contact,
		OriginalUploader: originalUploader,
	}
	return nil, dataObject, nil
}

func (c *ClientMock) SendStore(ctx context.Context, contactWeRequest *contact.Contact, data string, originalUploader *contact.Contact) error {
	candidateNode := c.testNodes[contactWeRequest.Address]
	key := kademliaid.NewKademliaIDFromData(data)
	candidateNode.store.SetValue(key.String(), data, time.Hour, originalUploader)
	return nil
}

func (c *ClientMock) SendRefreshTTL(ctx context.Context, key string, contactWeRequest *contact.Contact) error {
	panic("TODO")
}

func (c *ClientMock) SendNewStoredLocation(ctx context.Context, key string, originalUploader, newContactStoringData *contact.Contact) error {
	candidateNode := c.testNodes[originalUploader.Address]
	candidateNode.store.AddStoreLocation(key, newContactStoringData)
	return nil
}

func TestFindNode(t *testing.T) {
	// Simplify testing
	env.BucketSize = 4
	testNodes := populateTestNodes()

	t.Run("findNode", func(t *testing.T) {
		expectedNodes := []*contact.Contact{thirteen, twelve, fifteen, five}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// We are 18
		node := Node{
			RoutingTable: testNodes[":18"].routingTable,
			Client:       newClientMock(testNodes, testNodes[":18"].contact),
		}

		// Trying to find 13
		nodesFound := node.findNode(ctx, thirteen)

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

		mockClient := newClientMock(testNodes, eighteen)
		mockClient.SetFindNodeSuccesfulCount(1)

		// We are 18
		node := Node{
			RoutingTable: testNodes[eighteen.Address].routingTable,
			Client:       mockClient,
		}

		// Trying to find 13
		nodesFound := node.findNode(ctx, thirteen)

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

		mockClient := newClientMock(testNodes, eighteen)
		mockClient.SetFindNodeSuccesfulCount(3)

		// We are 18
		node := Node{
			RoutingTable: testNodes[eighteen.Address].routingTable,
			Client:       mockClient,
		}

		// Trying to find 13
		nodesFound := node.findNode(ctx, thirteen)

		if len(nodesFound) < 3 {
			t.Fatalf("wrong number of nodes, expected at least 3, got %v", nodesFound)
		}

	})

}
