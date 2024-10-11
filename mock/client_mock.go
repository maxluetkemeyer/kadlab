package mock

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"fmt"
)

type MockClient struct {
	me             *contact.Contact
	pingSuccessful bool
}

func NewClientMock(me *contact.Contact) *MockClient {
	return &MockClient{
		me: me,
	}
}

func (c *MockClient) SetPingResult(result bool) {
	c.pingSuccessful = result
}

func (c *MockClient) SendPing(ctx context.Context, targetIpWithPort string) (*contact.Contact, error) {
	if c.pingSuccessful {
		return c.me, nil
	}

	return nil, fmt.Errorf("failed to ping")
}

func (c *MockClient) SendFindNode(ctx context.Context, contactWeRequest, contactWeAreSearchingFor *contact.Contact) ([]*contact.Contact, error) {
	panic("TODO")
}

func (c *MockClient) SendFindValue(ctx context.Context, contactWeRequest *contact.Contact, hash string) ([]*contact.Contact, string, error) {
	panic("TODO")
}

func (c *MockClient) SendStore(ctx context.Context, contactWeRequest *contact.Contact, data string) error {
	panic("TODO")
}
