package main

import (
	"d7024e/kademlia"
	"fmt"
	"time"
)

//trace print format: file: [insert file name] function: [insert function]

func main() {
	fmt.Println("Starting Kademlia nodes...")

	// Create Kademlia instances for the two nodes, with IP and Port
	firstKademliaNode := kademlia.CreateKademliaNode("127.0.0.1:8000")
	secondKademliaNode := kademlia.CreateKademliaNode("127.0.0.1:8001")

	// Start the Kademlia nodes to process messages
	firstKademliaNode.Start()
	secondKademliaNode.Start()

	// Give some time for the nodes to start listening and processing messages
	time.Sleep(1 * time.Second)

	// Step 1: Send a ping message from the second node to the first node
	fmt.Printf("Sending ping from second node to first node: %v\n", firstKademliaNode.RoutingTable.me)
	secondKademliaNode.Network.SendPingMessage(firstKademliaNode.RoutingTable.me, "ping:"+firstKademliaNode.RoutingTable.me.ID.String()+":ping") //det här känns inte som att det borde testas här tbh

	// Step 2: Wait and give time for the ping-pong interaction to complete
	time.Sleep(1 * time.Second)

	// Step 3: Send a lookUpContact message from the second node to the first node
	fmt.Println("Sending lookUpContact from second node to first node...")
	lookupMessage := fmt.Sprintf("lookUpContact:%s:%s", secondKademliaNode.RoutingTable.me.ID.String(), secondKademliaNode.RoutingTable.me.ID.String())
	secondKademliaNode.Network.SendMessage(firstKademliaNode.RoutingTable.me, lookupMessage)

	// Step 4: Wait for lookUpContact to be processed and returnLookUpContact to be sent
	time.Sleep(4 * time.Second)

	fmt.Println("Kademlia nodes are running. Check logs for network activity.")
}
