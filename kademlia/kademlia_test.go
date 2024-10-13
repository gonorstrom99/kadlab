package kademlia

import (
	"fmt"
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
	storeValueCalled := false

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

	kademlia.HandleStoreValueSpy = func(contact *Contact, msg Message) {
		storeValueCalled = true
	}

	// Prepare test messages
	testMessages := []Message{
		{Command: "ping", SenderID: "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", SenderAddress: "127.0.0.1:8080"},
		{Command: "pong", SenderID: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", SenderAddress: "127.0.0.1:8081"},
		{Command: "LookupContact", SenderID: "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB", SenderAddress: "127.0.0.1:8082"},
		{Command: "StoreValue", SenderID: "CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC", SenderAddress: "127.0.0.1:8083"},
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

	if !storeValueCalled {
		t.Errorf("Expected StoreValue handler to be called")
	}
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
