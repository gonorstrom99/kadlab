package kademlia

import (
	"log"
	"net"
)

// Network is a node in the network
type Network struct {
	IP   string
	Port int
	Conn *net.UDPConn
}

// Listen listens on a UDP address and returns the received message and the contact's address.
func (network *Network) Listen() (string, string, error) {
	// Create a UDPAddr based on the Network's IP and Port
	addr := net.UDPAddr{
		IP:   net.ParseIP(network.IP),
		Port: network.Port,
	}

	// Start listening on the provided address
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Printf("Error starting UDP listener: %v\n", err)
		return "", "", err
	}
	defer conn.Close()

	// Assign the connection to the network's Conn field
	network.Conn = conn

	buffer := make([]byte, 1024)

	// Read incoming UDP message
	n, remoteAddr, err := conn.ReadFrom(buffer)
	if err != nil {
		log.Printf("Error reading from UDP connection: %v", err)
		return "", "", err
	}

	message := string(buffer[:n])
	contactAddress := remoteAddr.String()

	return message, contactAddress, nil
}

// SendPingMessage sends a "ping" message to the contact.
func (network *Network) SendPingMessage(contact *Contact) {
	// Parse the contact's address (expected to be in "IP:Port" format)
	addr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		log.Printf("Error resolving UDP address: %v", err)
		return
	}

	// Dial UDP to the contact's address
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("Error dialing UDP: %v", err)
		return
	}
	defer conn.Close()

	// Send the "ping" message
	_, err = conn.Write([]byte("ping"))
	if err != nil {
		log.Printf("Error sending ping: %v", err)
		return
	}
}

// SendPongMessage sends a "pong" message to the contact.
func (network *Network) SendPongMessage(contact *Contact) {
	// Parse the contact's address to send the pong message
	addr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		log.Printf("Error resolving address: %v", err)
		return
	}

	// Dial UDP to the contact's address
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("Error dialing UDP: %v", err)
		return
	}
	defer conn.Close()

	// Send the "pong" message
	_, err = conn.Write([]byte("pong"))
	if err != nil {
		log.Printf("Error sending pong: %v", err)
	}
}

// SendFindContactMessage sends a "FIND_NODE" message to a contact requesting nodes close to a target ID.
func (network *Network) SendFindContactMessage(contact *Contact, targetNodeID string) {
	// Parse the contact's address (expected to be in "IP:Port" format)
	addr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		log.Printf("Error resolving UDP address: %v", err)
		return
	}

	// Dial UDP to the contact's address
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("Error dialing UDP: %v", err)
		return
	}
	defer conn.Close()

	// Send the "FIND_NODE" message, including the targetNodeID
	message := "FIND_NODE:" + targetNodeID
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Printf("Error sending FIND_NODE message: %v", err)
		return
	}

	// Buffer to read the response
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("Error receiving FIND_NODE response: %v", err)
		return
	}

	response := string(buffer[:n])
	log.Printf("Received response: %s", response)
	// You would then parse the response for the closest contacts
}

// SendFindDataMessage sends a "FIND_DATA" message to search for data by hash.
func (network *Network) SendFindDataMessage(contact *Contact, hash string) {
	// Parse the contact's address (expected to be in "IP:Port" format)
	addr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		log.Printf("Error resolving UDP address: %v", err)
		return
	}

	// Dial UDP to the contact's address
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("Error dialing UDP: %v", err)
		return
	}
	defer conn.Close()

	// Send the "FIND_DATA" message, including the hash of the data
	message := "FIND_DATA:" + hash
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Printf("Error sending FIND_DATA message: %v", err)
		return
	}

	// Buffer to read the response
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("Error receiving FIND_DATA response: %v", err)
		return
	}

	response := string(buffer[:n])
	log.Printf("Received data or node list: %s", response)
	// You would then parse the response for the data or the closest nodes to the hash
}

// SendStoreMessage sends a "STORE" message to store data on a contact.
func (network *Network) SendStoreMessage(contact *Contact, data []byte, hash string) {
	// Parse the contact's address (expected to be in "IP:Port" format)
	addr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		log.Printf("Error resolving UDP address: %v", err)
		return
	}

	// Dial UDP to the contact's address
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("Error dialing UDP: %v", err)
		return
	}
	defer conn.Close()

	// Send the "STORE" message, including the hash of the data
	message := "STORE:" + hash + ":" + string(data)
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Printf("Error sending STORE message: %v", err)
		return
	}

	log.Printf("Stored data on contact %s", contact.Address)
}
