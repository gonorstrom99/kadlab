package kademlia

import (
	"log"
	"net"
)

type Network struct {
	Conn *net.UDPConn
}

func (n *Network) Listen(Node Contact) error {
	//addr := fmt.Sprintf("%s:%d", Node.IP, Node.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", Node.routingtable.me.Address) //addr)

	//udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Printf("Failed to resolve UDP address: %v", err)
		return err
	}

	n.Conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("Failed to start UDP listener: %v", err)
		return err
	}
	defer n.Conn.Close()

	buffer := make([]byte, 1024)
	for {
		nRead, remoteAddr, err := n.Conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Failed to read from UDP: %v", err)
			return err
		}
		receivedMessage := string(buffer[:nRead])
		log.Printf("Received '%s' from %s", receivedMessage, remoteAddr)
	}
}
func (n *Network) Ping(target Contact) error {
	//targetAddr := fmt.Sprintf("%s:%d", target.IP, target.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", target.routingtable.me.Address)
	if err != nil {
		log.Printf("Failed to resolve UDP address: %v", err)
		return err
	}

	message := "ping"
	_, err = n.Conn.WriteToUDP([]byte(message), udpAddr)
	if err != nil {
		log.Printf("Failed to send ping: %v", err)
		return err
	}
	log.Printf("Ping sent to %s", target.routingtable.me.Address)
	return nil
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
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