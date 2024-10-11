package routingtable

import (
	"d7024e_group04/env"
	"d7024e_group04/internal/client"
	"d7024e_group04/internal/kademlia/bucket"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/network"
	"sync"
)

// RoutingTable definition
// keeps a reference contact of me and an array of buckets
// 160 buckets with the current IDLength
type RoutingTable struct {
	mut     sync.RWMutex
	me      *contact.Contact
	client  network.ClientRPC
	buckets [env.IDLength * 8]*bucket.Bucket
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me *contact.Contact) *RoutingTable {
	// ignore err for now, we set this in runtime
	routingTable := &RoutingTable{
		me:     me,
		client: client.NewClient(me),
	}
	for i := 0; i < env.IDLength*8; i++ {
		routingTable.buckets[i] = bucket.NewBucket(env.BucketSize)
	}
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact *contact.Contact) {
	routingTable.mut.Lock()
	defer routingTable.mut.Unlock()

	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddContact(contact, routingTable.client)
}

// FindClosestContacts finds the 'count' closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target kademliaid.KademliaID, count int, blacklist ...kademliaid.KademliaID) []*contact.Contact {
	candidates := make([]*contact.Contact, 0)
	// Find in which bucket the target should be
	routingTable.mut.RLock()
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	// Get all contacts in the bucket with already calculated distances
	candidates = append(candidates, bucket.GetContactAndCalcDistance(target)...)

	// TODO: Put condition in extra function
	// If we do not have enough candidates, we check the two nearest buckets and so on and so on
	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < env.IDLength*8) && len(candidates) < count; i++ {
		// TODO: Can we follow DRY principle? dont repeat yourself
		// Add candidates of the smaller nearest bucket
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates = append(candidates, bucket.GetContactAndCalcDistance(target)...)
		}
		// Add candidates of the bigger nearest bucket
		if bucketIndex+i < env.IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates = append(candidates, bucket.GetContactAndCalcDistance(target)...)
		}
	}

	routingTable.mut.RUnlock()

	contact.SortContacts(candidates)

	// Maybe we have too little
	if count > len(candidates) {
		count = len(candidates)
	}

	// make sure blacklisted id is not in list
	for _, node := range blacklist {
		candidates = contact.RemoveID(candidates, node)
	}

	// If we have too much in our candidates, the get contacts function returns the right amount
	return candidates[:count]
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) getBucketIndex(id kademliaid.KademliaID) int {
	// distance to ourselves
	distance := id.CalcDistance(routingTable.me.ID)

	// TODO: Remove abbreviations
	// TODO: Simplify loop
	for i := 0; i < env.IDLength; i++ {
		for j := 0; j < 8; j++ {
			// Loop through each bit of the id
			// TODO: Maybe stick to byte type?

			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return env.IDLength*8 - 1
}

func (routingTable *RoutingTable) Me() *contact.Contact {
	routingTable.mut.RLock()
	defer routingTable.mut.RUnlock()
	return routingTable.me
}
