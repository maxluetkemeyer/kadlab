package node

import (
	"context"
	"d7024e_group04/internal/kademlia/kademliaid"
	"testing"
	"time"
)

func TestNode_PutObject(t *testing.T) {
	testNodes := populateTestNodes()
	data := "some_data"

	node := Node{
		RoutingTable: testNodes[":18"].routingTable,
		Client:       newClientMock(testNodes, testNodes[":18"].contact),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := node.PutObject(ctx, data)
	if err != nil {
		t.Fatalf("failed to store data, err: %v", err)
	}
}

func TestNode_GetObject(t *testing.T) {
	testNodes := populateTestNodes()

	node := Node{
		RoutingTable: testNodes[":18"].routingTable,
		Store:        testNodes[":18"].store,
		Client:       newClientMock(testNodes, testNodes[":18"].contact),
	}

	node2 := Node{
		RoutingTable: testNodes[":1"].routingTable,
		Store:        testNodes[":1"].store,
		Client:       newClientMock(testNodes, testNodes[":1"].contact),
	}

	t.Run("find data that does not exist", func(t *testing.T) {
		hash := kademliaid.NewKademliaIDFromData("non-existent").String()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		value, candidates, err := node.GetObject(ctx, hash)

		if err != nil {
			t.Fatalf("error while running GetObject, err: %v", err)
		}

		if value != nil {
			t.Fatalf("found non-existent data: %v", value)
		}

		if candidates == nil {
			t.Fatalf("did not return any candidates")
		}

	})

	data := "some_data"
	key := kademliaid.NewKademliaIDFromData(data).String()

	node.Store.SetValue(key, data)
	node2.Store.SetValue(key, data)

	t.Run("find data on sending node", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		value, candidates, err := node.GetObject(ctx, key)
		if err != nil {
			t.Fatalf("error while running GetObject, err: %v", err)
		}

		if candidates != nil {
			t.Fatalf("returned candidates")
		}

		if value == nil {
			t.Fatalf("could not find data")
		}

		if value.DataValue != data {
			t.Fatalf("found invalid data, got %v, expected %v", value.DataValue, data)
		}

		if !value.NodeWithValue.ID.Equals(node.RoutingTable.Me().ID) {
			t.Fatalf("found data in wrong node, expected %v, got %v", node.RoutingTable.Me(), value.NodeWithValue)
		}
	})

	data2 := "some_other_data"
	key2AsKademliaID := kademliaid.NewKademliaIDFromData(data2)

	node2.Store.SetValue(key2AsKademliaID.String(), data2)
	t.Run("find data on another node", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		value, candidates, err := node.GetObject(ctx, key2AsKademliaID.String())
		if err != nil {
			t.Fatalf("error while running GetObject, err: %v", err)
		}

		if candidates != nil {
			t.Fatalf("returned candidates")
		}

		if value == nil {
			t.Fatalf("could not find data")
		}

		if value.DataValue != data2 {
			t.Fatalf("found invalid data, got %v, expected %v", value.DataValue, data2)
		}

		if !value.NodeWithValue.ID.Equals(node2.RoutingTable.Me().ID) {
			t.Fatalf("found data in wrong node, expected %v, got %v", node2.RoutingTable.Me(), value.NodeWithValue)
		}
	})

	t.Run("replicate data on closest node", func(t *testing.T) {
		data3 := "new_data"
		key3AsKademliaID := kademliaid.NewKademliaIDFromData(data3)
		node2.Store.SetValue(key3AsKademliaID.String(), data3)

		// check closest node
		candidates := node2.RoutingTable.FindClosestContacts(key3AsKademliaID, 1, node2.Me().ID)

		closestNode := testNodes[candidates[0].Address]
		valueFromClosestNode, err := closestNode.store.GetValue(key3AsKademliaID.String())

		if err == nil {
			t.Fatalf("GetValue found data that should not exist, data: %v", valueFromClosestNode)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		value, _, err := node.GetObject(ctx, key3AsKademliaID.String())

		if err != nil {
			t.Fatalf("failed to GetObject, err: %v", err)
		}

		if value.DataValue != data3 {
			t.Fatalf("wrong data in node, exepected %v, got %v", data3, value.DataValue)
		}

		// data should now be in closest node
		valueFromClosestNode, err = closestNode.store.GetValue(key3AsKademliaID.String())
		if err != nil {
			t.Fatalf("err fetching data: %v", err)
		}

		if valueFromClosestNode != data3 {
			t.Fatalf("failed to replicate to closest node, got %v, expected %v", valueFromClosestNode, data3)
		}

	})
}
