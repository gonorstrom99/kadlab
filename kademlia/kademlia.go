package kademlia

import (
	"fmt"
	"log"
	"math/rand"
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

// use newCommandID to get a command ID, even though it's just a random int
// för att kolla om ett ID finns i listan, använd slices.contains(listan, ID)
// kan ha flera listor för olika commands om man vill (en för lookupcontact etc)
var commandIDlist []int

var pongList []ponged

// Kademlia node
type Kademlia struct {
	Network      *Network
	RoutingTable *RoutingTable
	// Tasks        []Task
}

// NewKademlia creates and initializes a new Kademlia node
func NewKademlia(network *Network, routingTable *RoutingTable) *Kademlia {
	return &Kademlia{
		Network:      network,
		RoutingTable: routingTable,
	}
}

// CreateKademliaNode make a new kademlia node
func CreateKademliaNode(address string) *Kademlia {
	ID := NewRandomKademliaID()
	contact := NewContact(ID, address)
	routingTable := NewRoutingTable(contact)
	messageCh := make(chan Message)
	network := &Network{
		ID:        *ID,
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

/*
Message structure :
<command>:<senderID>:<commandID><commandInfo>
*/
func (kademlia *Kademlia) StartLookUpContact(recipient Contact, lookupTarget Contact) {

	lookupMessage := fmt.Sprintf("lookUpContact:%s:%d:%s", kademlia.Network.ID.String(), NewCommandID(), lookupTarget.ID.String())
	//  lookupMessage := fmt.Sprintf("lookUpContact:%s:%d:%s", KademliaNode1.Network.ID.String(), kademlia.NewCommandID(), KademliaNode1.Network.ID.String())

	kademlia.Network.SendMessage(&recipient, lookupMessage)
	//  KademliaNode1.Network.SendMessage(KademliaNode2.RoutingTable.GetMe(), lookupMessage)

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
	chPong <- contact.ID.String()
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
	kademlia.updateRoutingTable(contact)

	log.Printf("(File: kademlia: Function: HandleLookupContact) Sent returnLookUpContact to %s with contacts: %s", contact.Address, responseMessage)
}

// handleReturnLookUpContact processes a "returnLookUpContact" message
func (kademlia *Kademlia) handleReturnLookUpContact(contact *Contact, msg Message) {
	log.Printf("(File: kademlia: Function: HandleReturnLookupContact) Handling returnLookUpContact from %s", contact.Address)

	// Split the contact list by commas to get individual contact strings
	contactStrings := strings.Split(msg.CommandInfo, ",")

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
		//TODO if com.id in commandlist run lookupcontact again
		//if command exists
		//lookupcontact(newcontact, targetContact)
		//log.Printf("(File: kademlia: Function: HandleReturnLookupContact) called updateRoutingTable for a contact in returnLookUpContact message: %s", commandInfo)

	}
	// Optionally, log that the contacts have been added to the routing table
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
	var pong bool = false //gets set to true if handlePongMessage is called (somehow)

	//det var ngt mer jag skulle göra med pongList men har hjärnsläpp atm och kommer förhoppningsvis på det strax
	//ponglist finns specifikt för att för att hantera ifall pongs kommer i "fel" ordning, om man vill kolla statusen på flera kontakter

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
		if hasPonged.hasPonged == true {
			pong = true
			return pong
		}
	}
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

func NewCommandID() int {
	return rand.Int()
}

func removeFromCommandIDList(ID int) []int {
	for i, IDs := range commandIDlist {
		if IDs == ID {
			if i == -1 {
				fmt.Println("index out of range")
				return commandIDlist
			}
			// Replace the current element with the last one and then truncate the slice
			commandIDlist[i] = commandIDlist[len(commandIDlist)-1]
			return commandIDlist[:len(commandIDlist)-1]
		}
	}
	return commandIDlist
}
