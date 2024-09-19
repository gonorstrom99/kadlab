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

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
