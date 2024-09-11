package kademlia

import (
	"log"
	"net"
)

type Network struct {
	NodeId string
	IP     string
	Port   int
	Conn   *net.UDPConn
}

// Listen listens on a UDP address and responds with "pong" when it receives "ping".
func (network *Network) Listen() {
	// Create a UDPAddr based on the Network's IP and Port
	addr := net.UDPAddr{
		IP:   net.ParseIP(network.IP),
		Port: network.Port,
	}

	// Start listening on the provided address
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Printf("Error starting UDP listener: %v\n", err)
		return
	}
	defer conn.Close()

	// Assign the connection to the network's Conn field
	network.Conn = conn

	buffer := make([]byte, 1024)

	for {
		// Read incoming UDP message
		n, remoteAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading from UDP connection: %v", err)
			continue
		}

		message := string(buffer[:n])
		log.Printf("Received message: %s from %s", message, remoteAddr.String())

		// If the message is "ping", respond with "pong"
		if message == "ping" {
			contact := Contact{
				Address: remoteAddr.String(), // Use the remote address (IP:Port) as Contact's Address
			}
			network.SendPongMessage(&contact)
		}
	}
}

// send ping to contact
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

	// Buffer to read the "pong" response
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("Error receiving pong: %v", err)
		return
	}

	response := string(buffer[:n])
	if response == "pong" {
		log.Printf("Received pong from %s", contact.Address)
	} else {
		log.Printf("Unexpected response: %s", response)
	}
}

// Sends a pong message to a contact
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

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
