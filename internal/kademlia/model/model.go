package model

import "d7024e_group04/internal/kademlia/contact"

type ValueObject struct {
	DataValue     string
	NodeWithValue *contact.Contact
}
