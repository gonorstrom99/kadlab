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
		log.Printf("Kademlia processing message: '%s' from %s", msg.Content, msg.Address)
		contact := &Contact{Address: msg.Address}

		// Handle different message types
		switch msg.Content {
		case "ping":
			kademlia.handlePing(contact) // Respond with "pong"
		case "pong":
			kademlia.handlePong(contact)
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

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
