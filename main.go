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

	// Create message channels for each network
	messageCh1 := make(chan kademlia.Message)
	messageCh2 := make(chan kademlia.Message)

	// Create Networks for the two nodes with message channels
	network1 := &kademlia.Network{ // Use pointers
		IP:        "127.0.0.1",
		Port:      8000,
		MessageCh: messageCh1,
	}
	network2 := &kademlia.Network{ // Use pointers
		IP:        "127.0.0.1",
		Port:      8001,
		MessageCh: messageCh2,
	}

	// Create Kademlia instances for the two nodes with network and routing table references
	firstKademliaNode := kademlia.NewKademlia(network1, routingTable)
	secondKademliaNode := kademlia.NewKademlia(network2, secondRoutingTable)

	// Start the Kademlia nodes to process messages
	go firstKademliaNode.Start()
	go secondKademliaNode.Start()

	// Give some time for the nodes to start listening and processing messages
	time.Sleep(1 * time.Second)

	// Send a ping message from the second node to the first node
	fmt.Printf("Sending ping from second node to first node: %v\n", contact)
	secondKademliaNode.Network.SendPingMessage(&contact)

	// Give some time for the ping-pong interaction to complete
	time.Sleep(4 * time.Second)

	fmt.Println("Kademlia nodes are running. Check logs for network activity.")
}
