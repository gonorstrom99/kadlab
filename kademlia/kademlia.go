package kademlia

import (
	"fmt"
	"log"
	"time"
)

const pongTimer = 5 //sekunder
var pongList []string

// Kademlia node
type Kademlia struct {
	Network      *Network
	RoutingTable *RoutingTable
}

// NewKademlia creates and initializes a new Kademlia node
func NewKademlia(network *Network, routingTable *RoutingTable) *Kademlia {
	return &Kademlia{
		Network:      network,
		RoutingTable: routingTable,
	}
}

// Start starts the Kademlia node, processing incoming messages from the network channel
func (kademlia *Kademlia) Start() {
	// Start processing messages from the network's channel
	go func() {
		err := kademlia.Network.Listen(kademlia.RoutingTable.me)
		if err != nil {
			log.Printf("Error in network listener: %v", err)
		}
	}()
	go kademlia.processMessages()
}

// processMessages listens to the Network's channel and handles messages
func (kademlia *Kademlia) processMessages() {
	for msg := range kademlia.Network.MessageCh {
		log.Printf("Kademlia processing message: '%s' from %s with ID: %s", msg.Content, msg.Address, msg.ID)

		contact := &Contact{ID: NewKademliaID(msg.ID), Address: msg.Address}

		// Handle different message types
		switch msg.Content {
		case "ping":
			kademlia.Network.SendPongMessage(contact)

		case "pong":
			kademlia.handlePongMessage(contact)
			log.Printf("Received pong from %s", msg.Address)

		case "lookUpContact":
			kademlia.handleLookUpContact(contact)

		case "findValue":
			kademlia.handleFindValue(contact)

		case "storeValue":
			kademlia.handleStoreValue(contact)

		case "returnLookUpContact":
			kademlia.handleReturnLookUpContact(contact)

		case "returnFindValue":
			kademlia.handleReturnFindValue(contact)

		case "returnStoreValue":
			kademlia.handleReturnStoreValue(contact)

		default:
			log.Printf("Received unknown message type '%s' from %s", msg.Content, msg.Address)
		}
	}
}

// handlePing processes a "ping" message
func (kademlia *Kademlia) handlePing(contact *Contact) {
	log.Printf("Received ping from %s", contact)

	// Send a pong message back to the contact
	kademlia.Network.SendPongMessage(contact)
}

// handleLookUpContact processes a "lookUpContact" message
func (kademlia *Kademlia) handleLookUpContact(contact *Contact) {
	// TODO: Implement the logic for handling a "lookUpContact" message
	log.Printf("Handling lookUpContact from %s", contact.Address)
}

// handleFindValue processes a "findValue" message
func (kademlia *Kademlia) handleFindValue(contact *Contact) {
	// TODO: Implement the logic for handling a "findValue" message
	log.Printf("Handling findValue from %s", contact.Address)
}

// handleStoreValue processes a "storeValue" message
func (kademlia *Kademlia) handleStoreValue(contact *Contact) {
	// TODO: Implement the logic for handling a "storeValue" message
	log.Printf("Handling storeValue from %s", contact.Address)
}

// handleReturnLookUpContact processes a "returnLookUpContact" message
func (kademlia *Kademlia) handleReturnLookUpContact(contact *Contact) {
	// TODO: Implement the logic for handling a "returnLookUpContact" message
	log.Printf("Handling returnLookUpContact from %s", contact.Address)
}

// handleReturnFindValue processes a "returnFindValue" message
func (kademlia *Kademlia) handleReturnFindValue(contact *Contact) {
	// TODO: Implement the logic for handling a "returnFindValue" message
	log.Printf("Handling returnFindValue from %s", contact.Address)
}

// handleReturnStoreValue processes a "returnStoreValue" message
func (kademlia *Kademlia) handleReturnStoreValue(contact *Contact) {
	// TODO: Implement the logic for handling a "returnStoreValue" message
	log.Printf("Handling returnStoreValue from %s", contact.Address)
}

var chPong chan string

// CheckContactStatus pings a contact and returns true if its alive and false if not
func (kademlia *Kademlia) CheckContactStatus(contact *Contact) bool {
	kademlia.Network.SendPingMessage(contact)
	pongList = append(pongList, contact.ID.String())
	chPong = make(chan string)
	timeOut := time.After(pongTimer * time.Second)
	waitTime := time.Second
	var pong bool = false //gets set to true if handlePongMessage is called (somehow)

	//det var ngt mer jag skulle göra med pongList men har hjärnsläpp atm och kommer förhoppningsvis på det strax
	for {
		select {
		case ID := <-chPong:
			if ID == contact.ID.String() {
				pong = true
				fmt.Println("The correct contact answered")
				removeFromList(pongList, findListIndex(pongList, contact.ID.String()))
				return pong
			} else {
				fmt.Println("Pong recieved from incorrect contact")
			}
		case <-timeOut:
			fmt.Println("Waited five seconds, contact presumed dead")
			return pong
		default:
			fmt.Println("still waiting for pong")
		}
		time.Sleep(waitTime)
	}
}

func (kademlia *Kademlia) handlePongMessage(contact *Contact) {
	chPong <- contact.ID.String()
}

func removeFromList(s []string, i int) []string {
	if i == -1 {
		fmt.Println("index out of range")
		return s
	}
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findListIndex(s []string, ID string) int {
	for i, IDs := range s {
		if IDs == ID {
			return i
		}
	}
	return -1
}
