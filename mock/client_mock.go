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
	me                *contact.Contact
	pingSuccessful    bool
	findNodeSuccesful bool
	nodes             map[string]*MockNode
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
	}, nil
}

func (c *mockClient) SetPingResult(result bool) {
	c.pingSuccessful = result
}

func (c *mockClient) SetFindNodeResult(result bool) {
	c.findNodeSuccesful = result
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
	if !c.findNodeSuccesful {
		return nil, fmt.Errorf("findNodeSuccesful is false")
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
