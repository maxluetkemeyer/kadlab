package store

type Store interface {
	SetValue(key string, value []byte)
	GetValue(key string) (value []byte, error error)
}
