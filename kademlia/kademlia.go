package kademlia

import (
	"fmt"
	"log"
	"strings"
)

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
			log.Printf("Received pong from %s", msg.Address)

		case "lookUpContact":
			kademlia.handleLookUpContact(contact)

		case "returnLookUpContact":
			kademlia.handleReturnLookUpContact(contact, msg.Content)

		case "findValue":
			kademlia.handleFindValue(contact, msg.Content)

		case "returnFindValue":
			kademlia.handleReturnFindValue(contact, msg.Content)

		case "storeValue":
			kademlia.handleStoreValue(contact, msg.Content)

		case "returnStoreValue":
			kademlia.handleReturnStoreValue(contact, msg.Content)

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
	log.Printf("Handling lookUpContact from %s", contact.Address)

	// Extract the target ID from the contact or message (assuming contact.ID here for simplicity)
	targetID := contact.ID

	// Find the 3 closest contacts to the target ID in the routing table
	closestContacts := kademlia.RoutingTable.FindClosestContacts(targetID, 3) // Find 3 closest contacts

	// Prepare the response message by concatenating the three closest contacts
	var responseMessage string
	for i, c := range closestContacts {

		contactStr := fmt.Sprintf("%s:%s", c.ID.String(), c.Address)

		// Append to the response message
		responseMessage += contactStr

		// Add a comma after each contact except the last one
		if i < len(closestContacts)-1 {
			responseMessage += ","
		}
	}

	// Send the response back to the original contact with the message format "returnLookUpContact:<contacts>"
	finalMessage := "returnLookUpContact:" + responseMessage
	kademlia.Network.sendMessage(contact, finalMessage)

	// Log the response message
	log.Printf("Sent 'returnLookUpContact' to %s with contacts: %s", contact.Address, responseMessage)
}

// handleReturnLookUpContact processes a "returnLookUpContact" message
func (kademlia *Kademlia) handleReturnLookUpContact(contact *Contact, message string) {
	log.Printf("Handling returnLookUpContact from %s", contact.Address)

	// Split the contact list by commas to get individual contact strings
	contactStrings := strings.Split(message, ",")

	// Iterate over the contact strings to parse and add them to the routing table
	for _, contactStr := range contactStrings {
		// Split each contact string into ID and address using ":"
		parts := strings.Split(contactStr, ":")
		if len(parts) != 2 {
			log.Printf("Invalid contact format: %s", contactStr)
			continue
		}

		// Create a new contact
		newContact := NewContact(NewKademliaID(parts[0]), parts[1]) // parts[0] is the ID, parts[1] is the address

		// Add the contact to the routing table
		kademlia.RoutingTable.AddContact(newContact)
	}

	// Optionally, log that the contacts have been added to the routing table
	log.Printf("Added contacts to the routing table from returnLookUpContact message: %s", message)
}

// handleFindValue processes a "findValue" message
func (kademlia *Kademlia) handleFindValue(contact *Contact, message string) {
	// TODO: Implement the logic for handling a "findValue" message
	log.Printf("Handling findValue from %s", contact.Address)
}

// handleReturnFindValue processes a "returnFindValue" message
func (kademlia *Kademlia) handleReturnFindValue(contact *Contact, message string) {
	// TODO: Implement the logic for handling a "returnFindValue" message
	log.Printf("Handling returnFindValue from %s", contact.Address)
}

// handleStoreValue processes a "storeValue" message
func (kademlia *Kademlia) handleStoreValue(contact *Contact, message string) {
	// TODO: Implement the logic for handling a "storeValue" message
	log.Printf("Handling storeValue from %s", contact.Address)
}

// handleReturnStoreValue processes a "returnStoreValue" message
func (kademlia *Kademlia) handleReturnStoreValue(contact *Contact, message string) {
	// TODO: Implement the logic for handling a "returnStoreValue" message
	log.Printf("Handling returnStoreValue from %s", contact.Address)
}
