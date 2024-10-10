package model

import "d7024e_group04/internal/kademlia/contact"

type FindValueSuccessfulResponse struct {
	DataValue     string
	NodeWithValue *contact.Contact
}
