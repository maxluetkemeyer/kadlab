package store

import "time"

type Store interface {
	SetValue(key string, value string)
	GetValue(key string) (value string, err error)
}

type TTLStore interface {
	SetValue(key string, value string, ttl time.Duration)
	GetValue(key string) (value string, error error)
	SetTTL(key string, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
}
