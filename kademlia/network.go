package kademlia

import (
	"log"
	"net"
)

// Network is a node in the network
type Network struct {
	// IP        string
	// Port      int
	Conn      *net.UDPConn
	MessageCh chan Message // Channel for UDP messages
}

// Message is a simple struct to hold a message and its sender's address.
type Message struct {
	Content string
	Address string
	ID      string
}

// Listen listens on a UDP address and stores incoming messages in the channel.
func (network *Network) Listen(Node Contact) error {
	// Create a UDPAddr based on the Network's IP and Port
	// addr := net.UDPAddr{
	// 	IP:   net.ParseIP(network.IP),
	// 	Port: network.Port,
	// }
	udpAddr, err := net.ResolveUDPAddr("udp", Node.Address) //addr)

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

		message := string(buffer[:n])
		contactAddress := remoteAddr.String()
		//log.Printf("Network received message: %s from %s", message, contactAddress)

		// Send the message to the channel for the Kademlia class to process
		network.MessageCh <- Message{Content: message, Address: contactAddress}
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
