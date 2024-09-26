package main

import (
	"d7024e/kademlia"
	"fmt"
	"time"
)

func main() {
	fmt.Println("(file: main) Starting Kademlia nodes...")
	KademliaNode1 := kademlia.CreateKademliaNode("127.0.0.1:8000")
	KademliaNode2 := kademlia.CreateKademliaNode("127.0.0.1:8001")
	KademliaNode3 := kademlia.CreateKademliaNode("127.0.0.1:8002")
	KademliaNode4 := kademlia.CreateKademliaNode("127.0.0.1:8003")
	KademliaNode5 := kademlia.CreateKademliaNode("127.0.0.1:8004")
	KademliaNode6 := kademlia.CreateKademliaNode("127.0.0.1:8005")

	// Start the Kademlia nodes to process messages
	KademliaNode1.Start()
	KademliaNode2.Start()
	KademliaNode3.Start()
	KademliaNode4.Start()
	KademliaNode5.Start()
	KademliaNode6.Start()

	// Give some time for the nodes to start listening and processing messages
	time.Sleep(1 * time.Second)
	KademliaNode1.RoutingTable.AddContact(*KademliaNode2.RoutingTable.GetMe())
	KademliaNode2.RoutingTable.AddContact(*KademliaNode3.RoutingTable.GetMe())
	KademliaNode2.RoutingTable.AddContact(*KademliaNode4.RoutingTable.GetMe())
	KademliaNode2.RoutingTable.AddContact(*KademliaNode5.RoutingTable.GetMe())
	// Step 1: Send a ping message from the second node to the first node
	//fmt.Printf("Sending ping from second node to first node: %v\n", KademliaNode2.)
	//KademliaNode2.Network.SendPingMessage(KademliaNode1, "ping:"+secondID.String()+":ping")

	// Step 2: Wait and give time for the ping-pong interaction to complete
	time.Sleep(1 * time.Second)

	//node 1 is sending a lookupmesseage to node 2 to look up node 1, it then adds node 3, 4, and 5 to its routingtable
	fmt.Println("(file: main) Sending lookUpContact from first node to second node...")
	lookupMessage := fmt.Sprintf("lookUpContact:%s:%s", KademliaNode1.Network.ID.String(), KademliaNode1.Network.ID.String())
	KademliaNode1.Network.SendMessage(KademliaNode2.RoutingTable.GetMe(), lookupMessage)

	// Step 4: Wait for lookUpContact to be processed and returnLookUpContact to be sent
	time.Sleep(2 * time.Second)
	fmt.Println("(file: main) Value of iscontactinroutingtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode2.RoutingTable.GetMe()))

	time.Sleep(2 * time.Second)

}
