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
		log.Printf("Kademlia processing message: '%s' from %s with ID: %s", msg.Command, msg.SenderAddress, msg.SenderID)

		// Create a contact using the sender's ID and address
		contact := &Contact{
			ID:      NewKademliaID(msg.SenderID), // Convert the sender's ID to a KademliaID
			Address: msg.SenderAddress,           // The sender's IP and port
		}

		// Handle different message types based on the "Command" field
		switch msg.Command {
		case "ping":
			// Respond with "pong" to a ping message
			kademlia.Network.SendPongMessage(contact)

		case "pong":
			// Log that a pong message was received
			log.Printf("Received pong from %s", msg.SenderAddress)

		case "lookUpContact":
			// Call the handleLookUpContact function, passing the contact
			kademlia.handleLookUpContact(contact, msg.CommandInfo)

		case "returnLookUpContact":
			// Handle the return lookup contact, passing commandInfo for processing
			kademlia.handleReturnLookUpContact(contact, msg.CommandInfo)

		case "findValue":
			// Handle the findValue command, using commandInfo as additional data
			kademlia.handleFindValue(contact, msg.CommandInfo)

		case "returnFindValue":
			// Handle the return of a found value, using commandInfo as additional data
			kademlia.handleReturnFindValue(contact, msg.CommandInfo)

		case "storeValue":
			// Handle storing a value, with commandInfo containing the value to be stored
			kademlia.handleStoreValue(contact, msg.CommandInfo)

		case "returnStoreValue":
			// Handle the return of a stored value
			kademlia.handleReturnStoreValue(contact, msg.CommandInfo)

		default:
			// Log unknown command types
			log.Printf("Received unknown message type '%s' from %s", msg.Command, msg.SenderAddress)
		}
	}
}

// handlePing processes a "ping" message
func (kademlia *Kademlia) handlePing(contact *Contact) {
	log.Printf("Received ping from %s", contact)

	// Send a pong message back to the contact
	kademlia.Network.SendPongMessage(contact)
}

func (kademlia *Kademlia) handleLookUpContact(contact *Contact, targetID string) {
	log.Printf("Handling lookUpContact from %s with target ID: %s", contact.Address, targetID)

	// Find the 3 closest contacts to the target ID in the routing table
	closestContacts := kademlia.RoutingTable.FindClosestContacts(NewKademliaID(targetID), 3)

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

	// Send the response message back to the requesting contact
	// The command for the response is 'returnLookUpContact'
	kademlia.Network.SendMessage(contact, fmt.Sprintf("returnLookUpContact:%s", responseMessage))

	log.Printf("Sent returnLookUpContact to %s with contacts: %s", contact.Address, responseMessage)
}

// handleReturnLookUpContact processes a "returnLookUpContact" message
func (kademlia *Kademlia) handleReturnLookUpContact(contact *Contact, commandInfo string) {
	log.Printf("Handling returnLookUpContact from %s", contact.Address)

	// Split the contact list by commas to get individual contact strings
	contactStrings := strings.Split(commandInfo, ",")

	// Iterate over the contact strings to parse and add them to the routing table
	for _, contactStr := range contactStrings {
		// Split each contact string into ID and address using ":"
		parts := strings.Split(contactStr, ":")
		if len(parts) != 2 {
			log.Printf("Invalid contact format: %s", contactStr)
			continue
		}

		// Create a new contact using the ID and the address
		newContact := NewContact(NewKademliaID(parts[0]), parts[1]) // parts[0] is the ID, parts[1] is the address

		// Add the contact to the routing table
		kademlia.RoutingTable.AddContact(newContact)
	}

	// Optionally, log that the contacts have been added to the routing table
	log.Printf("Added contacts to the routing table from returnLookUpContact message: %s", commandInfo)
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
