package kademlia

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"
)

// WaitingContact represents a contact we're waiting for a response from, along with the time the message was sent.
type WaitingContact struct {
	SentTime time.Time // The time when the message was sent
	Contact  Contact   // The contact we are waiting for a response from
}

// Task represents an ongoing task for the Kademlia node.
type Task struct {
	CommandType       string           // The type of task (e.g., "lookUpContact", "findValue")
	CommandID         int              // The unique identifier for the command
	TargetID          *KademliaID      // The ID we're looking for (if applicable)
	StartTime         time.Time        // When the job was started (for timeouts)
	ClosestContacts   []Contact        // for storing nodes that are closest
	ContactedNodes    []Contact        // All the contacted nodes that we shouldn't contact again
	WaitingForReturns []WaitingContact // Alpha number of nodes that we're waiting for returns from
}

// UpdateTaskFromMessage takes a Message struct and a contact, parses the information,
// and updates the task to the Kademlia node's task list.
func (kademlia *Kademlia) UpdateTaskFromMessage(msg Message, contact *Contact) {
	// Parse the message fields from the Message struct
	commandType := msg.Command
	commandID, err := strconv.Atoi(msg.CommandID)
	if err != nil {
		log.Printf("Invalid command ID: %s", msg.CommandID)
		return
	}
	targetID := NewKademliaID(msg.CommandInfo) // Assuming CommandInfo is the targetID in this case

	// Create a new task and set its fields
	task := Task{
		CommandType:       commandType,
		CommandID:         commandID,
		TargetID:          targetID,
		StartTime:         time.Now(),
		ClosestContacts:   make([]Contact, 0),        // Initialize as an empty slice
		ContactedNodes:    make([]Contact, 0),        // Initialize as an empty slice
		WaitingForReturns: make([]WaitingContact, 0), // Initialize as an empty slice
	}

	// Log the creation of the task for debugging purposes
	log.Printf("Task added: CommandType=%s, CommandID=%d, TargetID=%s", commandType, commandID, targetID.String())

	// Add the task to the node's task list (assuming you have a task list)
	kademlia.Tasks = append(kademlia.Tasks, task)
}

// FindTaskByCommandID takes a Message and looks for a matching Task with the same CommandID in the task list
func (kademlia *Kademlia) FindTaskByCommandID(commandID int) (*Task, error) {
	for _, task := range kademlia.Tasks {
		if task.CommandID == commandID {
			return &task, nil
		}
	}
	return nil, fmt.Errorf("task with CommandID %d not found", commandID)
}

// RemoveContactFromTask removes a contact from the WaitingForReturns list when they respond
func (kademlia *Kademlia) RemoveContactFromTask(commandID int, contact Contact) {
	task, err := kademlia.FindTaskByCommandID(commandID)
	if err != nil {
		log.Printf("Task with CommandID %d not found, cannot remove contact", commandID)
		return
	}

	// Remove contact from WaitingForReturns
	for i, waitingContact := range task.WaitingForReturns {
		if waitingContact.Contact.ID.Equals(contact.ID) {
			// Remove the contact from the list
			task.WaitingForReturns = append(task.WaitingForReturns[:i], task.WaitingForReturns[i+1:]...)
			log.Printf("Contact %s removed from WaitingForReturns in task %d", contact.ID.String(), commandID)
			break
		}
	}
}

// MarkTaskAsCompleted updates the task when all contacts have responded or it is considered done
func (kademlia *Kademlia) MarkTaskAsCompleted(commandID int) {
	task, err := kademlia.FindTaskByCommandID(commandID)
	if err != nil {
		log.Printf("Task with CommandID %d not found", commandID)
		return
	}

	// Task is completed when no more waiting for returns
	if len(task.WaitingForReturns) == 0 {
		log.Printf("Task %d is completed", commandID)
		// Optionally, you can remove the task from the task list here
		kademlia.RemoveTask(commandID)
	}
}

// SortContactsByDistance sorts the contacts in the task's ClosestContacts based on their XOR distance to the TargetID.
func (task *Task) SortContactsByDistance() {
	if task.TargetID == nil {
		log.Printf("Task %d has no TargetID, unable to sort contacts", task.CommandID)
		return
	}

	// Define a custom sort function that compares XOR distances
	sort.Slice(task.ClosestContacts, func(i, j int) bool {
		distanceToI := task.ClosestContacts[i].ID.CalcDistance(task.TargetID)
		distanceToJ := task.ClosestContacts[j].ID.CalcDistance(task.TargetID)

		// Compare distances: return true if the distance to I is less than the distance to J
		return distanceToI.Less(distanceToJ)
	})

	log.Printf("Contacts in Task %d sorted by distance to TargetID", task.CommandID)
}

// RemoveTask removes a completed task by commandID
func (kademlia *Kademlia) RemoveTask(commandID int) {
	var remainingTasks []Task
	for _, task := range kademlia.Tasks {
		if task.CommandID != commandID {
			remainingTasks = append(remainingTasks, task)
		}
	}
	kademlia.Tasks = remainingTasks
}

// HasTaskTimedOut checks whether a task has timed out based on a given duration.
func (task *Task) HasTaskTimedOut(timeout time.Duration) bool {
	return time.Since(task.StartTime) > timeout
}
