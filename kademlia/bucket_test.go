package kademlia

import (
	"testing"
)

// TestAddContact tests adding contacts to the bucket
func TestAddContact(t *testing.T) {
	bucket := newBucket()

	// Create some contacts
	contact1 := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	contact2 := NewContact(NewRandomKademliaID(), "127.0.0.1:8001")

	// Add contacts to the bucket
	bucket.AddContact(contact1)
	bucket.AddContact(contact2)

	// Check if the contacts were added
	if !bucket.IsContactInBucket(&contact1) {
		t.Errorf("Expected contact1 to be in the bucket")
	}
	if !bucket.IsContactInBucket(&contact2) {
		t.Errorf("Expected contact2 to be in the bucket")
	}

	// Check if the length of the bucket is 2go test -p
	if bucket.Len() != 2 {
		t.Errorf("Expected bucket length to be 2, but got %d", bucket.Len())
	}
}

// Test that adding an existing contact moves it to the front of the bucket
func TestAddExistingContactMovesToFront(t *testing.T) {
	bucket := newBucket()

	// Create contacts
	contact1 := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	contact2 := NewContact(NewRandomKademliaID(), "127.0.0.1:8001")

	// Add both contacts to the bucket
	bucket.AddContact(contact1)
	bucket.AddContact(contact2)

	// Check that contact2 is at the front (as it was the last added)
	if bucket.list.Front().Value.(Contact).ID != contact2.ID {
		t.Errorf("Expected contact2 to be at the front")
	}

	// Re-add contact1 (should move it to the front)
	bucket.AddContact(contact1)

	// Check that contact1 is now at the front
	if bucket.list.Front().Value.(Contact).ID != contact1.ID {
		t.Errorf("Expected contact1 to be moved to the front")
	}
}

// Test that the bucket returns the correct length
func TestBucketLen(t *testing.T) {
	bucket := newBucket()

	// Check that the bucket is initially empty
	if bucket.Len() != 0 {
		t.Errorf("Expected bucket length to be 0, but got %d", bucket.Len())
	}

	// Add contacts to the bucket
	contact1 := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	contact2 := NewContact(NewRandomKademliaID(), "127.0.0.1:8001")
	bucket.AddContact(contact1)
	bucket.AddContact(contact2)

	// Check the length after adding contacts
	expectedLen := 2
	if bucket.Len() != expectedLen {
		t.Errorf("Expected bucket length to be %d, but got %d", expectedLen, bucket.Len())
	}
}

// Test GetContactAndCalcDistance to ensure distance is calculated correctly
func TestGetContactAndCalcDistance(t *testing.T) {
	bucket := newBucket()

	// Create a contact and a target KademliaID
	target := NewRandomKademliaID()
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")

	// Add contact to the bucket
	bucket.AddContact(contact)

	// Get contacts with calculated distances
	contacts := bucket.GetContactAndCalcDistance(target)

	// There should be one contact returned
	if len(contacts) != 1 {
		t.Errorf("Expected 1 contact, but got %d", len(contacts))
	}

	// The contact's distance should be calculated
	if contacts[0].distance == nil {
		t.Errorf("Expected distance to be calculated, but it was nil")
	}
}

// Test IsContactInBucket checks if a contact is correctly identified in the bucket
func TestIsContactInBucket(t *testing.T) {
	bucket := newBucket()

	// Create contacts
	contact1 := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	contact2 := NewContact(NewRandomKademliaID(), "127.0.0.1:8001")

	// Add contact1 to the bucket
	bucket.AddContact(contact1)

	// Check if contact1 is in the bucket
	if !bucket.IsContactInBucket(&contact1) {
		t.Errorf("Expected contact1 to be in the bucket")
	}

	// Check if contact2 is in the bucket (it should not be)
	if bucket.IsContactInBucket(&contact2) {
		t.Errorf("Did not expect contact2 to be in the bucket")
	}
}
