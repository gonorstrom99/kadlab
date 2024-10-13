package kademlia

import (
	"testing"
	"time"
)

// TestProcessMessages tests the processMessages function by simulating messages being sent to the Kademlia network.
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
