package client

import (
	"context"
	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	pb "d7024e_group04/proto"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	me   *contact.Contact
	opts []grpc.DialOption
}

// NewClient returns a client with optional dial options
func NewClient(me *contact.Contact, opts ...grpc.DialOption) *Client {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &Client{
		me:   me,
		opts: opts,
	}
}

// initConnection returns a grpc connection and client to the target address
// It is callers responsibility to close the connection after use to prevent leakage
func (c *Client) connectTo(address string, opts ...grpc.DialOption) (*grpc.ClientConn, pb.KademliaClient, error) {
	opts = append(opts, c.opts...)
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create connection to address: %v, err: %v", address, err)
	}

	kademliaClient := pb.NewKademliaClient(conn)

	return conn, kademliaClient, nil
}

// SendPingMessage sends a Ping rpc call to the target. Returns the contact of target node.
func (c *Client) SendPing(ctx context.Context, targetIpWithPort string) (*contact.Contact, error) {
	conn, kademliaClient, err := c.connectTo(targetIpWithPort)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	payload := contactToPbNode(c.me)

	responseNode, err := kademliaClient.Ping(ctx, payload)

	if err != nil {
		return nil, fmt.Errorf("failed to ping address %v, err: %v", targetIpWithPort, err)
	}

	contact := pbNodeToContact(responseNode)

	return contact, nil
}

// SendFindNode sends a FindNode rpc call to the candidate node. Returns the closest nodes to target node from the candidate.
func (c *Client) SendFindNode(ctx context.Context, contactWeRequest, contactWeAreSearchingFor *contact.Contact) ([]*contact.Contact, error) {
	conn, kademliaClient, err := c.connectTo(contactWeRequest.Address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	payload := &pb.FindNodeRequest{
		TargetID:       contactWeAreSearchingFor.ID.Bytes(),
		RequestingNode: contactToPbNode(c.me),
	}

	resp, err := kademliaClient.FindNode(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to send FIND_NODE RPC to address: %v, err: %v", contactWeRequest.Address, err)
	}

	contacts := make([]*contact.Contact, 0, len(resp.Nodes))
	for _, node := range resp.Nodes {
		if len(node.ID) != env.IDLength {
			log.Printf("%v\n", err)
			continue
		}
		id := kademliaid.NewKademliaIDFromBytes([env.IDLength]byte(node.ID))

		newContact := contact.NewContact(id, node.IPWithPort)
		contacts = append(contacts, newContact)
	}

	return contacts, nil
}

// SendFindValue sends a FindNode rpc call to the target contact with a hash value.
// Returns the data if it is found on the target contact, otherwise a slice of candidate nodes closest to the hash value.
func (c *Client) SendFindValue(ctx context.Context, contactWeRequest *contact.Contact, hash string) (candidates []*contact.Contact, data string, err error) {
	conn, kademliaClient, err := c.connectTo(contactWeRequest.Address)
	if err != nil {
		return nil, "", err
	}
	defer conn.Close()

	payload := &pb.FindValueRequest{
		Hash:           kademliaid.NewKademliaID(hash).Bytes(),
		RequestingNode: contactToPbNode(c.me),
	}

	resp, err := kademliaClient.FindValue(ctx, payload)

	if err != nil {
		return nil, "", fmt.Errorf("rpc server returned err: %v", err)
	}

	switch respValue := resp.Value.(type) {
	case *pb.FindValueResult_Data:
		return nil, respValue.Data, nil

	case *pb.FindValueResult_Nodes:
		candidates := make([]*contact.Contact, 0, len(respValue.Nodes.Nodes))
		for _, node := range respValue.Nodes.Nodes {
			candidates = append(candidates, pbNodeToContact(node))
		}

		return candidates, "", nil

	default:
		return nil, "", fmt.Errorf("response type invalid, resp: %v", respValue)
	}
}

func (c *Client) SendStore(ctx context.Context, contactWeRequest *contact.Contact, data string) error {
	conn, kademliaClient, err := c.connectTo(contactWeRequest.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// key is hash of data
	key := kademliaid.NewKademliaIDFromData(data)

	storeRequest := &pb.StoreRequest{
		Key:            key.Bytes(),
		Value:          data,
		RequestingNode: contactToPbNode(c.me),
	}

	storeResult, err := kademliaClient.Store(ctx, storeRequest)

	if err != nil || !storeResult.Success {
		return err
	}

	return nil
}

func pbNodeToContact(node *pb.Node) *contact.Contact {
	return contact.NewContact((kademliaid.KademliaID)(node.ID), node.IPWithPort)
}

func contactToPbNode(contact *contact.Contact) *pb.Node {
	return &pb.Node{ID: contact.ID.Bytes(), IPWithPort: contact.Address}
}
