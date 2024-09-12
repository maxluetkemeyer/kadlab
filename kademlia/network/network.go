package network

import (
	"context"
	"d7024e_group04/kademlia/contact"
)

type NetworkClient interface {
	SendPing(ctx context.Context, contact *contact.Contact) error
}
