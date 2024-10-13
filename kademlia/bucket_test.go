package kademlia

import (
	"strconv"
	"testing"
)

// TestNewBucket checks if a new bucket is properly initialized
func TestNewBucket(t *testing.T) {
	bucket := newBucket()
	if bucket.list == nil {
		t.Error("Expected new bucket to have a non-nil list")
	}
	if bucket.Len() != 0 {
		t.Error("Expected new bucket to be empty")
	}
}

// TestAddContact checks if contacts are added properly to the bucket
func TestAddContact(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")

	// Add the contact and verify it is in the bucket
	bucket.AddContact(contact)
	if bucket.Len() != 1 {
		t.Errorf("Expected bucket size to be 1, got %d", bucket.Len())
	}
	if !bucket.IsContactInBucket(&contact) {
		t.Error("Expected contact to be in the bucket but it wasn't")
	}

	// Add the same contact again and check that the size doesn't change
	bucket.AddContact(contact)
	if bucket.Len() != 1 {
		t.Errorf("Expected bucket size to remain 1, got %d", bucket.Len())
	}
}

// TestAddMultipleContacts checks if multiple contacts can be added to the bucket
func TestAddMultipleContacts(t *testing.T) {
	bucket := newBucket()

	for i := 0; i < bucketSize; i++ {
		contact := NewContact(NewRandomKademliaID(), "127.0.0.1:"+strconv.Itoa(8000+i))
		bucket.AddContact(contact)
		if bucket.Len() != i+1 {
			t.Errorf("Expected bucket size to be %d, got %d", i+1, bucket.Len())
		}
	}

	// Check that the bucket does not exceed its capacity
	if bucket.Len() > bucketSize {
		t.Errorf("Bucket size exceeded bucketSize: got %d", bucket.Len())
	}
}

// TestGetContactAndCalcDistance checks if the contacts are returned with calculated distance
func TestGetContactAndCalcDistance(t *testing.T) {
	bucket := newBucket()
	targetID := NewRandomKademliaID()

	for i := 0; i < 5; i++ {
		contact := NewContact(NewRandomKademliaID(), "127.0.0.1:"+strconv.Itoa(8000+i))
		bucket.AddContact(contact)
	}

	contacts := bucket.GetContactAndCalcDistance(targetID)

	if len(contacts) != bucket.Len() {
		t.Errorf("Expected %d contacts, but got %d", bucket.Len(), len(contacts))
	}

	// Verify that the distance has been calculated for each contact
	for _, contact := range contacts {
		if contact.distance == nil {
			t.Error("Expected contact to have a calculated distance but it didn't")
		}
	}
}

// TestIsContactInBucket checks if the function correctly identifies if a contact is in the bucket
func TestIsContactInBucket(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	bucket.AddContact(contact)

	if !bucket.IsContactInBucket(&contact) {
		t.Error("Expected contact to be in the bucket but it wasn't")
	}

	// Test for a contact not in the bucket
	otherContact := NewContact(NewRandomKademliaID(), "127.0.0.1:8001")
	if bucket.IsContactInBucket(&otherContact) {
		t.Error("Did not expect the contact to be in the bucket but it was")
	}
}

// TestBucketLen checks if the Len function returns the correct size
func TestBucketLen(t *testing.T) {
	bucket := newBucket()

	if bucket.Len() != 0 {
		t.Errorf("Expected bucket length to be 0, got %d", bucket.Len())
	}

	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	bucket.AddContact(contact)

	if bucket.Len() != 1 {
		t.Errorf("Expected bucket length to be 1, got %d", bucket.Len())
	}
}
