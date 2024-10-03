package node

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"d7024e_group04/env"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/network"
	"d7024e_group04/internal/store"
)

type Node struct {
	Client       network.ClientRPC
	RoutingTable *routingtable.RoutingTable
	Store        store.Store
	kNet         network.Network
}

func New(client network.ClientRPC, routingTable *routingtable.RoutingTable, store store.Store, kNet network.Network) *Node {
	return &Node{
		Client:       client,
		RoutingTable: routingTable,
		Store:        store,
		kNet:         kNet,
	}
}

/*
- A node joins the network as follows:
 1. if it does not already have a nodeID n, it generates one
 2. it inserts the value of some known node c into the appropriate bucket as its first contact
 3. it does an iterativeFindNode for n
 4. it refreshes all buckets further away than its closest neighbor, which will be in the occupied bucket with the lowest index.
*/
func (n *Node) Bootstrap(rootCtx context.Context) error {
	var recipient *contact.Contact

	logger := slog.Default().With("me", n.RoutingTable.Me().Address)
	ctx, cancel := context.WithTimeout(rootCtx, env.BootstrapTimeout)
	for {
		if ctx.Err() != nil {
			defer cancel()
			return ctx.Err()
		}

		ips, err := n.kNet.ResolveDNS(env.NodesProxyDomain)
		if err != nil {
			logger.Warn("Unable to resolve DNS", slog.Any("domain", env.NodesProxyDomain))
			time.Sleep(100 * time.Millisecond)
			continue
		}

		ourAddress := n.RoutingTable.Me().Address
		ips = removeAddress(ips, ourAddress)

		logger := logger.With(slog.Any("ips", ips))

		logger.Debug("pinging")

		recipient, err = n.pingIPsAndGetContact(ctx, ips)
		if err == nil {
			break
		}

		logger.Warn("Unable to ping any ip", slog.Any("err", err))
		slog.Warn("Retrying bootstrap in 100 milis")
		time.Sleep(100 * time.Millisecond)
	}

	cancel()

	n.RoutingTable.AddContact(recipient)
	// TODO: iterative search, should we update once for the list or for each node visited in findNode?
	me := n.RoutingTable.Me()
	closestContacts := n.findNode(rootCtx, me)
	for _, contact := range closestContacts {
		n.RoutingTable.AddContact(contact)
	}

	logger.Debug("FOUND NODES", slog.Any("nodes", closestContacts))

	return nil
}

// Put takes content of the file and outputs the hash of the object
func (n *Node) PutObject() {
	panic("TODO")
}

// Get takes hash and outputs the contents of the object and the node it was retrieved
func (n *Node) GetObject(rootCtx context.Context, hash string) (data string, err error) {
	panic("TODO")
}

// Ping each contact in <contacts> until one responeses and returns it.
func (n *Node) pingIPsAndGetContact(ctx context.Context, targetIPs []string) (*contact.Contact, error) {
	for _, targetIP := range targetIPs {
		contact, err := n.Client.SendPing(ctx, fmt.Sprintf("%v:%v", targetIP, env.Port))
		if err == nil {
			return contact, nil
		}
	}

	return nil, fmt.Errorf("unable to ping any contacts")
}

func removeAddress(ips []string, addressToRemove string) []string {
	portLength := 5 + /* colon */ 1
	ipToRemove := addressToRemove[:len(addressToRemove)-portLength]

	for idx, ip := range ips {
		if ip == ipToRemove {
			// Remove the ip
			ips = append(ips[:idx], ips[idx+1:]...)
			break
		}
	}

	return ips
}
