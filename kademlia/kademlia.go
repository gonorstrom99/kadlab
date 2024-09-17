package kademlia

import "log"

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
	go kademlia.processMessages()
}

// processMessages listens to the Network's channel and handles messages
func (kademlia *Kademlia) processMessages() {
	for msg := range kademlia.Network.MessageCh {
		log.Printf("Kademlia processing message: '%s' from %s", msg.Content, msg.Address)
		contact := &Contact{Address: msg.Address}

		// Handle different message types
		switch msg.Content {
		case "ping":
			kademlia.Network.SendPongMessage(contact) // Respond with "pong"
		case "pong":
			// Handle "pong" message (e.g., update routing table)
			log.Printf("Received pong from %s", msg.Address)
			// Add more cases for other Kademlia-specific messages
		}
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
	//kademlia.Network.SendFindContactMessage(targetContact, kademliaID)
}
