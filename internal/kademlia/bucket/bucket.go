package bucket

// least-recently seen (head) (old)
// most-recently seen (tail) (new)

import (
	"container/list"
)

//TODO: Read paper
//TODO: Generic List

// Bucket definition
// contains a List
// TODO: Need to be public for "routingtable.go", make this private
type Bucket struct {
	list *list.List
	size int
}

// NewBucket returns a new instance of a Bucket
func NewBucket(size int) *Bucket {
	// ignore error since we set this env in runtime
	bucket := &Bucket{size: size}
	bucket.list = list.New()
	return bucket
}

// Len return the size of the bucket
func (bucket *Bucket) Len() int {
	return bucket.list.Len()
}
