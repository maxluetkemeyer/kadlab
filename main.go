// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"d7024e_group04/kademlia/bucket"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	"fmt"
)

func main() {
	fmt.Println("Pretending to run the kademlia app...")
	// Using stuff from the kademlia package here. Something like...
	id := kademliaid.NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	myContact := contact.NewContact(id, "localhost:8000")
	myBucket := bucket.NewBucket()
	myBucket.Len()

	fmt.Println(myContact.String())
	fmt.Printf("%v\n", myContact)
}
