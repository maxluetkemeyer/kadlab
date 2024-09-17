package store

type Store interface {
	SetValue(key string, value string)
	GetValue(key string) (value string, error error)
}
