package kademlia

import (
	"log"
	"net"
	"strings"
)

// Network is a node in the network
type Network struct {
	Conn      *net.UDPConn
	MessageCh chan Message
}

// Message is a simple struct to hold a message and its sender's address.
type Message struct {
	Content string
	Address string
	ID      string
}

// Listen listens on a UDP address and stores incoming messages in the channel.
func (network *Network) Listen(contact Contact) error {
	// Create a UDPAddr based on the Network's IP and Port
	udpAddr, err := net.ResolveUDPAddr("udp", contact.Address)

	// Start listening on the provided address
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("Error starting UDP listener: %v\n", err)
		return err
	}
	network.Conn = conn // Store the connection for sending messages
	defer conn.Close()

	// Loop to continuously listen for incoming messages
	for {
		buffer := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP connection: %v", err)
			continue
		}

		// Convert the buffer into a string message
		message := string(buffer[:n])
		contactAddress := remoteAddr.String()

		// Split the message to extract the ID and the actual content
		parts := strings.SplitN(message, ":", 2)
		if len(parts) != 2 {
			log.Printf("Invalid message format received from %s: %s", contactAddress, message)
			continue
		}

		// Extract ID and content
		senderID := parts[0]
		content := parts[1]

		// Log and send the message to the channel for processing
		//log.Printf("Network received message: %s from %s (ID: %s)", content, contactAddress, senderID)
		network.MessageCh <- Message{ID: senderID, Content: content, Address: contactAddress}
	}
}

// sendMessage is a helper method to send messages to a contact using the existing UDP connection.
func (network *Network) sendMessage(contact *Contact, message string) {
	// Resolve the contact's address
	addr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		log.Printf("Error resolving UDP address: %v", err)
		return
	}

	message = contact.ID.String() + ":" + message
	// Use the existing connection to send the message
	_, err = network.Conn.WriteToUDP([]byte(message), addr)
	if err != nil {
		log.Printf("Error sending message to %s: %v", contact.Address, err)
		return
	}

	//log.Printf("Network sent message '%s' to %s", message, contact.Address)
}

// SendPingMessage sends a "ping" message to the contact.
func (network *Network) SendPingMessage(contact *Contact) {
	network.sendMessage(contact, "ping")
}

// SendPongMessage sends a "pong" message to the contact.
func (network *Network) SendPongMessage(contact *Contact) {
	network.sendMessage(contact, "pong")
}
