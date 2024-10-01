package kademlia

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const pongTimer = 5 //sekunder

// Kademlia node
type Kademlia struct {
	Network      *Network
	RoutingTable *RoutingTable
	Tasks        []Task
}

// NewKademlia creates and initializes a new Kademlia node
func NewKademlia(network *Network, routingTable *RoutingTable, taskList []Task) *Kademlia {
	return &Kademlia{
		Network:      network,
		RoutingTable: routingTable,
		Tasks:        taskList,
	}
}

// CreateKademliaNode make a new kademlia node
func CreateKademliaNode(address string) *Kademlia {
	ID := NewRandomKademliaID()
	contact := NewContact(ID, address)
	routingTable := NewRoutingTable(contact)
	messageCh := make(chan Message)
	network := &Network{
		MessageCh: messageCh,
	}
	taskList := make([]Task, 0)
	kademliaNode := NewKademlia(network, routingTable, taskList)
	return kademliaNode
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
		log.Printf("Kademlia processing message: '%s' from %s with nodeID: %s and commandID: %s", msg.Command, msg.SenderAddress, msg.SenderID, msg.CommandID)

		// Create a contact using the sender's ID and address
		contact := &Contact{
			ID:      NewKademliaID(msg.SenderID), // Convert the sender's ID to a KademliaID
			Address: msg.SenderAddress,           // The sender's IP and port
		}

		// Handle different message types based on the "Command" field
		switch msg.Command {
		case "ping":

			// id := kademlia.RoutingTable.me.ID.String()
			// Respond with "pong" to a ping message
			kademlia.handlePing(contact, msg)

		case "pong":
			kademlia.handlePongMessage(contact, msg)

			// Log that a pong message was received
			log.Printf("Received pong from %s", msg.SenderAddress)

		case "lookUpContact":
			// Call the handleLookUpContact function, passing the contact
			kademlia.handleLookUpContact(contact, msg)

		case "returnLookUpContact":
			// Handle the return lookup contact, passing commandInfo for processing
			kademlia.handleReturnLookUpContact(contact, msg)

		case "findValue":
			// Handle the findValue command, using commandInfo as additional data
			kademlia.handleFindValue(contact, msg)

		case "returnFindValue":
			// Handle the return of a found value, using commandInfo as additional data
			kademlia.handleReturnFindValue(contact, msg)

		case "storeValue":
			// Handle storing a value, with commandInfo containing the value to be stored
			kademlia.handleStoreValue(contact, msg)

		case "returnStoreValue":
			// Handle the return of a stored value
			kademlia.handleReturnStoreValue(contact, msg)

		default:
			// Log unknown command types
			log.Printf("Received unknown message type '%s' from %s and commandID: %s", msg.Command, msg.SenderAddress, msg.CommandID)
		}
	}
}

// handlePing processes a "ping" message
func (kademlia *Kademlia) handlePing(contact *Contact, msg Message) {
	log.Printf("Received ping from %s", contact.Address)

	// Prepare the pong message with the appropriate format
	// The format will be "pong:<senderID>:<senderAddress>"
	id := kademlia.RoutingTable.me.ID.String()

	pongMessage := fmt.Sprintf("pong:%s:%s:pong", msg.CommandID, id)

	// Send the pong message back to the contact
	kademlia.Network.SendMessage(contact, pongMessage)

	log.Printf("Sent pong to %s", contact.Address)
}

func (kademlia *Kademlia) handlePongMessage(contact *Contact, msg Message) {

}

func (kademlia *Kademlia) handleLookUpContact(contact *Contact, msg Message) {
	log.Printf("(File: kademlia: Function: HandleLookupContact) Handling lookUpContact from %s with target ID: %s and commandID: %s", contact.Address, msg.Command, msg.CommandID)

	// Find the bucketSize closest contacts to the target ID in the routing table
	closestContacts := kademlia.RoutingTable.FindClosestContacts(NewKademliaID(msg.CommandInfo), bucketSize)

	// Prepare the response message by concatenating the three closest contacts
	var responseMessage string
	myID := kademlia.RoutingTable.me.ID.String()
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

	kademlia.Network.SendMessage(contact, fmt.Sprintf("returnLookUpContact:%s:%s:%s", myID, msg.CommandID, responseMessage))

	log.Printf("(File: kademlia: Function: HandleLookupContact) Sent returnLookUpContact to %s with contacts: %s", contact.Address, responseMessage)
}

// handleReturnLookUpContact processes a "returnLookUpContact" message
func (kademlia *Kademlia) handleReturnLookUpContact(contact *Contact, msg Message) {
	log.Printf("(File: kademlia: Function: HandleReturnLookupContact) Handling returnLookUpContact from %s", contact.Address)

	// Parse the CommandID to find the task
	commandID, err := strconv.Atoi(msg.CommandID)
	if err != nil {
		log.Printf("Invalid CommandID in message: %s", msg.CommandID)
		return
	}

	// Find the task by CommandID
	task, err := kademlia.FindTaskByCommandID(commandID)
	if err != nil {
		log.Printf("Task with CommandID %d not found", commandID)
		return
	}

	// Remove this contact from the task's waitingForReturns list
	kademlia.RemoveContactFromTask(commandID, *contact)

	// Split the contact list by commas to get individual contact strings
	contactStrings := strings.Split(msg.CommandInfo, ",")

	// Iterate over the contact strings to parse and add them to the routing table
	for _, contactStr := range contactStrings {
		// Split each contact string into ID and address using ":"
		parts := strings.Split(contactStr, ":")
		if len(parts) != 2 {
			log.Printf("(File: kademlia: Function: HandleReturnLookupContact) Invalid contact format: %s", contactStr)
			continue
		}

		// Create a new contact using the ID and the address
		newContact := NewContact(NewKademliaID(parts[0]), parts[1]) // parts[0] is the ID, parts[1] is the address

		// Add the new contact to the routing table
		kademlia.updateRoutingTable(&newContact)

		// Add the contact to the task's closest contacts if it's not already there
		if !kademlia.isContactInList(task.ClosestContacts, newContact) {
			task.ClosestContacts = append(task.ClosestContacts, newContact)
		}
	}

	// Now, find the closest uncontacted contact from the task's closest contacts
	for _, closestContact := range task.ClosestContacts {
		if !kademlia.isContactInList(task.ContactedNodes, closestContact) {
			// Send a lookUpContact message to this closest uncontacted contact
			lookupMessage := fmt.Sprintf("lookUpContact:%s:%d:%s", kademlia.RoutingTable.me.ID.String(), task.CommandID, task.TargetID.String())
			kademlia.Network.SendMessage(&closestContact, lookupMessage)

			// Mark this contact as contacted
			task.ContactedNodes = append(task.ContactedNodes, closestContact)
			task.WaitingForReturns = append(task.WaitingForReturns, WaitingContact{
				SentTime: time.Now(),
				Contact:  closestContact,
			})

			log.Printf("Sent lookUpContact to %s", closestContact.Address)
			return
		}
	}

	// If no uncontacted nodes remain, we can consider the task complete or handle accordingly
	log.Printf("All contacts have been contacted for task %d", task.CommandID)

	// Optionally, mark the task as complete if all responses have been received
	kademlia.MarkTaskAsCompleted(task.CommandID)
}

// isContactInList checks if a contact is already in a list of contacts
func (kademlia *Kademlia) isContactInList(contacts []Contact, contact Contact) bool {
	for _, c := range contacts {
		if c.ID.Equals(contact.ID) {
			return true
		}
	}
	return false
}

// handleFindValue processes a "findValue" message
func (kademlia *Kademlia) handleFindValue(contact *Contact, msg Message) {
	// TODO: Implement the logic for handling a "findValue" message
	log.Printf("(File: kademlia: Function: HandleFindValue) Handling findValue from %s", contact.Address)
}

// handleReturnFindValue processes a "returnFindValue" message
func (kademlia *Kademlia) handleReturnFindValue(contact *Contact, msg Message) {
	// TODO: Implement the logic for handling a "returnFindValue" message
	log.Printf("(File: kademlia: Function: HandleReturnFindValue) Handling returnFindValue from %s", contact.Address)
}

// handleStoreValue processes a "storeValue" message
func (kademlia *Kademlia) handleStoreValue(contact *Contact, msg Message) {
	// TODO: Implement the logic for handling a "storeValue" message
	log.Printf("(File: kademlia: Function: HandleStoreValue) Handling storeValue from %s", contact.Address)
}

// handleReturnStoreValue processes a "returnStoreValue" message
func (kademlia *Kademlia) handleReturnStoreValue(contact *Contact, msg Message) {
	// TODO: Implement the logic for handling a "returnStoreValue" message
	log.Printf("(File: kademlia: Function: HandleReturnStoreValue) Handling returnStoreValue from %s", contact.Address)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

// CheckContactStatus creates a ping task for a contact and adds it to the task list.
func (kademlia *Kademlia) CheckContactStatus(contact *Contact) {
	// Generate a new command ID for the ping task
	commandID := NewCommandID()

	// Create a message string for the ping
	messageString := fmt.Sprintf("ping:%s:%d:ping", kademlia.RoutingTable.me.ID.String(), commandID)

	// Create a new task for this ping operation (without ClosestContacts or WaitingForReturns)
	task := Task{
		CommandType: "ping",
		CommandID:   commandID,
		TargetID:    nil, // No specific target for ping
		StartTime:   time.Now(),
		// We omit ClosestContacts, ContactedNodes, and WaitingForReturns since ping doesn't need them.
	}

	// Add the task to the node's task list
	kademlia.Tasks = append(kademlia.Tasks, task)

	// Send the ping message to the contact
	kademlia.Network.SendPingMessage(contact, messageString)

	log.Printf("Ping task added for contact %s with CommandID %d", contact.Address, commandID)
}

func (kademlia *Kademlia) updateRoutingTable(contact *Contact) {
	//if it should be added it is done in the if, if the oldest node is
	//alive it is moved to the front in the else, if the oldest node is
	//dead it is removed in the "shouldContactBeAddedToRoutingTable".
	if kademlia.RoutingTable.me == *contact {
		return
	} else if kademlia.shouldContactBeAddedToRoutingTable(contact) == true {
		kademlia.RoutingTable.AddContact(*contact)
	} else {
		bucketIndex := kademlia.RoutingTable.getBucketIndex(contact.ID)
		bucket := kademlia.RoutingTable.buckets[bucketIndex]
		bucket.list.MoveToFront(bucket.list.Back())

	}
}

func (kademlia *Kademlia) shouldContactBeAddedToRoutingTable(contact *Contact) bool {
	// checks if the contact is already in it's respective bucket.
	if kademlia.RoutingTable.IsContactInRoutingTable(contact) == true {
		return true
	}

	// if bucket is full - ping oldest contact to check if alive
	bucketIndex := kademlia.RoutingTable.getBucketIndex(contact.ID)
	bucket := kademlia.RoutingTable.buckets[bucketIndex]
	if kademlia.RoutingTable.IsBucketFull(bucket) == true {
		//ping amandas function
		//if oldest contact alive {
		oldContact := bucket.list.Back()
		if kademlia.CheckContactStatus(&oldContact) == true {
			return false
		}

		//If not alive
		//delete the dead contact
		bucket.list.Remove(bucket.list.Back())
		return true

	}

	return true
}

func NewCommandID() int {
	return rand.Int()
}
