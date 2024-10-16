package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"
	"time"
)

// TestHashKademliaID tests if the KademliaID is correctly hashed using SHA-1
func TestHashKademliaID(t *testing.T) {
	// Input string for the test
	input := "this is a test input"

	// Manually compute the expected SHA-1 hash for the input string
	hasher := sha1.New()
	hasher.Write([]byte(input))
	expectedHash := hasher.Sum(nil)

	// Call HashKademliaID to generate the KademliaID hash
	kademliaID := HashKademliaID(input)

	// Check if the length of the hashed ID is correct (should be IDLength bytes, which is 20 bytes for SHA-1)
	if len(kademliaID) != IDLength {
		t.Fatalf("Expected KademliaID length to be %d, but got %d", IDLength, len(kademliaID))
	}

	// Compare the result with the expected SHA-1 hash byte-by-byte
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != expectedHash[i] {
			t.Errorf("Expected byte %d to be %x, but got %x", i, expectedHash[i], kademliaID[i])
		}
	}
}

// TestNewKademliaID verifies that a new KademliaID is created correctly from a string
func TestNewKademliaID(t *testing.T) {
	// Updated input to be 40 hex characters (20 bytes) long
	data := "FFFFFFFFFF0000000000FFFFFFFFFF0000000000" // 40 characters = 20 bytes
	kademliaID := NewKademliaID(data)

	expectedBytes, _ := hex.DecodeString(data)

	// Ensure the expectedBytes length matches the IDLength
	if len(expectedBytes) != IDLength {
		t.Fatalf("Expected length of decoded bytes to be %d, but got %d", IDLength, len(expectedBytes))
	}

	// Verify each byte in kademliaID matches the expected value
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != expectedBytes[i] {
			t.Errorf("Expected byte %d to be %x, but got %x", i, expectedBytes[i], kademliaID[i])
		}
	}
}

// TestNewRandomKademliaID verifies that a random KademliaID is created and has the expected length
func TestNewRandomKademliaID(t *testing.T) {
	kademliaID := NewRandomKademliaID()

	if len(kademliaID) != IDLength {
		t.Errorf("Expected KademliaID length to be %d, but got %d", IDLength, len(kademliaID))
	}
}

// TestLess checks the Less function for KademliaID comparison
func TestLess(t *testing.T) {

	id1 := NewKademliaID("0000000000000000000000000000000000000000")
	id2 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id3 := NewKademliaID("0000000000000000000000000000000000000001") // Close to id1 for edge case

	// Test if id1 is less than id2
	if !id1.Less(id2) {
		t.Errorf("Expected id1 to be less than id2")
	}

	// Test if id2 is not less than id1
	if id2.Less(id1) {
		t.Errorf("Did not expect id2 to be less than id1")
	}

	// Test if id1 is less than id3
	if !id1.Less(id3) {
		t.Errorf("Expected id1 to be less than id3")
	}

	// Test if id3 is not less than id1
	if id3.Less(id1) {
		t.Errorf("Did not expect id3 to be less than id1")
	}

	// Test if identical KademliaID's return false for both comparisons
	id4 := NewKademliaID("0000000000000000000000000000000000000000")
	if id1.Less(id4) || id4.Less(id1) {
		t.Errorf("Expected identical KademliaIDs to be equal in comparison")
	}
}

// TestEquals checks if two KademliaIDs are correctly identified as equal
func TestEquals(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF0000000000FFFFFFFF0000000000FFFF")
	id2 := NewKademliaID("FFFFFFFF0000000000FFFFFFFF0000000000FFFF")
	id3 := NewKademliaID("00000000FFFFFFFFFFFFFFFF0000000000000000")

	// Test if id1 equals id2 (same value)
	if !id1.Equals(id2) {
		t.Errorf("Expected id1 to equal id2")
	}

	// Test if id1 does not equal id3 (different value)
	if id1.Equals(id3) {
		t.Errorf("Did not expect id1 to equal id3")
	}

	// Test if id1 equals itself (self-comparison)
	if !id1.Equals(id1) {
		t.Errorf("Expected id1 to equal itself")
	}

	// Test if id3 equals itself (self-comparison)
	if !id3.Equals(id3) {
		t.Errorf("Expected id3 to equal itself")
	}
}

// TestCalcDistance checks the CalcDistance function for XOR distance calculation
func TestCalcDistance(t *testing.T) {
	// Use full-length 20-byte Kademlia IDs (40 hexadecimal characters)
	id1 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF") // All bits set in the first half
	id2 := NewKademliaID("0000000000000000000000000000000000000000") // All bits set in the second half

	// The expected XOR distance (should be all F's as the difference between id1 and id2)
	expectedDistance := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	// Calculate the XOR distance between id1 and id2
	calculatedDistance := id1.CalcDistance(id2)

	// Check if the calculated distance matches the expected distance
	if !calculatedDistance.Equals(expectedDistance) {
		t.Errorf("Expected distance to be %s, but got %s", expectedDistance.String(), calculatedDistance.String())
	}
}

// TestString verifies that the String function returns the correct hexadecimal representation
func TestString(t *testing.T) {
	// Full-length 20-byte KademliaID (40 hexadecimal characters)
	id := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFF0000000000000000")
	expectedString := "ffffffffffffffffffffffff0000000000000000" // Expected string in lowercase

	if id.String() != expectedString {
		t.Errorf("Expected string to be %s, but got %s", expectedString, id.String())
	}
}

func TestCheckTTLs(t *testing.T) {
	// Create a Kademlia instance with a real routing table and network
	me := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	routingTable := NewRoutingTable(me)
	storage := NewStorage()
	network := &Network{
		ID:        *me.ID,
		MessageCh: make(chan Message),
	}
	kademlia := NewKademlia(network, routingTable, []Task{}, storage)

	// Set up a task with WaitingForReturns contacts
	task := Task{
		CommandID:      1,
		CommandType:    "ping",
		TargetID:       NewRandomKademliaID(),
		ReplaceContact: NewContact(NewRandomKademliaID(), "127.0.0.1:8001"),
		ClosestContacts: []Contact{
			NewContact(NewRandomKademliaID(), "127.0.0.1:8002"),
			NewContact(NewRandomKademliaID(), "127.0.0.1:8003"),
		},
		WaitingForReturns: []WaitingContact{
			// First contact is older than TTL (it should time out)
			{SentTime: time.Now().Add(-time.Duration(TTL+1) * time.Millisecond), Contact: NewContact(NewRandomKademliaID(), "127.0.0.1:8004")},
			// Second contact is within TTL (it should not time out)

		},
	}

	// Add the task to the Kademlia node's task list
	kademlia.Tasks = append(kademlia.Tasks, task)
	if len(kademlia.Tasks) != 1 {
		t.Errorf("Expected 1 task in tasks, got %d", len(kademlia.Tasks))

	}
	// Call the checkTTLs method to trigger TTL check
	kademlia.checkTTLs()
	if len(kademlia.Tasks) != 0 {
		t.Errorf("Expected 0 task in tasks, got %d", len(kademlia.Tasks))

	}

}
