package model

import (
	"d7024e_group04/internal/kademlia/contact"
	"time"
)

type FindValueSuccessfulResponse struct {
	DataValue        string           `json:"Value"`
	NodeWithValue    *contact.Contact `json:"Node"`
	OriginalUploader *contact.Contact `json:"OriginalUploader"`
}

type DataWithOriginalUploader struct {
	Data    string
	Contact *contact.Contact
}
type RefreshTTLRequest struct {
	Key string
	TTL time.Duration
}
