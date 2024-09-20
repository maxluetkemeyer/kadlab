package client

import (
	"context"
	"fmt"
	"log"

	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	pb "d7024e_group04/proto"

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
func (c *Client) SendPing(ctx context.Context, me *contact.Contact, target string) (contact contact.Contact, err error) {
	conn, err := initConnection(target, c.opts...)
	if err != nil {
		return contact, fmt.Errorf("failed to create connection address: %v, err: %v", target, err)
	}

	defer conn.Close()

	grpc := pb.NewKademliaClient(conn)

	payload := &pb.Node{
		ID:         &pb.KademliaID{Value: me.ID.Bytes()},
		IPWithPort: me.Address,
	}

	resp, err := grpc.Ping(ctx, payload)

	if err != nil {
		return contact, fmt.Errorf("failed to ping address %v, err: %v", target, err)
	}

	contact = pbNodeToContact(resp)

	return contact, nil
}

func (c *Client) SendFindNode(ctx context.Context, target *contact.Contact) ([]contact.Contact, error) {
	conn, err := initConnection(target.Address, c.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to address: %v, err: %v", target.Address, err)
	}

	defer conn.Close()

	grpc := pb.NewKademliaClient(conn)

	payload := &pb.KademliaID{
		Value: target.ID.Bytes(),
	}

	resp, err := grpc.Find_Node(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to send FIND_NODE RPC to address: %v, err: %v", target.Address, err)
	}

	var contacts []contact.Contact
	for _, node := range resp.Node {
		id, err := kademliaid.NewKademliaIDFromBytes(node.ID.Value)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			continue
		}

		newContact := contact.NewContact(id, node.IPWithPort)
		contacts = append(contacts, newContact)
	}

	return contacts, nil
}

func (c *Client) SendFindValue(ctx context.Context, hash string) (string, error) {
	panic("TODO")
}

func (c *Client) SendStore(ctx context.Context, data string) error {
	panic("TODO")
}

func pbNodeToContact(node *pb.Node) contact.Contact {
	return contact.NewContact((*kademliaid.KademliaID)(node.ID.Value), node.IPWithPort)
}
