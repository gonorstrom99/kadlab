package kademlia

import (
	"encoding/hex"
	"testing"
)

// TestNewKademliaID verifies that a new KademliaID is created correctly from a string
func TestNewKademliaID(t *testing.T) {
	data := "FFFFFFFF0000000000"
	kademliaID := NewKademliaID(data)

	expectedBytes, _ := hex.DecodeString(data)
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
	id1 := NewKademliaID("00000000000000000000")
	id2 := NewKademliaID("FFFFFFFFFFFFFFFFFFFF")

	if !id1.Less(id2) {
		t.Errorf("Expected id1 to be less than id2")
	}

	if id2.Less(id1) {
		t.Errorf("Did not expect id2 to be less than id1")
	}
}

// TestEquals checks if two KademliaIDs are correctly identified as equal
func TestEquals(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF0000000000")
	id2 := NewKademliaID("FFFFFFFF0000000000")
	id3 := NewKademliaID("00000000FFFFFFFFFF")

	if !id1.Equals(id2) {
		t.Errorf("Expected id1 to equal id2")
	}

	if id1.Equals(id3) {
		t.Errorf("Did not expect id1 to equal id3")
	}
}

// TestCalcDistance checks the CalcDistance function for XOR distance calculation
func TestCalcDistance(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF0000000000")
	id2 := NewKademliaID("00000000FFFFFFFFFF")

	expectedDistance := NewKademliaID("FFFFFFFFFFFFFFFFFF")

	calculatedDistance := id1.CalcDistance(id2)

	if !calculatedDistance.Equals(expectedDistance) {
		t.Errorf("Expected distance to be %s, but got %s", expectedDistance.String(), calculatedDistance.String())
	}
}

// TestString verifies that the String function returns the correct hexadecimal representation
func TestString(t *testing.T) {
	id := NewKademliaID("FFFFFFFF0000000000")
	expectedString := "ffffffff0000000000"

	if id.String() != expectedString {
		t.Errorf("Expected string to be %s, but got %s", expectedString, id.String())
	}
}
