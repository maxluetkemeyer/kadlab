package bucket

// least-recently seen (head) (old)
// most-recently seen (tail) (new)

import (
	"container/list"
	"fmt"
	"time"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
)

// Bucket definition
// contains a List
type Bucket struct {
	list        *list.List
	size        int
	lastRefresh time.Time
}

// NewBucket returns a new instance of a Bucket
func NewBucket(size int) *Bucket {
	// ignore error since we set this env in runtime
	bucket := &Bucket{size: size}
	bucket.list = list.New()
	bucket.Refresh()
	return bucket
}

// Len return the size of the bucket
func (bucket *Bucket) Len() int {
	return bucket.list.Len()
}

func (bucket *Bucket) NeedsRefresh() bool {
	now := time.Now()
	nextRefresh := bucket.lastRefresh.Add(env.TRefresh)
	return nextRefresh.Compare(now) <= 0
}

func (bucket *Bucket) CheckRefresh(refreshChannel chan<- Bucket) {
	if bucket.NeedsRefresh() {
		// If nextRefresh has already happen, or is happening right now
		go func() {
			refreshChannel <- *bucket
		}()
		bucket.Refresh()
	}
}

func (bucket *Bucket) Refresh() {
	bucket.lastRefresh = time.Now()
}

func (bucket *Bucket) GetContact(idx int) (*contact.Contact, error) {
	i := 0
	for contactElement := bucket.list.Front(); contactElement != nil; contactElement = contactElement.Next() {
		if idx == i {
			contactFromList := *contactElement.Value.(*contact.Contact) // use generics with list
			return &contactFromList, nil
		}
		i++
	}

	return nil, fmt.Errorf("couldn't find contact with idx=%v in bucket", idx)
}
