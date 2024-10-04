package kademlia

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

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

// processMessages listens to the Network's channel and handles messages or performs other tasks if no messages are available
func (kademlia *Kademlia) processMessages() {
	for {
		select {
		case msg := <-kademlia.Network.MessageCh: // If there's a message in the channel
			log.Printf("Kademlia processing message: '%s' from %s with nodeID: %s and commandID: %s", msg.Command, msg.SenderAddress, msg.SenderID, msg.CommandID)

			// Create a contact using the sender's ID and address
			contact := &Contact{
				ID:      NewKademliaID(msg.SenderID), // Convert the sender's ID to a KademliaID
				Address: msg.SenderAddress,           // The sender's IP and port
			}

			// Handle different message types based on the "Command" field
			switch msg.Command {
			case "ping":
				// Respond with "pong" to a ping message
				kademlia.handlePing(contact, msg)

			case "pong":
				kademlia.handlePongMessage(contact, msg)
				log.Printf("Received pong from %s", msg.SenderAddress)

			case "lookUpContact":
				kademlia.handleLookUpContact(contact, msg)

			case "returnLookUpContact":
				kademlia.handleReturnLookUpContact(contact, msg)

			case "findValue":
				kademlia.handleFindValue(contact, msg)

			case "returnFindValue":
				kademlia.handleReturnFindValue(contact, msg)

			case "storeValue":
				kademlia.handleStoreValue(contact, msg)

			case "returnStoreValue":
				kademlia.handleReturnStoreValue(contact, msg)

			default:
				// Log unknown command types
				log.Printf("Received unknown message type '%s' from %s and commandID: %s", msg.Command, msg.SenderAddress, msg.CommandID)
			}

		case <-time.After(100 * time.Millisecond): // If no message is received after 100 ms
			// Perform some other action when no messages are received
			kademlia.checkTTLs()
		}
	}
}

// checkTTLs checks the task list for contacts that haven't responded within the TTL limit
func (kademlia *Kademlia) checkTTLs() {
	ttl := 255 * time.Millisecond // Set the TTL limit

	for _, task := range kademlia.Tasks {
		var updatedWaitingForReturns []WaitingContact // To store contacts that are still within TTL

		for _, waitingContact := range task.WaitingForReturns {
			// Check if the contact has exceeded the TTL
			if time.Since(waitingContact.SentTime) > ttl {
				// Contact has timed out, remove it from WaitingForReturns
				log.Printf("Contact %s in task %d has timed out", waitingContact.Contact.ID.String(), task.CommandID)
				if task.CommandType == "ping" {
					bucketIndex := kademlia.RoutingTable.getBucketIndex(task.TargetID)
					bucket := kademlia.RoutingTable.buckets[bucketIndex]
					bucket.AddContact(task.ReplaceContact)

					kademlia.RemoveTask(task.CommandID)
					continue
				}
				// Try to find the next closest contact that hasn't been contacted yet
				for _, closestContact := range task.ClosestContacts {
					if !kademlia.isContactInList(task.ContactedNodes, closestContact) {
						// Send a new lookup to this uncontacted node
						lookupMessage := fmt.Sprintf("lookUpContact:%s:%d:%s", kademlia.RoutingTable.me.ID.String(), task.CommandID, task.TargetID.String())
						kademlia.Network.SendMessage(&closestContact, lookupMessage)

						// Mark this contact as contacted and add it to WaitingForReturns
						task.ContactedNodes = append(task.ContactedNodes, closestContact)
						task.WaitingForReturns = append(task.WaitingForReturns, WaitingContact{
							SentTime: time.Now(),
							Contact:  closestContact,
						})

						log.Printf("Sent lookUpContact to %s", closestContact.Address)
						break // Stop looking for the next contact once one is found
					}
				}
			} else {
				// Contact is still within TTL, keep it in the list
				updatedWaitingForReturns = append(updatedWaitingForReturns, waitingContact)
			}
		}

		// Update the task's WaitingForReturns list to only include those that haven't timed out
		task.WaitingForReturns = updatedWaitingForReturns

		// Optionally, you could check if the task is now complete and remove it
		if len(task.WaitingForReturns) == 0 {
			kademlia.MarkTaskAsCompleted(task.CommandID)
		}
	}
}

// Helper function to check if a contact is in a given list
func (kademlia *Kademlia) isContactInList(contacts []Contact, contact Contact) bool {
	for _, c := range contacts {
		if c.ID.Equals(contact.ID) {
			return true
		}
	}
	return false
}

// handlePing processes a "ping" message
func (kademlia *Kademlia) handlePing(contact *Contact, msg Message) {
	log.Printf("Received ping from %s", contact.Address)

	// Prepare the pong message with the appropriate format
	// The format will be "pong:<senderID>:<senderAddress>:pong"
	id := kademlia.RoutingTable.me.ID.String()

	pongMessage := fmt.Sprintf("pong:%s:%s:pong", msg.CommandID, id)

	// Send the pong message back to the contact
	kademlia.Network.SendMessage(contact, pongMessage)

	log.Printf("Sent pong to %s", contact.Address)
}

func (kademlia *Kademlia) handlePongMessage(contact *Contact, msg Message) {
	// Parse the CommandID from the message
	commandID, err := strconv.Atoi(msg.CommandID)
	if err != nil {
		log.Printf("Invalid CommandID in pong message: %s", msg.CommandID)
		return
	}

	// Try to find the task with the matching CommandID (you can omit `task` if it's not needed)
	if _, err := kademlia.FindTaskByCommandID(commandID); err != nil {
		log.Printf("Task with CommandID %d not found", commandID)
		return
	}

	// Task is found, now remove it
	kademlia.RemoveTask(commandID)
	log.Printf("Task with CommandID %d removed after pong received", commandID)

	// Use AddContact from the routing table to handle updating or adding the contact
	kademlia.RoutingTable.AddContact(*contact)
	log.Printf("Contact %s added or updated in the routing table", contact.Address)
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

func (kademlia *Kademlia) updateRoutingTable(newContact *Contact) {
	//if it should be added it is done in the if, if the oldest node is
	//alive it is moved to the front in the else, if the oldest node is
	if kademlia.RoutingTable.me == *newContact {
		return
	}
	if kademlia.RoutingTable.IsContactInRoutingTable(newContact) {
		bucketIndex := kademlia.RoutingTable.getBucketIndex(newContact.ID)
		bucket := kademlia.RoutingTable.buckets[bucketIndex]
		bucket.AddContact(*newContact)
	}

	// if bucket is full - ping oldest contact to check if alive and creat ping task
	bucketIndex := kademlia.RoutingTable.getBucketIndex(newContact.ID)
	bucket := kademlia.RoutingTable.buckets[bucketIndex]
	if kademlia.RoutingTable.IsBucketFull(bucket) {
		// Get the oldest contact (back of the list)
		oldestElement := bucket.list.Back()
		if oldestElement != nil {
			// Extract *Contact from *list.Element
			oldContact := oldestElement.Value.(*Contact)

			// Ping the oldest contact to check if it's alive
			kademlia.CheckContactStatus(oldContact, newContact)
		}
	}
}

// CheckContactStatus creates a ping task for a contact and adds it to the task list.
func (kademlia *Kademlia) CheckContactStatus(oldContact *Contact, newContact *Contact) {
	// Generate a new command ID for the ping task
	commandID := NewCommandID()

	// Create a message string for the ping
	messageString := fmt.Sprintf("ping:%s:%d:ping", kademlia.RoutingTable.me.ID.String(), commandID)

	// Create a new task for this ping operation (without ClosestContacts or WaitingForReturns)
	task := Task{
		CommandType: "ping",
		CommandID:   commandID,
		TargetID:    oldContact.ID,
		StartTime:   time.Now(),
		// We omit ClosestContacts, ContactedNodes, and WaitingForReturns since ping doesn't need them.
		ReplaceContact: *newContact,
	}

	// Add the task to the node's task list
	kademlia.Tasks = append(kademlia.Tasks, task)

	// Send the ping message to the contact
	kademlia.Network.SendPingMessage(oldContact, messageString)

	log.Printf("Ping task added for contact %s with CommandID %d", oldContact.Address, commandID)
}

// NewCommandID give a new command ID random int
func NewCommandID() int {
	return rand.Int()
}
