package mock

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/model"
)

type mockNode struct {
	me    *contact.Contact
	Store map[string]string
}

func NewNodeMock(me *contact.Contact) *mockNode {
	return &mockNode{
		me:    me,
		Store: make(map[string]string),
	}
}

func (n *mockNode) Me() *contact.Contact {
	return n.me
}

func (n *mockNode) Bootstrap(rootCtx context.Context) error {
	panic("TODO")
}

func (n *mockNode) PutObject(ctx context.Context, data string) (hashAsHex string, err error) {
	panic("TODO")
}

func (n *mockNode) GetObject(rootCtx context.Context, hash string) (FindValueSuccessfulResponse *model.FindValueSuccessfulResponse, candidates []*contact.Contact, err error) {
	if data, found := n.Store[hash]; found {
		return &model.FindValueSuccessfulResponse{DataValue: data, NodeWithValue: n.me}, nil, nil
	}

	fakeContact := contact.NewContact(kademliaid.NewRandomKademliaID(), "address")

	return nil, []*contact.Contact{fakeContact}, nil
}
