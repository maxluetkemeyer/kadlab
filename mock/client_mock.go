package mock

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"fmt"
)

type mockClient struct {
	me             *contact.Contact
	pingSuccessful bool
}

func NewClientMock(me *contact.Contact) *mockClient {
	return &mockClient{
		me: me,
	}
}

func (c *mockClient) SetPingResult(result bool) {
	c.pingSuccessful = result
}

func (c *mockClient) SendPing(ctx context.Context, targetIpWithPort string) (*contact.Contact, error) {
	if c.pingSuccessful {
		return c.me, nil
	}

	return nil, fmt.Errorf("failed to ping")
}

func (c *mockClient) SendFindNode(ctx context.Context, contactWeRequest, contactWeAreSearchingFor *contact.Contact) ([]*contact.Contact, error) {
	panic("TODO")
}

func (c *mockClient) SendFindValue(ctx context.Context, contactWeRequest *contact.Contact, hash string) ([]*contact.Contact, string, error) {
	panic("TODO")
}

func (c *mockClient) SendStore(ctx context.Context, contactWeRequest *contact.Contact, data string) error {
	panic("TODO")
}
func (c *mockClient) SendRefreshTTL(ctx context.Context, key string, contactWeRequest *contact.Contact) error {
	panic("TODO")
}
