package store

import (
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/model"
	"time"
)

type Store interface {
	SetValue(key string, value string, uploader *contact.Contact)
	GetValue(key string) (dataObject model.DataWithOriginalUploader, err error)
	GetOriginalUploader(key string) (*contact.Contact, error)
}

type TTLStore interface {
	SetValue(key string, value string, ttl time.Duration, uploader *contact.Contact)
	GetValue(key string) (dataObject model.DataWithOriginalUploader, err error)
	SetTTL(key string, ttl time.Duration)
	GetTTL(key string) time.Duration
	GetOriginalUploader(key string) (*contact.Contact, error)
	GetStoreLocations(key string) []*contact.Contact
	AddStoreLocation(key string, contact *contact.Contact)
	RemoveRefreshContact(key string)
}
