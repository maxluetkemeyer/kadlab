package mock

import (
	"context"
	"fmt"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/node"
	"d7024e_group04/internal/server"
	"d7024e_group04/internal/store"
)

type MockNode struct {
	server       *server.Server
	contact      *contact.Contact
	routingTable *routingtable.RoutingTable
}

type mockClient struct {
	me                     *contact.Contact
	pingSuccessful         bool
	findNodeCountUntilFail int
	findNodeSuccesfulCount int
	nodes                  map[string]*MockNode
}

func NewNodeMock(me *contact.Contact, contacts []*contact.Contact) *MockNode {
	routingTable := routingtable.NewRoutingTable(me)

	for _, c := range contacts {
		routingTable.AddContact(c)
	}

	return &MockNode{
		server:       server.NewServer(routingTable, store.NewMemoryStore()),
		contact:      me,
		routingTable: routingTable,
	}
}

func NewClientMock(me *contact.Contact) *mockClient {
	return &mockClient{
		me: me,
	}
}

func NewClientMockWithNodes(nodeAddr string, nodes map[string]*MockNode) (*mockClient, error) {
	n, ok := nodes[nodeAddr]
	if !ok {
		return nil, fmt.Errorf("Couldn't find node with addr=%v", nodeAddr)
	}

	return &mockClient{
		me:    n.contact,
		nodes: nodes,
		findNodeCountUntilFail: 0,
		findNodeSuccesfulCount: 0,
	}, nil
}

func (c *mockClient) SetPingResult(result bool) {
	c.pingSuccessful = result
}

// Set the number of requests that will success until one fails
// Eg. Setting it to 1 will make all requests fail, while
//     setting it to 3 will make 2 requests work and 1 fail.
// Setting it to 0 will disable this feature
func (c *mockClient) SetFindNodeSuccesfulCount(count int) {
	c.findNodeCountUntilFail = count
}

func (c *mockClient) GetNode(nodeAddr string) (*node.Node, error) {
	n, ok := c.nodes[nodeAddr]
	if !ok {
		return nil, fmt.Errorf("Couldn't find node with addr=%v", nodeAddr)
	}

	return &node.Node{
		Client:       c,
		RoutingTable: n.routingTable,
		Store:        store.NewMemoryStore(),
	}, nil
}

func (c *mockClient) SendPing(ctx context.Context, targetIpWithPort string) (*contact.Contact, error) {
	if c.pingSuccessful {
		return c.me, nil
	}

	return nil, fmt.Errorf("failed to ping")
}

func (c *mockClient) SendFindNode(ctx context.Context, contactWeRequest, contactWeAreSearchingFor *contact.Contact) ([]*contact.Contact, error) {
	if c.findNodeCountUntilFail != 0 && c.findNodeSuccesfulCount >= c.findNodeCountUntilFail-1 {
		c.findNodeSuccesfulCount = 0
		return nil, fmt.Errorf("bad network (not a real error)")
	} else {
		c.findNodeSuccesfulCount++
	}

	candidateNode := c.nodes[contactWeRequest.Address]
	return candidateNode.routingTable.FindClosestContacts(contactWeAreSearchingFor.ID, env.BucketSize), nil
}

func (c *mockClient) SendFindValue(ctx context.Context, contactWeRequest *contact.Contact, hash string) ([]*contact.Contact, string, error) {
	panic("TODO")
}

func (c *mockClient) SendStore(ctx context.Context, contactWeRequest *contact.Contact, data string) error {
	panic("TODO")
}
