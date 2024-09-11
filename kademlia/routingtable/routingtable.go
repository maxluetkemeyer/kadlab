package routingtable

import (
	"d7024e_group04/kademlia/bucket"
	"d7024e_group04/kademlia/contact"
	"d7024e_group04/kademlia/kademliaid"
	"os"
	"strconv"
)

// RoutingTable definition
// keeps a reference contact of me and an array of buckets
// TODO: Do we want to store our address here?
// 160 buckets with the current IDLength
type RoutingTable struct {
	me      contact.Contact
	buckets [kademliaid.IDLength * 8]*bucket.Bucket
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me contact.Contact) *RoutingTable {
	// ignore err for now, we set this in runtime
	size, _ := strconv.Atoi(os.Getenv("BUCKET_SIZE"))

	routingTable := &RoutingTable{}
	for i := 0; i < kademliaid.IDLength*8; i++ {
		routingTable.buckets[i] = bucket.NewBucket(size)
	}
	routingTable.me = me
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact contact.Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddContact(contact)
}

// FindClosestContacts finds the 'count' closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *kademliaid.KademliaID, count int) []contact.Contact {
	// TODO: It is a slice of contacts
	var candidates contact.ContactCandidates
	// Find in which bucket the target should be
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	// Get all contacts in the bucket with already calculated distances
	candidates.Append(bucket.GetContactAndCalcDistance(target))

	// TODO: Put condition in extra function
	// If we do not have enough candidates, we check the two nearest buckets and so on and so on
	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < kademliaid.IDLength*8) && candidates.Len() < count; i++ {
		// TODO: Can we follow DRY principle? dont repeat yourself
		// Add candidates of the smaller nearest bucket
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		// Add candidates of the bigger nearest bucket
		if bucketIndex+i < kademliaid.IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	// Maybe we have too little
	if count > candidates.Len() {
		count = candidates.Len()
	}

	// If we have to much in our candidates, the get contacts function returns the right amount
	return candidates.GetContacts(count)
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) getBucketIndex(id *kademliaid.KademliaID) int {
	// distance to ourselves
	distance := id.CalcDistance(routingTable.me.ID)

	// TODO: Remove abbreviations
	// TODO: Simplify loop
	for i := 0; i < kademliaid.IDLength; i++ {
		for j := 0; j < 8; j++ {
			// Loop thorugh each bit of the id
			// TODO: Maybe stick to byte type?

			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return kademliaid.IDLength*8 - 1
}
