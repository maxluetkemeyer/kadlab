package kademliaid

import (
	"encoding/hex"
	"math/rand"

	"d7024e_group04/env"
)

// The static number of bytes in a KademliaID. 160 / 8 = 20
type KademliaID [env.IDLength]byte

// NewKademliaID returns a new KademliaID based on the string input
func NewKademliaID(data string) KademliaID {
	// TODO: data is stored as hex at the moment
	// byte and error
	decoded, _ := hex.DecodeString(data)

	// new variable, only declared (initialized with the "zero" value)
	newKademliaID := KademliaID{}
	// TODO: this for loop literally makes no sense, you can directly assign it
	for i := 0; i < env.IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}

	// the address of the new internal id
	return newKademliaID
}

// NewRandomKademliaID returns a new random KademliaID,
func NewRandomKademliaID() KademliaID {
	newKademliaID := KademliaID{}
	for i := range env.IDLength {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return newKademliaID
}

func NewKademliaIDFromBytes(data [env.IDLength]byte) KademliaID {
	return data
}

// Less returns true if kademliaID < otherKademliaID (bitwise)
func (kademliaID KademliaID) Less(otherKademliaID KademliaID) bool {
	for i := 0; i < env.IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals returns true if kademliaID == otherKademliaID (bitwise)
func (kademliaID KademliaID) Equals(otherKademliaID KademliaID) bool {
	for i := 0; i < env.IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

// String returns a simple hex string representation of a KademliaID
func (kademliaID KademliaID) String() string {
	return hex.EncodeToString(kademliaID.Bytes())
}

// Bytes returns the kademliaID as a byte array
func (kademliaID KademliaID) Bytes() []byte {
	bytes := make([]byte, 0, env.IDLength)
	for _, b := range kademliaID {
		bytes = append(bytes, b)
	}
	return bytes
}
