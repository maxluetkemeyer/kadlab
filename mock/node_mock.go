package mock

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/model"
)

type MockNode struct {
	me    *contact.Contact
	Store map[string]string
}

/*
func NewNodeMock(me *contact.Contact) *MockNode {
	return &MockNode{
		me:    me,
		Store: make(map[string]string),
	}
}*/

func (n *MockNode) Me() *contact.Contact {
	return n.me
}

func (n *MockNode) Bootstrap(_ context.Context) error {
	panic("TODO")
}

func (n *MockNode) PutObject(_ context.Context, _ string) (hashAsHex string, err error) {
	panic("TODO")
}

func (n *MockNode) GetObject(_ context.Context, hash string) (FindValueSuccessfulResponse *model.FindValueSuccessfulResponse, candidates []*contact.Contact, err error) {
	if data, found := n.Store[hash]; found {
		return &model.FindValueSuccessfulResponse{DataValue: data, NodeWithValue: n.me}, nil, nil
	}

	fakeContact := contact.NewContact(kademliaid.NewRandomKademliaID(), "address")

	return nil, []*contact.Contact{fakeContact}, nil
}
