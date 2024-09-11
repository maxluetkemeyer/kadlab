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
}

// NewBucket returns a new instance of a Bucket
func NewBucket() *Bucket {
	bucket := &Bucket{}
	bucket.list = list.New()
	return bucket
}

// Len return the size of the bucket
func (bucket *Bucket) Len() int {
	return bucket.list.Len()
}
