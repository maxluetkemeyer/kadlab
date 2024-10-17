package utils

import (
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	pb "d7024e_group04/proto"
)

func PbNodeToContact(node *pb.Node) *contact.Contact {
	return contact.NewContact((kademliaid.KademliaID)(node.ID), node.IPWithPort)
}

func ContactToPbNode(contact *contact.Contact) *pb.Node {
	return &pb.Node{ID: contact.ID.Bytes(), IPWithPort: contact.Address}
}
