package main

import (
	"d7024e/kademlia"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting Kademlia nodes...")

	// Create KademliaIDs for the nodes
	id := kademlia.NewRandomKademliaID()
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
		MessageCh: messageCh1,
	}
	network2 := &kademlia.Network{ // Use pointers
		MessageCh: messageCh2,
	}

	// Create Kademlia instances for the two nodes with network and routing table references
	firstKademliaNode := kademlia.NewKademlia(network1, routingTable)
	secondKademliaNode := kademlia.NewKademlia(network2, secondRoutingTable)

	// Start the Kademlia nodes to process messages
	firstKademliaNode.Start()
	secondKademliaNode.Start()

	// Give some time for the nodes to start listening and processing messages
	time.Sleep(1 * time.Second)

	// Step 1: Send a ping message from the second node to the first node
	fmt.Printf("Sending ping from second node to first node: %v\n", contact)
	secondKademliaNode.Network.SendPingMessage(&contact, "ping:"+contact.ID.String()+":ping")

	// Step 2: Wait and give time for the ping-pong interaction to complete
	time.Sleep(1 * time.Second)

	// Step 3: Send a lookUpContact message from the second node to the first node
	fmt.Println("Sending lookUpContact from second node to first node...")
	lookupMessage := fmt.Sprintf("lookUpContact:%s:%s", secondContact.ID.String(), secondContact.ID.String())
	secondKademliaNode.Network.SendMessage(&contact, lookupMessage)

	// Step 4: Wait for lookUpContact to be processed and returnLookUpContact to be sent
	time.Sleep(4 * time.Second)

	fmt.Println("Kademlia nodes are running. Check logs for network activity.")
}
