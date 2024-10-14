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
	log.Printf("my address is: %s", KademliaNode1.RoutingTable.GetMe().Address)

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
					log.Printf("value found? %s%t", value, found)
				} else if command == "get" {
					// fileToStore := argument
					// hashedFile := kademlia.HashKademliaID(fileToStore)
					// hashedFile := argument
					hashedFile := kademlia.NewKademliaID(argument)
					KademliaNode1.StartTask(hashedFile, "FindValue", "")

					// KademliaNode1.StartTask(&hashedFile, "FindValue", fileToStore)

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

}
