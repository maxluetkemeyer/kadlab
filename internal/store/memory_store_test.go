package store

import "testing"

func TestMemoryStore(t *testing.T) {
	t.Run("Store a value and retrieve it", func(t *testing.T) {
		memStore := NewMemoryStore()

		key := "myKey"
		value := "myValue"
		memStore.SetValue(key, value, nil)

		want := value
		got, err := memStore.GetValue(key)

		if err != nil {
			t.Fatal(err)
		}

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}

	})

	t.Run("Store a value and retrieve a nonexistent key", func(t *testing.T) {
		memStore := NewMemoryStore()

		key := "myKey"
		value := "myValue"
		memStore.SetValue(key, value, nil)

		_, err := memStore.GetValue("anotherKey")

		if err == nil {
			t.Errorf("got nil, want error, %s", err)
		}
	})
}
