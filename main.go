// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"d7024e/kademlia"
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Println("Pretending to run the kademlia app...")
	// Using stuff from the kademlia package here. Something like...
	id := kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	contact := kademlia.NewContact(id, "localhost:8000")
	fmt.Println(contact.String())
	fmt.Printf("%v\n", contact)
	//networkthing := &kademlia.Network{}
	node1 := &kademlia.Kademlia{
		NodeID: "node1",
		IP:     "127.0.0.1",
		Port:   8000,
		//network: *networkthing,

	}
	node2 := &kademlia.Kademlia{
		NodeID: "node2",
		IP:     "127.0.0.1",
		Port:   8001,
		//network: *network,

	}
	/*node3 := &kademlia.Network{
		NodeID: "node2",
		IP:     "127.0.0.1",
		Port:   8002,
	}
	node4 := &kademlia.Network{
		NodeID: "node2",
		IP:     "127.0.0.1",
		Port:   8003,
	}*/

	go func() {
		if err := node1.Listen(); err != nil {
			log.Fatal("Node 1 Listen error:", err)
		}
	}()

	go func() {
		if err := node2.Listen(); err != nil {
			log.Fatal("Node 2 Listen error:", err)
		}
	}()
	/*go func() {
		if err := node3.Listen(); err != nil {
			log.Fatal("Node 3 Listen error:", err)
		}
	}()

	go func() {
		if err := node4.Listen(); err != nil {
			log.Fatal("Node 4 Listen error:", err)
		}
	}()*/
	// Wait a moment for both listeners to initialize
	time.Sleep(2 * time.Second)

	// Ping from node1 to node2
	if err := node1.Ping(node2); err != nil {
		log.Fatal("Ping failed:", err)
	}
	/*
		if err := node1.Ping("127.0.0.1", 8001); err != nil {
			log.Fatal("Ping failed:", err)
		}
		if err := node1.Ping("127.0.0.1", 8001); err != nil {
			log.Fatal("Ping failed:", err)
		}
		if err := node1.Ping("127.0.0.1", 8001); err != nil {
			log.Fatal("Ping failed:", err)
		}*/
}