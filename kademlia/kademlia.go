package kademlia

import (
	"context"
	"d7024e_group04/kademlia/network"

	"golang.org/x/sync/errgroup"
)

type Node struct {
	store  map[string][]byte
	Client network.ClientRPC
	server network.ServerRPC
}

func NewNode(address string, client network.ClientRPC, server network.ServerRPC) *Node {
	return &Node{
		store:  make(map[string][]byte),
		Client: client,
		server: server,
	}
}

func (n *Node) Start(ctx context.Context) error {
	// TODO eventually start more stuff here that should be running during runtime
	errGroup, errCtx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return n.server.Start(errCtx)
	})

	return errGroup.Wait()
}
