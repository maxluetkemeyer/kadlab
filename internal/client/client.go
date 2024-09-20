package client

import (
	"context"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	pb "d7024e_group04/proto"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb.KademliaClient
	opts []grpc.DialOption
}

// TODO: Lets start with a simulated network
// TODO: Define a network interface to have a simulated on and a real one (and maybe a spy test one)

func NewClient(opts ...grpc.DialOption) *Client {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &Client{
		opts: opts,
	}
}

// initConnection returns a grpc connection and client to the target address
// It is callers responsibility to close the connection after use to prevent leakage
func (c *Client) NewConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, c.opts...)
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// SendPingMessage sends an rpc call to the target contact. If a reply is received the bucket is updated with the target contact.
func (c *Client) SendPing(ctx context.Context, grpc pb.KademliaClient, me, target *contact.Contact) (contact contact.Contact, err error) {
	payload := &pb.Node{
		ID:         me.ID.Bytes(),
		IPWithPort: me.Address,
	}

	resp, err := grpc.Ping(ctx, payload)

	if err != nil {
		return contact, fmt.Errorf("failed to ping address %v, err: %v", target.Address, err)
	}

	contact = pbNodeToContact(resp)

	return contact, nil
}

func (c *Client) SendFindNode(ctx context.Context, contact contact.Contact) ([]contact.Contact, error) {
	panic("TODO")
}

func (c *Client) SendFindValue(ctx context.Context, me, target contact.Contact, hash string) (candidates *contact.ContactCandidates, data string, err error) {
	conn, err := c.NewConnection(target.Address)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create connection to address: %v, err: %v", target.Address, err)
	}

	defer conn.Close()

	grpc := pb.NewKademliaClient(conn)

	payload := &pb.FindValueRequest{
		Hash:     kademliaid.NewKademliaID(hash).Bytes(),
		SenderID: me.ID.Bytes(),
	}

	resp, err := grpc.FindValue(ctx, payload)

	if err != nil {
		return nil, "", fmt.Errorf("rpc server returned err: %v", err)
	}

	switch value := resp.Value.(type) {
	case *pb.NodesOrData_Data:
		return nil, value.Data, nil

	case *pb.NodesOrData_Nodes:
		contacts := make([]contact.Contact, 0, len(value.Nodes.Node))
		for _, node := range value.Nodes.Node {
			contacts = append(contacts, pbNodeToContact(node))
		}
		candidates.Append(contacts)

		return candidates, "", nil

	default:
		return nil, "", fmt.Errorf("response type invalid")
	}
}

func (c *Client) SendStore(ctx context.Context, data string) error {
	panic("TODO")
}

func pbNodeToContact(node *pb.Node) contact.Contact {
	return contact.NewContact((*kademliaid.KademliaID)(node.ID), node.IPWithPort)
}
