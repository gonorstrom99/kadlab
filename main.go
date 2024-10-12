package main

//$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
// entrypointnode docker run -it  --network host -e ID=container1 -e ADDRESS=127.0.0.1:8001 kadlab
// other nodes "docker run -it  --network host -e ID=container2 -e ADDRESS=127.0.0.1:8002 kadlab", 8002 to 8051 i  suppose

import (
	// "d7024e/kademlia"

	"bufio"
	"kadlab/kademlia"
	"time"

	"log"
	"os"
	"strings"

	"fmt"
)

func main() {
	fmt.Println("Docker node is running. Type 'put' followed by your command:")
	address := os.Getenv("ADDRESS")
	KademliaNode1 := kademlia.CreateKademliaNode(address)
	KademliaNode1.Start()
	time.Sleep(1 * time.Second)

	if address == "127.0.0.1:8001" {
		// KademliaNode1.Network.ID = *kademlia.NewKademliaID("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		log.Printf("in the if %s", KademliaNode1.RoutingTable.GetMe().Address)
		// KademliaNode1.Network.SendMessageWithAddress(address, "ping")
		// KademliaNode1.Network.SendMessageWithAddress(address, "ping:"+"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA:"+":ping")

	} else {
		contact := kademlia.NewContact(kademlia.NewKademliaID("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"), "127.0.0.1:8001")
		KademliaNode1.RoutingTable.AddContact(contact)
		KademliaNode1.StartTask(&KademliaNode1.Network.ID, "LookupContact", "")
	}

	// Create a new scanner to read input from stdin.
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if scanner.Scan() {
			// Read the input and trim any extra spaces.
			input := strings.TrimSpace(scanner.Text())

			// Split the input into words.
			parts := strings.Fields(input)

			// Ensure that we have exactly two parts.
			if len(parts) == 2 {
				command := parts[0]
				argument := parts[1]

				// Check if the input starts with the "put" command.
				if command == "put" {
					// Execute the command or handle the "put" command.
					// contact := kademlia.NewContact(kademlia.NewKademliaID("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"), "127.0.0.1:8001")
					// KademliaNode1.RoutingTable.AddContact(contact)
					fileToStore := argument
					hashedFile := kademlia.HashKademliaID(fileToStore)
					KademliaNode1.StartTask(&hashedFile, "StoreValue", fileToStore)

				} else if command == "show" {
					fileToShow := argument
					hashedFile := kademlia.HashKademliaID(fileToShow)
					value, found := KademliaNode1.Storage.GetValue(hashedFile)
					log.Printf("value found? %s%b", value, found)
				} else {
					fmt.Println("Unknown command. Type command followed by an argument.")
				}
			} else if strings.HasPrefix(input, "exit") {
				break
			} else {
				fmt.Println("Unknown command. Type command followed by an argument.")

			}
		} else {
			// If there's an error or EOF, break the loop.
			break
		}
	}

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

	//time.Sleep(2 * time.Second)

}
