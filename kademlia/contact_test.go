package kademlia

import (
	"testing"
)

// TestNewContact checks that a new Contact is created correctly
func TestNewContact(t *testing.T) {
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	address := "127.0.0.1:8080"
	contact := NewContact(id, address)

	if contact.ID != id {
		t.Errorf("Expected contact ID to be %s, but got %s", id.String(), contact.ID.String())
	}

	if contact.Address != address {
		t.Errorf("Expected contact Address to be %s, but got %s", address, contact.Address)
	}
}

// TestCalcDistance checks that the CalcDistance method calculates the correct XOR distance
func TestCalcDistance(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("00000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	contact := NewContact(id1, "127.0.0.1:8080")

	contact.CalcDistance(id2)
	expectedDistance := id1.CalcDistance(id2)

	if !contact.distance.Equals(expectedDistance) {
		t.Errorf("Expected distance to be %s, but got %s", expectedDistance.String(), contact.distance.String())
	}
}

// TestLess checks if the Less method correctly compares two contacts based on their distance
func TestLess(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("00000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id3 := NewKademliaID("000000000000000000000000FFFFFFFFFFFFFFFF")
	contact1 := NewContact(id1, "127.0.0.1:8080")
	contact2 := NewContact(id2, "127.0.0.1:8081")
	contact3 := NewContact(id3, "127.0.0.1:8082")

	contact1.CalcDistance(id3)
	contact2.CalcDistance(id3)

	if !contact1.Less(&contact2) {
		t.Errorf("Expected contact1 to be less than contact2")
	}

	contact2.CalcDistance(id3)
	contact3.CalcDistance(id3)

	if contact3.Less(&contact2) {
		t.Errorf("Did not expect contact3 to be less than contact2")
	}
}

// TestString checks the string representation of a contact
func TestString(t *testing.T) {
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	contact := NewContact(id, "127.0.0.1:8080")
	expectedString := `contact("ffffffff00000000000000000000000000000000", "127.0.0.1:8080")`

	if contact.String() != expectedString {
		t.Errorf("Expected string to be %s, but got %s", expectedString, contact.String())
	}
}
