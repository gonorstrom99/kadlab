package main

import (
	"d7024e/kademlia"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting Kademlia nodes...")

	// Create KademliaIDs for the nodes
	id := kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	secondID := kademlia.NewRandomKademliaID()

	// Create Contacts for the two nodes
	contact := kademlia.NewContact(id, "127.0.0.1:8000")
	secondContact := kademlia.NewContact(secondID, "127.0.0.1:8001")

	// Create RoutingTables for the two nodes
	routingTable := kademlia.NewRoutingTable(contact)
	secondRoutingTable := kademlia.NewRoutingTable(secondContact)

	// Create Networks for the two nodes
	network := kademlia.Network{
		IP:   "127.0.0.1",
		Port: 8000, // Port for the first node
	}
	secondNetwork := kademlia.Network{
		IP:   "127.0.0.1",
		Port: 8001, // Port for the second node
	}

	// Create Kademlia instances for the two nodes
	firstKademliaNode := kademlia.NewKademlia(network, *routingTable)
	secondKademliaNode := kademlia.NewKademlia(secondNetwork, *secondRoutingTable)

	// Start listening on the networks in separate goroutines
	go func() {
		firstKademliaNode.Listen() // Start listening on the first node
	}()

	go func() {
		secondKademliaNode.Listen() // Start listening on the second node
	}()

	// Give some time for the nodes to start listening
	time.Sleep(1 * time.Second)

	// Send a ping message from the second node to the first node
	fmt.Printf("Sending ping from second node to first node: %v\n", contact)
	secondKademliaNode.Network.SendPingMessage(&contact)

	// Give some time for the ping-pong interaction to complete
	time.Sleep(1 * time.Second)

	fmt.Println("Kademlia nodes are running. Check logs for network activity.")

	time.Sleep(3 * time.Second)

}
