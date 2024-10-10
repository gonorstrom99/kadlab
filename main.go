package main

import (
	// "d7024e/kademlia"

	"bufio"
	"kadlab/kademlia"

	"log"
	"os"
	"strings"
	"time"

	"fmt"
)

func startprogram(ID string) {
	log.Printf("asd")
}
func create(ID string, Address string) {
	log.Printf("Creating with ID: %s, Address: %s\n", ID, Address)
}

func main() {
	fmt.Println("Docker node is running. Type 'put' followed by your command:")
	id := os.Getenv("ID")
	address := os.Getenv("ADDRESS")
	kademlia.TestPrinter("aaaaaaaaaaaaaa")
	//kademlia.TestPrinter(address)

	// KademliaNode1 := kademlia.CreateKademliaNode(address)
	// log.Printf(KademliaNode1.RoutingTable.GetMe().Address)
	// KademliaNode1.Start()
	// if address == "127.0.0.1:8001" {
	// 		KademliaNode1.Network.ID = *kademlia.NewKademliaID("Node1")
	// }

	// if address == "127.0.0.1:8002" {
	// 	contact := kademlia.NewContact(kademlia.NewKademliaID("Node1"), "127.0.0.1:8001")

	// 	KademliaNode1.RoutingTable.AddContact(contact)

	// 	KademliaNode1.StartTask(contact.ID, "LookupContact", "---")

	// }
	// Call the create function with the provided values.
	create(id, address)
	// Create a new scanner to read input from stdin.
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	log.Printf(address)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			// Read the input and trim any extra spaces.
			input := strings.TrimSpace(scanner.Text())
			// Check if the input starts with the "put" command.
			if strings.HasPrefix(input, "put") {
				// Execute the command or handle the "put" command.
				log.Printf(id)
				startprogram(input)
			} else if strings.HasPrefix(input, "exit") {
				break
			} else {
				fmt.Println("Unknown command. Type 'put' followed by a command.")
			}
		} else {
			// If there's an error or EOF, break the loop.
			break
		}
	}

	// KademliaNode1 := kademlia.CreateKademliaNode("127.0.0.1:8001")
	// //KademliaNode1 := kademlia.CreateKademliaNode("127.0.0.1:portfromdocker")
	// //if kademlia.port == 8001 {
	// //kademlia.id="newkademliaID(AAAAAAA)"
	// //}else {
	// //	lookup på sig själv till (AAAAAAA, 127.0.0.1:8001) node}
	// //
	// //

	// KademliaNode2 := kademlia.CreateKademliaNode("127.0.0.1:8002")
	// KademliaNode3 := kademlia.CreateKademliaNode("127.0.0.1:8003")
	// KademliaNode4 := kademlia.CreateKademliaNode("127.0.0.1:8004")
	// KademliaNode5 := kademlia.CreateKademliaNode("127.0.0.1:8005")
	// KademliaNode6 := kademlia.CreateKademliaNode("127.0.0.1:8006")
	// KademliaNode7 := kademlia.CreateKademliaNode("127.0.0.1:8007")
	// KademliaNode8 := kademlia.CreateKademliaNode("127.0.0.1:8008")
	// KademliaNode9 := kademlia.CreateKademliaNode("127.0.0.1:8009")
	// KademliaNode10 := kademlia.CreateKademliaNode("127.0.0.1:8010")
	// KademliaNode11 := kademlia.CreateKademliaNode("127.0.0.1:8011")

	// // Start the Kademlia nodes to process messages
	// KademliaNode1.Start()
	// KademliaNode2.Start()
	// KademliaNode3.Start()
	// KademliaNode4.Start()
	// KademliaNode5.Start()
	// KademliaNode6.Start()
	// KademliaNode7.Start()
	// KademliaNode8.Start()
	// KademliaNode9.Start()
	// KademliaNode10.Start()
	// KademliaNode11.Start()
	// // Give some time for the nodes to start listening and processing messages
	// time.Sleep(1 * time.Second)
	// KademliaNode1.RoutingTable.AddContact(*KademliaNode2.RoutingTable.GetMe())
	// KademliaNode2.RoutingTable.AddContact(*KademliaNode3.RoutingTable.GetMe())
	// KademliaNode2.RoutingTable.AddContact(*KademliaNode4.RoutingTable.GetMe())
	// KademliaNode2.RoutingTable.AddContact(*KademliaNode5.RoutingTable.GetMe())
	// KademliaNode5.RoutingTable.AddContact(*KademliaNode6.RoutingTable.GetMe())

	// KademliaNode6.RoutingTable.AddContact(*KademliaNode7.RoutingTable.GetMe())
	// KademliaNode6.RoutingTable.AddContact(*KademliaNode5.RoutingTable.GetMe())
	// KademliaNode7.RoutingTable.AddContact(*KademliaNode4.RoutingTable.GetMe())
	// KademliaNode7.RoutingTable.AddContact(*KademliaNode2.RoutingTable.GetMe())
	// KademliaNode7.RoutingTable.AddContact(*KademliaNode8.RoutingTable.GetMe())
	// KademliaNode8.RoutingTable.AddContact(*KademliaNode11.RoutingTable.GetMe())
	// KademliaNode8.RoutingTable.AddContact(*KademliaNode9.RoutingTable.GetMe())
	// KademliaNode9.RoutingTable.AddContact(*KademliaNode10.RoutingTable.GetMe())
	// KademliaNode10.RoutingTable.AddContact(*KademliaNode11.RoutingTable.GetMe())

	// //KademliaNode1.UpdateRoutingTable(KademliaNode6.RoutingTable.GetMe())

	// // Step 1: Send a ping message from the second node to the first node
	// //fmt.Printf("Sending ping from second node to first node: %v\n", KademliaNode2.)
	// //KademliaNode2.Network.SendPingMessage(KademliaNode1, "ping:"+secondID.String()+":ping")

	// // Step 2: Wait and give time for the ping-pong interaction to complete
	// time.Sleep(1 * time.Second)

	// //node 1 is sending a lookupmesseage to node 2 to look up node 1, it then adds node 3, 4, and 5 to its routingtable
	// fileToStore := "gustav e bonk, honing is bonking. Messa with the honk and you get the bonkTHIS IS A STORED FILE HOPEFULLY"
	// hashedFile := kademlia.HashKademliaID(fileToStore)
	// KademliaNode1.StartTask(&hashedFile, "StoreValue", fileToStore)
	// // Step 4: Wait for lookUpContact to be processed and returnLookUpContact to be sent
	// time.Sleep(2 * time.Second)
	// // fmt.Println("(file: main) Value of iscontactinroutiongtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode2.RoutingTable.GetMe()))

	// fmt.Println("(file: main) Value of iscontactinroutiongtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode3.RoutingTable.GetMe()))
	// fmt.Println("(file: main) Value of iscontactinroutiongtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode4.RoutingTable.GetMe()))
	// fmt.Println("(file: main) Value of iscontactinroutiongtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode5.RoutingTable.GetMe()))
	// fmt.Println("(file: main) Value of iscontactinroutiongtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode6.RoutingTable.GetMe()))
	// fmt.Println("(file: main) Value of iscontactinroutiongtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode11.RoutingTable.GetMe()))

	// //fmt.Println("(file: main) Value of iscontactinroutiongtable:", KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode11.RoutingTable.GetMe()))

	time.Sleep(2 * time.Second)

}
