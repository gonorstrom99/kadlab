package main

import (
	"d7024e/kademlia"

	"fmt"
	//"time"
)

func main() {

	// skall endast skapa en bootstrap-nod
	fmt.Println("(file: main) Starting Kademlia nodes...")
	KademliaNode1 := kademlia.CreateKademliaNode("127.0.0.1:8000")
	KademliaNode1.Start()

	//hårdkodad bootstrap med icke-random KademliaID
	// kolla om du är boot
	//annars, lägg till (pinga) boot i routingtable

	//starta CLI

}
