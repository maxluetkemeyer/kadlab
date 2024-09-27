package node

import (
	"d7024e_group04/internal/kademlia/contact"
	"slices"
	"sync"
)

type kClosestList struct {
	mut     sync.RWMutex
	list    []*contact.Contact
	updated bool
}

func (kClosestList *kClosestList) Has(targetContact *contact.Contact) bool {
	kClosestList.mut.RLock()
	defer kClosestList.mut.RUnlock()

	for _, c := range kClosestList.list {
		if c.ID.Equals(targetContact.ID) {
			return true
		}
	}
	return false
}

func (kClosestList *kClosestList) isSubset(set *contact.ContactSet) bool {
	kClosestList.mut.RLock()
	defer kClosestList.mut.RUnlock()

	for _, contact := range kClosestList.list {
		if !set.Has(contact) {
			return false
		}
	}
	return true
}

func (kClosestList *kClosestList) sort() {
	kClosestList.mut.Lock()
	defer kClosestList.mut.Unlock()
	slices.SortStableFunc(kClosestList.list, func(a, b *contact.Contact) int {
		if a.Less(b) {
			return -1
		} else {
			return 1
		}
	})
}

func (kClosestList *kClosestList) remove(target *contact.Contact) {
	kClosestList.mut.Lock()
	defer kClosestList.mut.Unlock()

	var contactList []*contact.Contact

	for _, contact := range kClosestList.list {
		if !contact.ID.Equals(target.ID) {
			contactList = append(contactList, contact)
		}
	}

	kClosestList.list = contactList
}

func (kClosestList *kClosestList) addContacts(contacts []*contact.Contact, referenceContact *contact.Contact) {
	for _, contact := range contacts {
		contact.CalcDistance(referenceContact.ID)

		if kClosestList.Has(contact) {
			continue
		}

		if len(kClosestList.list) < k {
			kClosestList.list = append(kClosestList.list, contact)
			kClosestList.sort()
			kClosestList.updated = true
			continue
		}

		if contact.Less(kClosestList.list[k-1]) {
			kClosestList.list[k-1] = contact
			kClosestList.sort()
			kClosestList.updated = true
		}
	}
}
