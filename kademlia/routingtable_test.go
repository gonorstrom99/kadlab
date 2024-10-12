package kademlia

import (
	"testing"
)

// TestGetMe tests the GetMe method of the RoutingTable
func TestGetMe(t *testing.T) {
	// Create a test contact (me)
	me := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	routingTable := NewRoutingTable(me)

	// Call GetMe and check if it returns the correct contact
	retrievedContact := routingTable.GetMe()

	if !retrievedContact.ID.Equals(me.ID) {
		t.Errorf("Expected self contact ID %s, got %s", me.ID.String(), retrievedContact.ID.String())
	}

	if retrievedContact.Address != me.Address {
		t.Errorf("Expected self contact address %s, got %s", me.Address, retrievedContact.Address)
	}
}

// TestAddContact tests the AddContact method of the RoutingTable
func TestAddContactRoutingTable(t *testing.T) {
	// Create the RoutingTable with a test node as 'me'
	me := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	routingTable := NewRoutingTable(me)

	// Create a new contact to add to the routing table
	contact1 := NewContact(NewRandomKademliaID(), "127.0.0.1:8001")
	contact2 := NewContact(NewRandomKademliaID(), "127.0.0.1:8002")

	// Add the first contact
	routingTable.AddContact(contact1)

	// Verify that the contact is added to the correct bucket
	bucketIndex1 := routingTable.getBucketIndex(contact1.ID)
	bucket1 := routingTable.buckets[bucketIndex1]

	// Check if the contact is in the bucket
	if !bucket1.IsContactInBucket(&contact1) {
		t.Errorf("Contact %s not found in bucket after adding", contact1.Address)
	}

	// Add another contact to the routing table
	routingTable.AddContact(contact2)

	// Verify that the second contact is added to the correct bucket
	bucketIndex2 := routingTable.getBucketIndex(contact2.ID)
	bucket2 := routingTable.buckets[bucketIndex2]

	// Check if the second contact is in the bucket
	if !bucket2.IsContactInBucket(&contact2) {
		t.Errorf("Contact %s not found in bucket after adding", contact2.Address)
	}

	// Add the first contact again and ensure it moves to the front of the bucket
	routingTable.AddContact(contact1)

	// Verify that the first contact is now at the front of the bucket
	frontContact := bucket1.list.Front().Value.(Contact)
	if !frontContact.ID.Equals(contact1.ID) {
		t.Errorf("Contact %s was not moved to the front of the bucket", contact1.Address)
	}

	// Check that the second contact is still in its correct bucket
	if !bucket2.IsContactInBucket(&contact2) {
		t.Errorf("Contact %s was removed unexpectedly", contact2.Address)
	}
}

// TestFindClosestContacts tests the FindClosestContacts method of the RoutingTable
func TestFindClosestContacts(t *testing.T) {
	// Create a test target ID to search for
	targetID := NewRandomKademliaID()

	// Create the RoutingTable with a test node as 'me'
	me := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	routingTable := NewRoutingTable(me)

	// Add some contacts to the routing table
	contacts := []Contact{
		NewContact(NewRandomKademliaID(), "127.0.0.1:8001"),
		NewContact(NewRandomKademliaID(), "127.0.0.1:8002"),
		NewContact(NewRandomKademliaID(), "127.0.0.1:8003"),
		NewContact(NewRandomKademliaID(), "127.0.0.1:8004"),
		NewContact(NewRandomKademliaID(), "127.0.0.1:8005"),
	}

	for _, contact := range contacts {
		routingTable.AddContact(contact)
	}

	// Define the number of closest contacts to find
	count := 3

	// Call FindClosestContacts to find the closest contacts
	foundContacts := routingTable.FindClosestContacts(targetID, count)

	// Check if the number of returned contacts is correct
	if len(foundContacts) != count {
		t.Errorf("Expected %d closest contacts, but got %d", count, len(foundContacts))
	}

	// Optionally, check that the contacts returned are sorted by XOR distance (depending on your Sort method)
	// Assuming GetContactAndCalcDistance sorts them by distance, you can also verify ordering correctness
	for i := 0; i < len(foundContacts)-1; i++ {
		dist1 := foundContacts[i].ID.CalcDistance(targetID)
		dist2 := foundContacts[i+1].ID.CalcDistance(targetID)

		if !dist1.Less(dist2) && !dist1.Equals(dist2) {
			t.Errorf("Contacts are not sorted by distance: contact %d is farther than contact %d", i, i+1)
		}
	}
}

// TestGetBucketIndex tests the getBucketIndex method of the RoutingTable
func TestGetBucketIndex(t *testing.T) {
	// Create a test contact (me)
	me := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	routingTable := NewRoutingTable(me)

	// Create another contact and calculate the bucket index
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8001")

	bucketIndex := routingTable.getBucketIndex(contact.ID)

	// Verify the bucket index is within the valid range (0 to IDLength * 8 - 1)
	if bucketIndex < 0 || bucketIndex >= IDLength*8 {
		t.Errorf("Bucket index %d is out of range for contact %s", bucketIndex, contact.ID.String())
	}
}

// TestGetBucket tests the getBucket method of the RoutingTable
func TestGetBucket(t *testing.T) {
	// Create a test contact (me)
	me := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	routingTable := NewRoutingTable(me)

	// Create another contact and get its bucket index
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8001")
	bucketIndex := routingTable.getBucketIndex(contact.ID)

	// Retrieve the bucket and check if it's nil
	bucket := routingTable.buckets[bucketIndex]
	if bucket == nil {
		t.Errorf("Bucket for index %d is nil", bucketIndex)
	}

	// Check if the bucket is valid by adding the contact and verifying it's stored
	bucket.AddContact(contact)

	// Ensure the contact is now in the bucket
	if !bucket.IsContactInBucket(&contact) {
		t.Errorf("Contact %s not found in bucket after adding", contact.ID.String())
	}
}
