package main

import (
	"d7024e/kademlia"
	"fmt"
	"time"
)

func main() {
	fmt.Println("(file: main) Starting Kademlia nodes...")

	// Create KademliaIDs for the nodes
	// id := kademlia.NewRandomKademliaID()
	// secondID := kademlia.NewRandomKademliaID()

	// // Print the Kademlia IDs for the nodes
	// //fmt.Printf("Kademlia ID of first node: %s\n", id.String())
	// //fmt.Printf("Kademlia ID of second node: %s\n", secondID.String())

	// // Create Contacts for the two nodes
	// contact := kademlia.NewContact(id, "127.0.0.1:8000")
	// secondContact := kademlia.NewContact(secondID, "127.0.0.1:8001")

	// // Create RoutingTables for the two nodes
	// routingTable := kademlia.NewRoutingTable(contact)
	// secondRoutingTable := kademlia.NewRoutingTable(secondContact)

	// // Create message channels for each network
	// messageCh1 := make(chan kademlia.Message)
	// messageCh2 := make(chan kademlia.Message)

	// // Create Networks for the two nodes with message channels
	// network1 := &kademlia.Network{
	// 	MessageCh: messageCh1,
	// }
	// network2 := &kademlia.Network{
	// 	MessageCh: messageCh2,
	// }

	// // Create Kademlia instances for the two nodes with network and routing table references
	// firstKademliaNode := kademlia.NewKademlia(network1, routingTable)
	// secondKademliaNode := kademlia.NewKademlia(network2, secondRoutingTable)
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
	fmt.Println("(file: main) Sending lookUpContact from second node to first node...")
	lookupMessage := fmt.Sprintf("lookUpContact:%s:%s", KademliaNode1.Network.ID.String(), KademliaNode1.Network.ID.String())
	KademliaNode1.Network.SendMessage(KademliaNode2.RoutingTable.GetMe(), lookupMessage)

	// Step 4: Wait for lookUpContact to be processed and returnLookUpContact to be sent
	time.Sleep(2 * time.Second)
	fmt.Println("(file: main) Value of iscontactinroutingtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode2.RoutingTable.GetMe()))

	fmt.Println("(file: main) Value of iscontactinroutingtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode3.RoutingTable.GetMe()))
	fmt.Println("(file: main) Value of iscontactinroutingtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode4.RoutingTable.GetMe()))
	fmt.Println("(file: main) Value of iscontactinroutingtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode5.RoutingTable.GetMe()))

	time.Sleep(2 * time.Second)

}
