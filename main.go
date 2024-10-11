// A Kademlia distributed data store
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"d7024e_group04/api"
	"d7024e_group04/cli"
	"d7024e_group04/env"
	"d7024e_group04/internal/client"
	"d7024e_group04/internal/kademlia/contact"
	"d7024e_group04/internal/kademlia/kademliaid"
	"d7024e_group04/internal/kademlia/routingtable"
	"d7024e_group04/internal/network"
	"d7024e_group04/internal/node"
	"d7024e_group04/internal/server"
	"d7024e_group04/internal/store"

	"golang.org/x/sync/errgroup"
)

func main() {
	rootCtx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	errGroup, errCtx := errgroup.WithContext(rootCtx)

	me := createOwnContact()
	routingTable := routingtable.NewRoutingTable(me)
	memoryStore := store.NewMemoryStore()
	simpleTtlStore := store.NewSimpleTTLStore(memoryStore)
	client := client.NewClient(me)
	kNetwork := &network.PublicNetwork{}
	node := node.New(client, routingTable, simpleTtlStore, kNetwork)

	startServer(errGroup, errCtx, routingTable, simpleTtlStore)
	startAPI(errGroup, errCtx, node)
	startCLI(errGroup, errCtx, cancelCtx, node)
	startBootstrapping(errGroup, errCtx, node)

	err := errGroup.Wait()
	cancelCtx()
	log.Printf("Node shutdown, reason: %v", err)
}

func createOwnContact() (me *contact.Contact) {
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}
	ip, err := net.LookupIP(host)
	if err != nil {
		log.Fatalf("Invalid IP: %v", ip)
	}

	ipWithPort := ip[0].String() + ":" + strconv.Itoa(env.Port)

	// Assumption: this will give unique IDs for the whole network
	id := kademliaid.NewRandomKademliaID()

	return contact.NewContact(id, ipWithPort)
}

func startServer(errGroup *errgroup.Group, errCtx context.Context, routingTable *routingtable.RoutingTable, store store.TTLStore) {
	newServer := server.NewServer(routingTable, store)
	errGroup.Go(func() error {
		log.Println("STARTING SERVER")
		return newServer.Start(errCtx)
	})
}

func startAPI(errGroup *errgroup.Group, errCtx context.Context, node node.NodeHandler) {
	var handler api.KademliaAPI = api.NewHandler(node)
	errGroup.Go(func() error {
		log.Println("STARTING API")
		return handler.ListenAndServe(errCtx)
	})
}

func startCLI(errGroup *errgroup.Group, errCtx context.Context, cancelCtx context.CancelFunc, node node.NodeHandler) {
	errGroup.Go(func() error {
		log.Println("STARTING CLI")
		return cli.InputLoop(errCtx, cancelCtx, node)
	})
}

func startBootstrapping(errGroup *errgroup.Group, errCtx context.Context, node node.NodeHandler) {
	errGroup.Go(func() error {
		err := node.Bootstrap(errCtx)
		if err != nil {
			log.Fatalf("BOOTSTRAP failed: %v\n", err)
		} else {
			log.Printf("BOOTSTRAP succeeded")
		}

		return err
	})
}
