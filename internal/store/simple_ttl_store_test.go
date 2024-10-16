package store

import (
	"log"
	"testing"
	"time"
)

func TestSimpleTTLStore(t *testing.T) {
	t.Run("set value with enough time", func(t *testing.T) {
		memStore := NewMemoryStore()
		tStore := NewSimpleTTLStore(memStore)

		want := "value0"

		tStore.SetValue("key0", want, time.Hour, nil)

		got, err := tStore.GetValue("key0")

		if err != nil {
			log.Fatalf("error getting value: %v", err)
		}

		if got != want {
			log.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("set value with zero time", func(t *testing.T) {
		memStore := NewMemoryStore()
		tStore := NewSimpleTTLStore(memStore)

		val := "value0"

		tStore.SetValue("key0", val, time.Second*0, nil)

		got, err := tStore.GetValue("key0")

		if got != "" {
			log.Fatalf("should get empty string, but got %v", got)
		}

		if err == nil {
			log.Fatalf("should error, but got nil")
		}

	})

	t.Run("set value with enough ttl and later it should be expired", func(t *testing.T) {
		memStore := NewMemoryStore()
		tStore := NewSimpleTTLStore(memStore)

		val := "value0"

		tStore.SetValue("key0", val, time.Hour, nil)

		got, err := tStore.GetValue("key0")

		if err != nil {
			log.Fatalf("error getting value: %v", err)
		}

		if got != val {
			log.Fatalf("TTL should be enough: got %q, want %q", got, val)
		}

		tStore.SetTTL("key0", time.Second*0)

		gotAfterExpiration, errAfterExpiration := tStore.GetValue("key0")

		if gotAfterExpiration != "" {
			log.Fatalf("should get empty string, but got %v", gotAfterExpiration)
		}

		if errAfterExpiration == nil {
			log.Fatalf("should error, but got nil")
		}
	})
}
