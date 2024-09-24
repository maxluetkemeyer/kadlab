package kademliaid

import (
	"bytes"
	. "d7024e_group04/env"
	"math/rand"
	"testing"
)

func TestKademliaID(t *testing.T) {
	randomBytes := generateRandomBytes(t, IDLength)

	t.Run("NewKademliaIDFromBytes and toBytes", func(t *testing.T) {
		goodBytes := [IDLength]byte{}
		//badBytes := [env.IDLength + 1]byte{}

		id := NewKademliaIDFromBytes(goodBytes)

		if len(id.Bytes()) != IDLength {
			t.Fatalf("Somehow the amount of bytes %v is different than the defined IDLength %v", len(id.Bytes()), IDLength)
		}
	})

	t.Run("to bytes", func(t *testing.T) {
		got := KademliaID(randomBytes)

		toBytes := got.Bytes()
		want := KademliaID(toBytes)

		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("to string", func(t *testing.T) {
		// Test is not configured for a different IDLength than 20!
		want := "1111111100000000000000000000000000000000"
		id := NewKademliaID(want)

		got := id.String()

		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("Equals", func(t *testing.T) {
		randomBytes1 := generateRandomBytes(t, IDLength)

		id0 := KademliaID(randomBytes)
		id1 := KademliaID(randomBytes)
		id2 := KademliaID(randomBytes1)

		if !id0.Equals(id1) {
			t.Fatalf("Id0 and Id2 should be the same! got %v, want %v", id1, id0)
		}

		if id0.Equals(id2) {
			t.Fatalf("Id0 cant be equal to id2! got %v, want %v", id2, id0)
		}

	})

	t.Run("Distance test high level", func(t *testing.T) {
		id0 := NewKademliaID("1111111100000000000000000000000000000000")
		id1 := NewKademliaID("1111111100000000000000000000000000000011")

		got := id0.CalcDistance(id1).String()
		want := "0000000000000000000000000000000000000011"

		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}

	})

	t.Run("XOR distance", func(t *testing.T) {
		bytes0 := []byte{0101}
		bytes1 := []byte{0100}

		want := []byte{0001}

		got, err := bitwiseXOR(bytes0, bytes1)

		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})

	t.Run("Less", func(t *testing.T) {
		id0 := NewKademliaID("1111111100000000000000000000000000000000")
		id1 := NewKademliaID("1111111100000000000000000000000000000010")

		smaller := id0.Less(id1)

		if !smaller {
			t.Fatalf("Id0 should be less than Id1")
		}

		isId1LessThanId0 := id1.Less(id0)

		if isId1LessThanId0 {
			t.Fatalf("Id1 should be bigger than Id0")
		}

		if id0.Equals(id1) {
			t.Fatalf("Id0 and Id1 should be different! got %v, want %v", id1, id0)
		}

	})

	t.Run("Random KademliaID", func(t *testing.T) {
		id0 := NewRandomKademliaID()
		id1 := NewRandomKademliaID()

		if id0.Equals(id1) {
			t.Fatalf("Id0 and Id1 should be different! got %v, want %v", id1, id0)
		}
	})
}

func generateRandomBytes(t *testing.T, n int) []byte {
	t.Helper()
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.Intn(256)) // Generate random number between 0 and 255
	}
	return b
}
