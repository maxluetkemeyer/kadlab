package node

import (
	"d7024e_group04/internal/kademlia/contact"
	"slices"
	"sync"
)

type NodeList struct {
	mut      sync.Mutex
	slice    []contact.Contact
	size     int
	modified bool
}

func NewNodeList(size int) NodeList {
	return NodeList{
		size:     size,
		modified: false,
	}
}

func (n *NodeList) AddNode(node contact.Contact, target *contact.Contact) {
	n.mut.Lock()
	defer n.mut.Unlock()

	if n.Contains(node) {
		return
	}

	node.CalcDistance(target.ID)

	if len(n.slice) < n.size {
		n.slice = append(n.slice, node)
		n.sort()
	} else {
		if node.Less(&n.slice[n.size-1]) {
			n.slice[n.size-1] = node
			n.sort()
		}
	}
}

func (n *NodeList) AddNodes(nodes []contact.Contact, target *contact.Contact) {
	for _, node := range nodes {
		n.AddNode(node, target)
	}
}

func (n *NodeList) Contains(node contact.Contact) bool {
	for _, candidate := range n.slice {
		if candidate.ID.Equals(node.ID) {
			return true
		}
	}
	return false
}

func (n *NodeList) GetClosest() []contact.Contact {
	return n.slice[:]
}

func (n *NodeList) ResetModifiedFlag() {
	n.modified = false
}


func (n *NodeList) HasBeenModified() bool {
	return n.modified
}

func (n *NodeList) sort() {
	slices.SortStableFunc(n.slice, func(a, b contact.Contact) int {
		if a.Less(&b) {
			return -1
		} else {
			return 1
		}
	})
}
