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

// TestAppend checks if contacts are correctly appended to the ContactCandidates
func TestAppend(t *testing.T) {
	candidates := ContactCandidates{}
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8080")
	contact2 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "127.0.0.1:8081")

	candidates.Append([]Contact{contact1, contact2})

	if len(candidates.contacts) != 2 {
		t.Errorf("Expected 2 contacts, but got %d", len(candidates.contacts))
	}

	if candidates.contacts[0] != contact1 || candidates.contacts[1] != contact2 {
		t.Errorf("Contacts were not appended correctly")
	}
}

// TestGetContacts checks if the correct number of contacts is returned by GetContacts
func TestGetContacts(t *testing.T) {
	candidates := ContactCandidates{}
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8080")
	contact2 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "127.0.0.1:8081")
	contact3 := NewContact(NewKademliaID("BBBBBBBB00000000000000000000000000000000"), "127.0.0.1:8082")

	candidates.Append([]Contact{contact1, contact2, contact3})

	contacts := candidates.GetContacts(2)
	if len(contacts) != 2 {
		t.Errorf("Expected 2 contacts, but got %d", len(contacts))
	}
	if contacts[0] != contact1 || contacts[1] != contact2 {
		t.Errorf("Returned contacts are not correct")
	}
}

// TestSort checks if contacts are sorted correctly by distance
func TestSort(t *testing.T) {
	id := NewKademliaID("0000000000000000000000000000000000000000")
	candidates := ContactCandidates{}
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8080")
	contact2 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "127.0.0.1:8081")

	contact1.CalcDistance(id)
	contact2.CalcDistance(id)

	candidates.Append([]Contact{contact2, contact1}) // Reverse order

	candidates.Sort()

	if !candidates.contacts[0].Less(&candidates.contacts[1]) {
		t.Errorf("Expected contact1 to be less than contact2 after sorting")
	}
}

// TestLen checks if the length of the ContactCandidates is returned correctly
func TestLen(t *testing.T) {
	candidates := ContactCandidates{}
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8080")
	contact2 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "127.0.0.1:8081")

	candidates.Append([]Contact{contact1, contact2})

	if candidates.Len() != 2 {
		t.Errorf("Expected length to be 2, but got %d", candidates.Len())
	}
}

// TestSwap checks if two contacts are swapped correctly
func TestSwap(t *testing.T) {
	candidates := ContactCandidates{}
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8080")
	contact2 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "127.0.0.1:8081")

	candidates.Append([]Contact{contact1, contact2})

	candidates.Swap(0, 1)

	if candidates.contacts[0] != contact2 || candidates.contacts[1] != contact1 {
		t.Errorf("Contacts were not swapped correctly")
	}
}

// TestLess checks if the Less function compares contacts correctly by distance
func TestLessInContact(t *testing.T) {
	id := NewKademliaID("0000000000000000000000000000000000000000")
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8080")
	contact2 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "127.0.0.1:8081")

	contact1.CalcDistance(id)
	contact2.CalcDistance(id)

	candidates := ContactCandidates{}
	candidates.Append([]Contact{contact1, contact2})

	if !candidates.Less(1, 0) {
		t.Errorf("Expected contact2 to be less than contact1")
	}
}
