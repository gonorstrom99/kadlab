package kademlia

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

// Task represents an ongoing task for the Kademlia node.
type Task struct {
	CommandType       string      // The type of task (e.g., "lookUpContact", "findValue")
	CommandID         int         //
	TargetID          *KademliaID // The ID or we're looking for (if applicable)
	StartTime         time.Time   // When the job was started (for timeouts)
	StartContact      *Contact    // The contact that initiated the job
	ClosestContacts   []Contact   // for storing nodes that are closest
	ContactedNodes    []Contact
	waitingForReturns []Contact
}

// AddTaskFromMessage takes a Message struct and a contact, parses the information,
// and adds the task to the Kademlia node's task list.
func (kademlia *Kademlia) AddTaskFromMessage(msg Message, contact *Contact) {
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
		CommandType:     commandType,
		CommandID:       commandID,
		TargetID:        targetID,
		StartTime:       time.Now(),
		Contact:         contact,
		ClosestContacts: make([]Contact, 0), // Initialize as an empty slice
	}

	// Log the creation of the task for debugging purposes
	log.Printf("Task added: CommandType=%s, CommandID=%d, TargetID=%s, Contact=%s", commandType, commandID, targetID.String(), contact.Address)

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

// CheckContactStatus pings a contact and returns true if it receives a pong response.
func (kademlia *Kademlia) CheckContactStatus(contact *Contact) bool {
	id := kademlia.RoutingTable.me.ID.String()
	commandID := NewCommandID() // Generate a new command ID
	messageString := fmt.Sprintf("ping:%s:%d", id, commandID)

	// Create a new task for this ping
	task := Task{
		CommandType: "ping",
		CommandID:   commandID,
		TargetID:    nil,
		StartTime:   time.Now(),
		Contact:     contact,
		IsCompleted: false,
	}

	// Add the task to the node's task list
	kademlia.Tasks = append(kademlia.Tasks, task)
	kademlia.Network.SendPingMessage(contact, messageString)

	// Wait for pong or timeout
	timeOut := time.After(pongTimer * time.Second)
	waitTime := time.Second
	var pongReceived bool = false

	for {
		select {
		case <-timeOut:
			fmt.Println("Waited five seconds, contact presumed dead")
			kademlia.MarkTaskAsCompleted(commandID, false)
			return pongReceived

		case <-time.After(waitTime):
			task, _ := kademlia.FindTaskByCommandID(commandID)
			if task.IsCompleted {
				fmt.Println("Pong received from correct contact")
				return true
			}
			fmt.Println("Still waiting for pong...")
		}
	}
}

// MarkTaskAsCompleted updates the status of a task by commandID
func (kademlia *Kademlia) MarkTaskAsCompleted(commandID int, success bool) {
	for i, task := range kademlia.Tasks {
		if task.CommandID == commandID {
			kademlia.Tasks[i].IsCompleted = success
			log.Printf("Task %d marked as completed: %t", commandID, success)
			return
		}
	}
}

// RemoveCompletedTasks removes tasks that are marked as completed.
func (kademlia *Kademlia) RemoveCompletedTasks() {
	var activeTasks []Task
	for _, task := range kademlia.Tasks {
		if !task.IsCompleted {
			activeTasks = append(activeTasks, task)
		}
	}
	kademlia.Tasks = activeTasks
}

// HasTaskTimedOut checks whether a task has timed out based on a given duration.
func (task *Task) HasTaskTimedOut(timeout time.Duration) bool {
	return time.Since(task.StartTime) > timeout
}
