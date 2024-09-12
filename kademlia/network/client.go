package network

import (
	"d7024e_group04/kademlia/contact"
	pb "d7024e_group04/proto"
)

type Client struct {
	pb.KademliaClient
}

// TODO: Lets start with a simulated network
// TODO: Define a network interface to have a simulated on and a real one (and maybe a spy test one)

func Listen(address string) {
	// TODO
}

func (c *Client) SendPingMessage(contact *contact.Contact) {
	// TODO
}

func (c *Client) SendFindContactMessage(contact *contact.Contact) {
	// TODO
}

func (c *Client) SendFindDataMessage(hash string) {
	// TODO
}

func (c *Client) SendStoreMessage(data []byte) {
	// TODO
}
