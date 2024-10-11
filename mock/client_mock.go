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

func (c *MockClient) SendPing(_ context.Context, _ string) (*contact.Contact, error) {
	if c.pingSuccessful {
		return c.me, nil
	}

	return nil, fmt.Errorf("failed to ping")
}

func (c *MockClient) SendFindNode(_ context.Context, _, _ *contact.Contact) ([]*contact.Contact, error) {
	panic("TODO")
}

func (c *MockClient) SendFindValue(_ context.Context, _ *contact.Contact, _ string) ([]*contact.Contact, string, error) {
	panic("TODO")
}

func (c *MockClient) SendStore(_ context.Context, _ *contact.Contact, _ string) error {
	panic("TODO")
}
