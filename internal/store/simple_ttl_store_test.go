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

		tStore.SetValue("key0", want, time.Hour)

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

		tStore.SetValue("key0", val, time.Second*0)

		got, err := tStore.GetValue("key0")

		if got == val {
			log.Fatalf("should get empty string, but got %v", got)
		}

		if err == nil {
			log.Fatalf("should error, but got nil")
		}

	})
}
