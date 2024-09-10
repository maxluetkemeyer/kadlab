package network

import "d7024e_group04/kademlia/contact"

type Network struct {
}

// TODO: Lets start with a simulated network
// TODO: Define a network interface to have a simulated on and a real one (and maybe a spy test one)

func Listen(ip string, port int) {
	// TODO
}

func (network *Network) SendPingMessage(contact *contact.Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *contact.Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
