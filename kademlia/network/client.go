package network

import (
	"context"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	"d7024e_group04/kademlia/routingtable"
	pb "d7024e_group04/proto"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb.KademliaClient
	id           *kademliaid.KademliaID
	address      string
	routingTable *routingtable.RoutingTable
	opts         []grpc.DialOption
}

// TODO: Lets start with a simulated network
// TODO: Define a network interface to have a simulated on and a real one (and maybe a spy test one)

func NewClient(address string, id *kademliaid.KademliaID, routingTable *routingtable.RoutingTable, opts ...grpc.DialOption) *Client {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &Client{
		id:           id,
		address:      address,
		routingTable: routingTable,
		opts:         opts,
	}
}

// initConnection returns a grpc connection to the target address
// It is callers responsibility to close the connection after use to prevent leakage
func initConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// SendPingMessage sends an rpc call to the target contact. If a reply is received the bucket is updated with the target contact.
func (c *Client) SendPing(ctx context.Context, contact *contact.Contact) error {
	conn, err := initConnection(contact.Address, c.opts...)
	if err != nil {
		return fmt.Errorf("failed to create connection to node: %v address: %v, err: %v", contact.ID, contact.Address, err)
	}

	defer conn.Close()

	grpc := pb.NewKademliaClient(conn)

	payload := &pb.Node{
		ID:      c.id.Bytes(),
		Address: c.address,
	}

	resp, err := grpc.Ping(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to ping node: %v address %v, err: %v", contact.ID, contact.Address, err)
	}

	if !contact.ID.Equals((*kademliaid.KademliaID)(resp.ID)) {
		return fmt.Errorf("invalid response from ping to node: %v, address %v, err: %v", contact.ID, contact.Address, err)
	}

	c.routingTable.AddContact(*contact)

	return nil
}

func (c *Client) SendFindContactMessage(ctx context.Context, contact *contact.Contact) {
	// TODO
}

func (c *Client) SendFindDataMessage(ctx context.Context, hash string) {
	// TODO
}

func (c *Client) SendStoreMessage(ctx context.Context, data []byte) {
	// TODO
}
