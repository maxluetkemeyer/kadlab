package utils

import (
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	pb "d7024e_group04/proto"
	"reflect"
	"testing"
)

func TestConverters_PbNodeToContact(t *testing.T) {
	originalContact := contact.NewContact(kademliaid.NewRandomKademliaID(), "address")

	pbNode := &pb.Node{ID: originalContact.ID.Bytes(), IPWithPort: originalContact.Address}

	convertedToContact := PbNodeToContact(pbNode)

	if !originalContact.ID.Equals(convertedToContact.ID) {
		t.Fatalf("mismatching id, expected %v, got %v", originalContact.ID, convertedToContact.ID)
	}

	if originalContact.Address != convertedToContact.Address {
		t.Fatalf("mismatching address, expected %v, got %v", originalContact.Address, convertedToContact.Address)
	}
}

func TestConverters_ContactToPbNode(t *testing.T) {
	originalContact := contact.NewContact(kademliaid.NewRandomKademliaID(), "address")

	pbNode := ContactToPbNode(originalContact)

	if !reflect.DeepEqual(pbNode.ID, originalContact.ID.Bytes()) {
		t.Fatalf("mismatching id, expected %v, got %v", originalContact.ID.Bytes(), pbNode.ID)
	}

	if originalContact.Address != pbNode.IPWithPort {
		t.Fatalf("mismatching address, expected %v, got %v", originalContact.Address, pbNode.IPWithPort)
	}
}
