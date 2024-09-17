package store

import "testing"

func TestMemoryStore(t *testing.T) {
	t.Run("Store a value and retrieve it", func(t *testing.T) {
		memStore := NewMemoryStore()

		key := "myKey"
		value := "myValue"
		memStore.SetValue(key, []byte(value))

		want := value
		got, err := memStore.GetValue(key)

		if err != nil {
			t.Fatal(err)
		}

		if string(got) != want {
			t.Errorf("got %s, want %s", string(got), want)
		}

	})

	t.Run("Store a value and retrieve a nonexistent key", func(t *testing.T) {
		memStore := NewMemoryStore()

		key := "myKey"
		value := "myValue"
		memStore.SetValue(key, []byte(value))

		_, err := memStore.GetValue("anotherKey")

		if err == nil {
			t.Errorf("got nil, want error, %s", err)
		}
	})
}
