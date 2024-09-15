package kademlia

import (
	"context"
	"d7024e_group04/kademlia/network"
	"golang.org/x/sync/errgroup"
)

type Node struct {
	Client network.ClientRPC
	server network.ServerRPC
}

func NewNode(client network.ClientRPC, server network.ServerRPC) *Node {
	return &Node{
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
