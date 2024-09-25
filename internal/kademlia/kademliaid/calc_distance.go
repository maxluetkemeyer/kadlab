package kademliaid

import (
	"fmt"
)

// CalcDistance returns a new instance of a KademliaID that is built
// through a bitwise XOR operation between kademliaID and target
func (kademliaID KademliaID) CalcDistance(target KademliaID) KademliaID {
	byteSlice, _ := bitwiseXOR(kademliaID.Bytes(), target.Bytes())

	return KademliaID(byteSlice)
}

func bitwiseXOR(byteSlice0, byteSlice1 []byte) ([]byte, error) {
	if len(byteSlice0) != len(byteSlice1) {
		return []byte{}, fmt.Errorf("the slices have different size: %d vs %d", len(byteSlice0), len(byteSlice1))
	}

	sliceLen := len(byteSlice0)
	result := make([]byte, sliceLen)

	for i := range sliceLen {
		// bitwise XOR
		result[i] = byteSlice0[i] ^ byteSlice1[i]
	}

	return result, nil
}
