package kademlia

import "log"

// Kademlia node
type Kademlia struct {
	Network      Network
	RoutingTable RoutingTable
}

// NewKademlia creates and initializes a new Kademlia node
func NewKademlia(network Network, routingTable RoutingTable) *Kademlia {
	return &Kademlia{
		Network:      network,
		RoutingTable: routingTable,
	}
}

// Listen listens for incoming network messages and handles them
func (kademlia *Kademlia) Listen() {
	for {
		// Call the Network's Listen method and handle the response
		message, contactAddress, err := kademlia.Network.Listen()
		if err != nil {
			log.Printf("Error listening on the network: %v", err)
			continue
		}

		// Create a contact with the address only (ID will be handled later with the RoutingTable)
		contact := &Contact{Address: contactAddress}

		// Handle the received message
		kademlia.handleMessage(message, contact)
	}
}

// handleMessage processes the incoming messages and performs actions based on them
func (kademlia *Kademlia) handleMessage(message string, contact *Contact) {
	log.Printf("Received %s from %s", message, contact)
	switch message {
	case "ping":
		// Handle the "ping" message
		kademlia.handlePing(contact)
	case "pong":
		// Handle the "pong" message
		kademlia.handlePong(contact)
	default:
		// Handle other types of messages (e.g., FIND_NODE)
	}
}

// handlePing processes a "ping" message
func (kademlia *Kademlia) handlePing(contact *Contact) {
	log.Printf("Received ping from %s", contact)

	// Send a pong message back to the contact
	kademlia.Network.SendPongMessage(contact)
}

// handlePong processes a "pong" message
func (kademlia *Kademlia) handlePong(contact *Contact) {
	log.Printf("Received pong from %s", contact)

	// Update the routing table with the sender's contact information
	// This will be expanded later when integrating the RoutingTable
}

// LookupContact sends a FIND_NODE message
func (kademlia *Kademlia) LookupContact(targetContact *Contact, kademliaID string) {
	kademlia.Network.SendFindContactMessage(targetContact, kademliaID)
	kademlia.Network.Listen()
}
