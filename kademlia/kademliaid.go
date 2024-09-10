package kademlia

import (
	"encoding/hex"
	"math/rand"
)

// the static number of bytes in a KademliaID
// TODO: Why use bytes and not bit? Does go has a bit interface? No it does not
// 160 / 8 = 20
// TODO: should it not be private?
const IDLength = 20

// type definition of a KademliaID
// 20 byte array
type KademliaID [IDLength]byte

// NewKademliaID returns a new instance of a KademliaID based on the string input
func NewKademliaID(data string) *KademliaID {
	// TODO: data is stored as hex at the moment
	// byte and error
	decoded, _ := hex.DecodeString(data)

	// new variable, only declared (initialized with the "zero" value)
	// u cant use := outside a function
	newKademliaID := KademliaID{}
	// TODO: this for loop literaly makes no sense, you can directly assign it
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}

	// the address of the new kademlia id
	return &newKademliaID
}

// NewRandomKademliaID returns a new instance of a random KademliaID,
// TODO: change this to a better version if you like
func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

// Less returns true if kademliaID < otherKademliaID (bitwise)
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals returns true if kademliaID == otherKademliaID (bitwise)
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

// CalcDistance returns a new instance of a KademliaID that is built
// through a bitwise XOR operation betweeen kademliaID and target
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		// bitwise XOR
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

// String returns a simple string representation of a KademliaID
// TODO: do we want to work with hex?
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
