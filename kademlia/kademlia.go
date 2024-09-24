package kademlia

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"time"
)

const pongTimer = 5 //sekunder
var chPong chan string

type ponged struct {
	ID        string
	hasPonged bool
}

var pongList []ponged

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
func CreateKademliaNode(address string) *Kademlia {
	ID := NewRandomKademliaID()
	contact := NewContact(ID, address)
	routingTable := NewRoutingTable(contact)
	messageCh := make(chan Message)
	network := &Network{
		MessageCh: messageCh,
	}
	kademliaNode := NewKademlia(network, routingTable)
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

/*Message structure :
<command>:<senderID>:<commandInfo>*/

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

			id := kademlia.RoutingTable.me.ID.String()
			// Respond with "pong" to a ping message
			kademlia.Network.SendPongMessage(contact, "pong:"+id+":pong")

		case "pong":
			kademlia.handlePongMessage(contact)

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
	log.Printf("Received ping from %s", contact.Address)

	// Prepare the pong message with the appropriate format
	// The format will be "pong:<senderID>:<senderAddress>"
	id := kademlia.RoutingTable.me.ID.String()

	pongMessage := fmt.Sprintf("pong:%s:pong", id)

	// Send the pong message back to the contact
	kademlia.Network.SendMessage(contact, pongMessage)

	log.Printf("Sent pong to %s", contact.Address)
}

func (kademlia *Kademlia) handleLookUpContact(contact *Contact, targetID string) {
	log.Printf("(File: kademlia: Function: HandleLookupContact) Handling lookUpContact from %s with target ID: %s", contact.Address, targetID)

	// Find the bucketSize closest contacts to the target ID in the routing table
	closestContacts := kademlia.RoutingTable.FindClosestContacts(NewKademliaID(targetID), bucketSize)

	// Prepare the response message by concatenating the three closest contacts
	var responseMessage string
	myID := kademlia.RoutingTable.me.ID.String()
	for i, c := range closestContacts {

		contactStr := fmt.Sprintf("%s:%s", c.ID.String(), c.Address)

		// Append to the response message
		responseMessage += contactStr

		// Add a comma after each contact except the last one
		if i < len(closestContacts)-1 { //TODO check for off by one error
			responseMessage += ","
		}
	}

	// Send the response message back to the requesting contact
	// The command for the response is 'returnLookUpContact'

	kademlia.Network.SendMessage(contact, fmt.Sprintf("returnLookUpContact:%s:%s", myID, responseMessage))

	log.Printf("(File: kademlia: Function: HandleLookupContact) Sent returnLookUpContact to %s with contacts: %s", contact.Address, responseMessage)
}

// handleReturnLookUpContact processes a "returnLookUpContact" message
func (kademlia *Kademlia) handleReturnLookUpContact(contact *Contact, commandInfo string) {
	log.Printf("(File: kademlia: Function: HandleReturnLookupContact) Handling returnLookUpContact from %s", contact.Address)

	// Split the contact list by commas to get individual contact strings
	contactStrings := strings.Split(commandInfo, ",")

	// Iterate over the contact strings to parse and add them to the routing table
	for _, contactStr := range contactStrings {
		// Split each contact string into ID and address using ":"
		parts := strings.Split(contactStr, ":")
		//log.Printf("(File: kademlia: Function: HandleReturnLookupContact) len(parts):", len(parts))
		if len(parts) != 3 {
			log.Printf("(File: kademlia: Function: HandleReturnLookupContact) Invalid contact format: %s", contactStr)
			continue
		}

		// Create a new contact using the ID and the address
		newContact := NewContact(NewKademliaID(parts[0]), parts[1]) // parts[0] is the ID, parts[1] is the address
		// Add the contact to the routing table
		kademlia.updateRoutingTable(&newContact)
		//log.Printf("(File: kademlia: Function: HandleReturnLookupContact) called updateRoutingTable for a contact in returnLookUpContact message: %s", commandInfo)

	}

	// Optionally, log that the contacts have been added to the routing table
}

// handleFindValue processes a "findValue" message
func (kademlia *Kademlia) handleFindValue(contact *Contact, message string) {
	// TODO: Implement the logic for handling a "findValue" message
	log.Printf("(File: kademlia: Function: HandleFindValue) Handling findValue from %s", contact.Address)
}

// handleReturnFindValue processes a "returnFindValue" message
func (kademlia *Kademlia) handleReturnFindValue(contact *Contact, message string) {
	// TODO: Implement the logic for handling a "returnFindValue" message
	log.Printf("(File: kademlia: Function: HandleReturnFindValue) Handling returnFindValue from %s", contact.Address)
}

// handleStoreValue processes a "storeValue" message
func (kademlia *Kademlia) handleStoreValue(contact *Contact, message string) {
	// TODO: Implement the logic for handling a "storeValue" message
	log.Printf("(File: kademlia: Function: HandleStoreValue) Handling storeValue from %s", contact.Address)
}

// handleReturnStoreValue processes a "returnStoreValue" message
func (kademlia *Kademlia) handleReturnStoreValue(contact *Contact, message string) {
	// TODO: Implement the logic for handling a "returnStoreValue" message
	log.Printf("(File: kademlia: Function: HandleReturnStoreValue) Handling returnStoreValue from %s", contact.Address)
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

// CheckContactStatus pings a contact and returns true if its alive and false if not
func (kademlia *Kademlia) CheckContactStatus(contact *Contact) bool {

	id := kademlia.RoutingTable.me.ID.String()
	messageString := fmt.Sprintf("ping:%s:ping", id)
	kademlia.Network.SendPingMessage(contact, messageString)
	contactID := contact.ID.String()
	hasPonged := ponged{
		ID:        contactID,
		hasPonged: false,
	}

	pongList = append(pongList, hasPonged)
	chPong = make(chan string)
	timeOut := time.After(pongTimer * time.Second)
	waitTime := time.Second
	var pong bool = false //gets set to true if handlePongMessage is called (somehow) //has changed but is still used and should work plsplspls

	for {
		select {
		case ID := <-chPong:
			IDPonged := ponged{
				ID:        ID,
				hasPonged: false,
			}
			if ID == contactID {
				pong = true
				fmt.Println("The correct contact answered")
				removeFromList(pongList, findListIndex(pongList, contactID))
				return pong
			} else {
				fmt.Println("Pong recieved from incorrect contact")
				if slices.Contains(pongList, IDPonged) {
					IDPonged.hasPonged = true
				}
			}
		case <-timeOut:
			fmt.Println("Waited five seconds, contact presumed dead")
			return pong
		default:
			fmt.Println("still waiting for pong")
		}
		time.Sleep(waitTime)
		if hasPonged.hasPonged {
			pong = true
			return pong
		}
	}
}

func (kademlia *Kademlia) handlePongMessage(contact *Contact) {
	chPong <- contact.ID.String()
}

func removeFromList(s []ponged, i int) []ponged {
	if i == -1 {
		fmt.Println("index out of range")
		return s
	}
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findListIndex(s []ponged, ID string) int {
	for i, IDs := range s {
		if IDs.ID == ID {
			return i
		}
	}
	return -1
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
		if kademlia.CheckContactStatus(oldContact.Value.(*Contact)) == true {
			return false

		}

		//If not alive
		//delete the dead contact
		bucket.list.Remove(bucket.list.Back())
		return true

	}

	return true
}
