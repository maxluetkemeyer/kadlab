package kademliaid

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"

	"d7024e_group04/env"
)

// The static number of bytes in a KademliaID. 160 / 8 = 20
type KademliaID [env.IDLength]byte

// NewKademliaID returns a new KademliaID based on the string input
// TODO: Does not perform length checking yet!
func NewKademliaID(dataAsHexString string) KademliaID {
	// []byte and error
	decoded, _ := hex.DecodeString(dataAsHexString)

	// If the hex is smaller than env.IDLength, just pad it with zeros
	for i := len(decoded); i < env.IDLength; i++ {
		decoded[i] = 0
	}

	// new variable, only declared (initialized with the "zero" value)
	newKademliaID := KademliaID(decoded)

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

func NewKademliaIDFromData(data string) KademliaID {
	// sha1 for now since it gives us 160 bit hash, we could use something better and truncate to 160 but specification mentioned sha1
	hasher := sha1.New()
	hasher.Write([]byte(data))

	hash := hasher.Sum(nil)

	return [env.IDLength]byte(hash)
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
