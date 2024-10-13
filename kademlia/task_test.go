package kademlia

import (
	"strconv"
	"testing"
	"time"
)

// TestUpdateTaskFromMessage checks if a task is updated correctly from a Message
func TestUpdateTaskFromMessage(t *testing.T) {
	kademlia := &Kademlia{}
	msg := Message{
		Command:     "StoreValue",
		CommandID:   "12345",
		CommandInfo: "FFFFFFFFFF0000000000FFFFFFFFFF0000000000", // Sample target ID
	}
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8080")

	kademlia.UpdateTaskFromMessage(msg, &contact)

	// Check if the task was added
	if len(kademlia.Tasks) == 0 {
		t.Fatalf("Expected a task to be added, but none was")
	}

	task := kademlia.Tasks[0]

	// Check if task details match the message
	if task.CommandType != "StoreValue" {
		t.Errorf("Expected CommandType to be 'StoreValue', but got %s", task.CommandType)
	}

	expectedID, _ := strconv.Atoi("12345")
	if task.CommandID != expectedID {
		t.Errorf("Expected CommandID to be %d, but got %d", expectedID, task.CommandID)
	}

	// Normalize both values to lowercase to avoid case sensitivity issues
	expectedTargetID := "ffffffffff0000000000ffffffffff0000000000" // lowercase expected value
	actualTargetID := task.TargetID.String()

	if actualTargetID != expectedTargetID {
		t.Errorf("Expected TargetID to be %s, but got %s", expectedTargetID, actualTargetID)
	}
}

// TestContactIsContacted verifies if the task correctly identifies already contacted nodes
func TestContactIsContacted(t *testing.T) {
	task := Task{
		ContactedNodes: []Contact{
			NewContact(NewRandomKademliaID(), "127.0.0.1:8080"),
			NewContact(NewRandomKademliaID(), "127.0.0.1:8081"),
		},
	}

	contact := task.ContactedNodes[0]

	if !task.ContactIsContacted(contact) {
		t.Errorf("Expected contact to be identified as contacted")
	}

	nonExistentContact := NewContact(NewRandomKademliaID(), "127.0.0.1:8082")

	if task.ContactIsContacted(nonExistentContact) {
		t.Errorf("Did not expect non-existent contact to be identified as contacted")
	}
}

// TestFindFirstNotContactedNodeIndex checks the method for finding the first non-contacted node
func TestFindFirstNotContactedNodeIndex(t *testing.T) {
	contact1 := NewContact(NewRandomKademliaID(), "127.0.0.1:8080")
	contact2 := NewContact(NewRandomKademliaID(), "127.0.0.1:8081")
	contact3 := NewContact(NewRandomKademliaID(), "127.0.0.1:8082")

	task := Task{
		ClosestContacts: []Contact{
			contact1,
			contact2,
			contact3,
		},
		ContactedNodes: []Contact{
			contact1, // Simulating that contact1 was already contacted
		},
	}

	index := task.FindFirstNotContactedNodeIndex()

	// Expect the first uncontacted node to be at index 1 (contact2)
	if index != 1 {
		t.Errorf("Expected first uncontacted node index to be 1, but got %d", index)
	}

	// Mark all nodes as contacted and check that index is -1
	task.ContactedNodes = append(task.ContactedNodes, contact2)
	task.ContactedNodes = append(task.ContactedNodes, contact3)

	index = task.FindFirstNotContactedNodeIndex()

	if index != -1 {
		t.Errorf("Expected -1 (no uncontacted nodes), but got %d", index)
	}
}

// TestAreClosestContactsContacted checks if all closest contacts have been contacted
func TestAreClosestContactsContacted(t *testing.T) {
	contact1 := NewContact(NewRandomKademliaID(), "127.0.0.1:8080")
	contact2 := NewContact(NewRandomKademliaID(), "127.0.0.1:8081")

	task := Task{
		ClosestContacts: []Contact{
			contact1,
			contact2,
		},
		ContactedNodes: []Contact{
			contact1, // Simulating that contact1 was already contacted
		},
	}

	// Initially, not all closest contacts have been contacted
	if task.AreClosestContactsContacted() {
		t.Errorf("Expected AreClosestContactsContacted to return false, but it returned true")
	}

	// Mark contact2 as contacted
	task.ContactedNodes = append(task.ContactedNodes, contact2)

	// Now, all closest contacts should be contacted
	if !task.AreClosestContactsContacted() {
		t.Errorf("Expected AreClosestContactsContacted to return true, but it returned false")
	}
}

// TestCreateTask checks if a task is created with the correct fields
func TestCreateTask(t *testing.T) {
	kademlia := &Kademlia{}
	commandType := "StoreValue"
	commandID := 12345
	targetID := NewRandomKademliaID()

	task := kademlia.CreateTask(commandType, commandID, targetID)

	// Check if the task fields are correctly initialized
	if task.CommandType != commandType {
		t.Errorf("Expected CommandType to be %s, but got %s", commandType, task.CommandType)
	}

	if task.CommandID != commandID {
		t.Errorf("Expected CommandID to be %d, but got %d", commandID, task.CommandID)
	}

	if task.TargetID.String() != targetID.String() {
		t.Errorf("Expected TargetID to be %s, but got %s", targetID.String(), task.TargetID.String())
	}

	if len(task.ClosestContacts) != 0 {
		t.Errorf("Expected ClosestContacts to be an empty slice, but got %d elements", len(task.ClosestContacts))
	}

	if len(task.ContactedNodes) != 0 {
		t.Errorf("Expected ContactedNodes to be an empty slice, but got %d elements", len(task.ContactedNodes))
	}

	if len(task.WaitingForReturns) != 0 {
		t.Errorf("Expected WaitingForReturns to be an empty slice, but got %d elements", len(task.WaitingForReturns))
	}
}

// TestRemoveContactFromWaitingForReturnsByTask checks if a contact is correctly removed from the WaitingForReturns list
func TestRemoveContactFromWaitingForReturnsByTask(t *testing.T) {
	task := Task{
		WaitingForReturns: []WaitingContact{
			{Contact: NewContact(NewRandomKademliaID(), "127.0.0.1:8080")},
			{Contact: NewContact(NewRandomKademliaID(), "127.0.0.1:8081")},
		},
	}

	contactID := task.WaitingForReturns[0].Contact.ID

	// Remove the first contact and verify it's removed
	task.RemoveContactFromWaitingForReturnsByTask(*contactID)

	// Check if the first contact is removed from the WaitingForReturns list
	if len(task.WaitingForReturns) != 1 {
		t.Errorf("Expected WaitingForReturns length to be 1, but got %d", len(task.WaitingForReturns))
	}

	if task.WaitingForReturns[0].Contact.ID.Equals(contactID) {
		t.Errorf("Expected contact %s to be removed, but it still exists in WaitingForReturns", contactID.String())
	}
}

// TestRemoveContactFromWaitingForReturns checks if a contact is removed from WaitingForReturns using commandID
func TestRemoveContactFromWaitingForReturns(t *testing.T) {
	kademlia := &Kademlia{
		Tasks: []Task{
			{
				CommandID: 12345,
				WaitingForReturns: []WaitingContact{
					{Contact: NewContact(NewRandomKademliaID(), "127.0.0.1:8080")},
					{Contact: NewContact(NewRandomKademliaID(), "127.0.0.1:8081")},
				},
			},
		},
	}

	contactID := kademlia.Tasks[0].WaitingForReturns[0].Contact.ID

	// Remove the first contact using commandID and verify it's removed
	kademlia.RemoveContactFromWaitingForReturns(12345, *contactID)

	task := kademlia.Tasks[0]

	// Check if the first contact is removed from the WaitingForReturns list
	if len(task.WaitingForReturns) != 1 {
		t.Errorf("Expected WaitingForReturns length to be 1, but got %d", len(task.WaitingForReturns))
	}

	if task.WaitingForReturns[0].Contact.ID.Equals(contactID) {
		t.Errorf("Expected contact %s to be removed, but it still exists in WaitingForReturns", contactID.String())
	}
}

// TestIsContactInClosestContacts checks if the contact is in the ClosestContacts list
func TestIsContactInClosestContacts(t *testing.T) {
	task := Task{
		ClosestContacts: []Contact{
			NewContact(NewRandomKademliaID(), "127.0.0.1:8080"),
			NewContact(NewRandomKademliaID(), "127.0.0.1:8081"),
		},
	}

	contact := task.ClosestContacts[0]

	// Check if the contact is found in ClosestContacts
	if !task.IsContactInClosestContacts(contact) {
		t.Errorf("Expected contact %s to be in ClosestContacts, but it was not found", contact.ID.String())
	}

	nonExistentContact := NewContact(NewRandomKademliaID(), "127.0.0.1:8082")
	// Check if a non-existent contact is correctly reported as not found
	if task.IsContactInClosestContacts(nonExistentContact) {
		t.Errorf("Did not expect non-existent contact to be in ClosestContacts")
	}
}

// TestRemoveContactFromTask checks if the contact is removed from the WaitingForReturns list
func TestRemoveContactFromTask(t *testing.T) {
	kademlia := &Kademlia{
		Tasks: []Task{
			{
				CommandID: 12345,
				WaitingForReturns: []WaitingContact{
					{Contact: NewContact(NewRandomKademliaID(), "127.0.0.1:8080")},
					{Contact: NewContact(NewRandomKademliaID(), "127.0.0.1:8081")},
				},
			},
		},
	}

	contact := kademlia.Tasks[0].WaitingForReturns[0].Contact

	// Remove the first contact
	kademlia.RemoveContactFromTask(12345, contact)

	// Verify if the contact has been removed from the task's WaitingForReturns list
	task := kademlia.Tasks[0]
	if len(task.WaitingForReturns) != 1 {
		t.Errorf("Expected WaitingForReturns length to be 1, but got %d", len(task.WaitingForReturns))
	}

	if task.WaitingForReturns[0].Contact.ID.Equals(contact.ID) {
		t.Errorf("Expected contact %s to be removed, but it still exists in WaitingForReturns", contact.ID.String())
	}
}

// TestMarkTaskAsCompleted checks if a task is marked as completed and removed from the task list
func TestMarkTaskAsCompleted(t *testing.T) {
	kademlia := &Kademlia{
		Tasks: []Task{
			{
				CommandID:         12345,
				CommandType:       "StoreValue",
				WaitingForReturns: []WaitingContact{},
			},
		},
	}

	// Mark the task as completed
	kademlia.MarkTaskAsCompleted(12345)

	// Check if the task has been removed from the task list
	if len(kademlia.Tasks) != 0 {
		t.Errorf("Expected task to be removed, but it still exists in the task list")
	}
}

// TestSortContactsByDistance checks if the contacts in ClosestContacts are sorted by XOR distance
func TestSortContactsByDistance(t *testing.T) {
	targetID := NewRandomKademliaID()

	task := Task{
		ClosestContacts: []Contact{
			NewContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "127.0.0.1:8080"),
			NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1:8081"),
			NewContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFE"), "127.0.0.1:8082"),
		},
		TargetID: targetID,
	}

	// Sort the contacts by distance
	task.SortContactsByDistance()

	// Check if contacts are sorted correctly by XOR distance
	distanceToFirst := task.ClosestContacts[0].ID.CalcDistance(targetID)
	distanceToSecond := task.ClosestContacts[1].ID.CalcDistance(targetID)

	if !distanceToFirst.Less(distanceToSecond) {
		t.Errorf("Expected the first contact to be closer to the target than the second one")
	}
}

// TestHasTaskTimedOut checks if the task is considered timed out after a given duration
func TestHasTaskTimedOut(t *testing.T) {
	task := Task{
		StartTime: time.Now().Add(-time.Minute), // Started one minute ago
	}

	timeout := 30 * time.Second // Set timeout to 30 seconds

	if !task.HasTaskTimedOut(timeout) {
		t.Errorf("Expected task to be timed out, but it was not")
	}

	task.StartTime = time.Now() // Reset start time to now

	if task.HasTaskTimedOut(timeout) {
		t.Errorf("Did not expect task to be timed out, but it was")
	}
}
