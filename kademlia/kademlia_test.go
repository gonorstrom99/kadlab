package kademlia

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type MockNetwork struct {
	MessageCh chan Message
	Network   Network
}

// Mock SendMessage function (which would normally send messages over the network)
func (network *MockNetwork) SendMessage(contact *Contact, message string) {
	// For the mock, we can log or store the message, or simply do nothing.
	fmt.Printf("MockNetwork: Message sent to %s with content: %s\n", contact.Address, message)
}

func TestProcessMessagesWithSpies(t *testing.T) {
	// Create a mock network
	mockNetwork := &MockNetwork{
		MessageCh: make(chan Message, 10), // Initialize a buffered message channel
		Network: Network{
			MessageCh: make(chan Message, 10), // Ensure the embedded Network also has a channel
		},
	}

	// Create a Kademlia instance with the mock network
	kademlia := &Kademlia{
		Network: &mockNetwork.Network,
		RoutingTable: &RoutingTable{
			me: Contact{
				ID:      NewRandomKademliaID(),
				Address: "127.0.0.1:8000",
			},
		},
	}

	// Initialize spy variables to track whether each handler was called
	pingCalled := false
	pongCalled := false
	lookupCalled := false
	returnlookupcontactCalled := false
	findValueCalled := false
	returnFindValueCalled := false
	storeValueCalled := false
	returnStoreValueCalled := false

	// Inject spies
	kademlia.HandlePingSpy = func(contact *Contact, msg Message) {
		pingCalled = true
	}

	kademlia.HandlePongSpy = func(contact *Contact, msg Message) {
		pongCalled = true
	}

	kademlia.HandleLookupSpy = func(contact *Contact, msg Message) {
		lookupCalled = true
	}

	kademlia.HandleReturnLookupContactSpy = func(contact *Contact, msg Message) {
		returnlookupcontactCalled = true
	}

	kademlia.HandleFindValueSpy = func(contact *Contact, msg Message) {
		findValueCalled = true
	}

	kademlia.HandleReturnFindValueSpy = func(contact *Contact, msg Message) {
		returnFindValueCalled = true
	}

	kademlia.HandleStoreValueSpy = func(contact *Contact, msg Message) {
		storeValueCalled = true
	}

	kademlia.HandleReturnStoreValueSpy = func(contact *Contact, msg Message) {
		returnStoreValueCalled = true
	}

	// Prepare test messages for different commands
	testMessages := []Message{
		{Command: "ping", SenderID: "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", SenderAddress: "127.0.0.1:8080"},
		{Command: "pong", SenderID: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", SenderAddress: "127.0.0.1:8081"},
		{Command: "LookupContact", SenderID: "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB", SenderAddress: "127.0.0.1:8082"},
		{Command: "returnLookupContact", SenderID: "CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC", SenderAddress: "127.0.0.1:8083"},
		{Command: "FindValue", SenderID: "DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD", SenderAddress: "127.0.0.1:8084"},
		{Command: "returnFindValue", SenderID: "EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE", SenderAddress: "127.0.0.1:8085"},
		{Command: "StoreValue", SenderID: "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", SenderAddress: "127.0.0.1:8086"},
		{Command: "returnStoreValue", SenderID: "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFAAAAAA", SenderAddress: "127.0.0.1:8087"},
	}

	// Feed the test messages into the message channel
	go func() {
		for _, msg := range testMessages {
			kademlia.Network.MessageCh <- msg
		}
	}()

	// Start processing messages
	go func() {
		kademlia.processMessages()
	}()

	// Wait for messages to be processed
	time.Sleep(100 * time.Millisecond)

	// Check if the correct handlers were called
	if !pingCalled {
		t.Errorf("Expected ping handler to be called")
	}

	if !pongCalled {
		t.Errorf("Expected pong handler to be called")
	}

	if !lookupCalled {
		t.Errorf("Expected LookupContact handler to be called")
	}

	if !returnlookupcontactCalled {
		t.Errorf("Expected returnLookupContact handler to be called")
	}

	if !findValueCalled {
		t.Errorf("Expected FindValue handler to be called")
	}

	if !returnFindValueCalled {
		t.Errorf("Expected returnFindValue handler to be called")
	}

	if !storeValueCalled {
		t.Errorf("Expected StoreValue handler to be called")
	}

	if !returnStoreValueCalled {
		t.Errorf("Expected returnStoreValue handler to be called")
	}
}

func TestMin(t *testing.T) {
	a := 3
	b := 2
	if min(a, b) != b {
		t.Fatalf("expected b to be less than a given b is 2 and a is 3")
	}
}
func TestCreateKademliaNode(t *testing.T) {
	KademliaNode1 := CreateKademliaNode("127.0.0.1:8050")
	if KademliaNode1.RoutingTable.me.Address != "127.0.0.1:8050" {
		t.Fatalf("expected the nodes address to be 127.0.0.1:8050")
	}
	KademliaNode2 := CreateKademliaNode("127.0.0.1:8001")
	print(KademliaNode2.Network.ID.String())
	if KademliaNode2.Network.ID.String() != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Fatalf("expected the nodes ID to be AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	}

}
func TestProcessMessages(t *testing.T) {
	// Create a mock network and Kademlia node without using Conn
	mockNetwork := &Network{
		ID:        *NewRandomKademliaID(),
		MessageCh: make(chan Message, 10), // Channel with buffer size 10 for messages
	}

	// Create a Kademlia instance with the mock network
	kademlia := &Kademlia{
		Network: mockNetwork,
	}

	// Start processing messages in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic: %v", r)
			}
		}()
		kademlia.processMessages()
	}()

	// Prepare test messages for different commands
	testMessages := []Message{
		{Command: "ping", SenderID: "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", SenderAddress: "127.0.0.1:8080"},
		{Command: "pong", SenderID: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", SenderAddress: "127.0.0.1:8081"},
		{Command: "LookupContact", SenderID: "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB", SenderAddress: "127.0.0.1:8082"},
		{Command: "StoreValue", SenderID: "CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC", SenderAddress: "127.0.0.1:8083"},
		{Command: "InvalidCommand", SenderID: "DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD", SenderAddress: "127.0.0.1:8084"}, // Invalid command
	}

	// Feed each message into the mock MessageCh
	for _, msg := range testMessages {
		func() {
			// Wrap each send operation in defer/recover to handle panics
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Recovered from panic while sending message: %v", r)
				}
			}()
			mockNetwork.MessageCh <- msg
		}()
	}

	// Allow time for messages to be processed
	time.Sleep(500 * time.Millisecond)

	// Optionally, you could assert conditions here, but the goal is to ensure that panics are handled and messages are processed.
}

// TestProcessMessagesTimeout tests that processMessages handles the TTL timeout case
func TestProcessMessagesTimeout(t *testing.T) {
	// Create a mock Kademlia instance with a mock network
	mockNetwork := &Network{
		MessageCh: make(chan Message, 1), // Buffer size of 1 to simulate channel input
	}
	kademlia := &Kademlia{
		Network: mockNetwork,
	}

	// Run processMessages in a goroutine
	go func() {
		kademlia.processMessages()
	}()

	// Wait longer than the TTL (assuming TTL is set to 100 milliseconds)
	time.Sleep(TTL * 2 * time.Millisecond)

	// At this point, the timeout case should have triggered, leading to kademlia.checkTTLs() being called
	// You can assert that behavior if `checkTTLs` has side effects you can observe or if you can mock `checkTTLs`.
}

func TestIsContactInList(t *testing.T) {
	// Create some sample contacts
	contact1 := NewContact(NewKademliaID("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"), "127.0.0.1:8000")
	contact2 := NewContact(NewKademliaID("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"), "127.0.0.1:8001")
	contact3 := NewContact(NewKademliaID("CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC"), "127.0.0.1:8002")

	// Create a Kademlia instance (mock network, routing table, etc.)
	kademlia := &Kademlia{}

	// Test case 1: Contact is in the list
	contacts := []Contact{contact1, contact2}
	if !kademlia.isContactInList(contacts, contact1) {
		t.Errorf("Expected contact1 to be in the list, but it wasn't found")
	}

	// Test case 2: Contact is not in the list
	if kademlia.isContactInList(contacts, contact3) {
		t.Errorf("Expected contact3 not to be in the list, but it was found")
	}

	// Test case 3: Empty contact list
	contacts = []Contact{}
	if kademlia.isContactInList(contacts, contact1) {
		t.Errorf("Expected contact1 not to be found in an empty list")
	}
}
func TestUpdateRoutingTable(t *testing.T) {
	// Setup a Kademlia instance
	me := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	routingTable := NewRoutingTable(me)
	network := &Network{
		ID:        *me.ID,
		MessageCh: make(chan Message),
	}
	storage := NewStorage()
	kademlia := NewKademlia(network, routingTable, []Task{}, storage)

	// Create a new contact (different from self)
	newContact := NewContact(NewRandomKademliaID(), "127.0.0.1:8001")

	// Test case 1: Add new contact to the routing table (bucket not full)
	kademlia.updateRoutingTable(&newContact)

	// Verify that the new contact was added
	if !routingTable.IsContactInRoutingTable(&newContact) {
		t.Errorf("Expected newContact to be added to the routing table, but it wasn't")
	}

	// Test case 2: Don't add yourself to the routing table
	kademlia.updateRoutingTable(&me)

	// Verify that the node didn't add itself
	if routingTable.IsContactInRoutingTable(&me) {
		t.Errorf("Expected not to add self to the routing table, but it was added")
	}

	// Test case 3: Move existing contact to the front of the bucket
	// Add the same contact again
	kademlia.updateRoutingTable(&newContact)

	// Verify that the contact was moved to the front (most recently used)
	bucketIndex := routingTable.getBucketIndex(newContact.ID)
	bucket := routingTable.buckets[bucketIndex]

	if bucket.list.Front().Value.(Contact).ID != newContact.ID {
		t.Errorf("Expected newContact to be moved to the front of the bucket, but it wasn't")
	}

	// Test case 4: Bucket is full, ping the oldest contact
	// Fill the bucket with contacts to simulate a full bucket
	for i := 0; i < bucketSize; i++ {
		contact := NewContact(NewRandomKademliaID(), "127.0.0.1:800"+strconv.Itoa(i))
		bucket.AddContact(contact)
	}

	// Add one more contact, which should cause the oldest contact to be pinged

	// Verify that the oldest contact was pinged (the oldest one should still be in the back)

}
