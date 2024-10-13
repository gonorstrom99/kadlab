package kademlia

// import (
// 	"net"
// 	"testing"
// 	"time"
// )

// // Test the SendMessage function by using real UDP connections
// func TestSendMessage(t *testing.T) {
// 	// Create the receiving UDP server (listener)
// 	addr := "127.0.0.1:0" // Use :0 to let the OS choose a free port
// 	udpAddr, err := net.ResolveUDPAddr("udp", addr)
// 	if err != nil {
// 		t.Fatalf("Failed to resolve UDP address: %v", err)
// 	}

// 	conn, err := net.ListenUDP("udp", udpAddr)
// 	if err != nil {
// 		t.Fatalf("Failed to start UDP listener: %v", err)
// 	}
// 	defer conn.Close()

// 	// Get the actual port assigned by the OS
// 	actualAddr := conn.LocalAddr().String()

// 	// Set up the network with the real UDP connection
// 	messageCh := make(chan Message)
// 	network := &Network{
// 		MessageCh: messageCh,
// 		Conn:      conn,
// 	}

// 	// Create a contact for testing
// 	contact := NewContact(NewRandomKademliaID(), actualAddr)

// 	// Set up a goroutine to simulate the receiver
// 	go func() {
// 		buffer := make([]byte, 1024)
// 		n, remoteAddr, err := conn.ReadFromUDP(buffer)
// 		if err != nil {
// 			t.Errorf("Error reading from UDP connection: %v", err)
// 			return
// 		}
// 		receivedMessage := string(buffer[:n])

// 		expectedMessage := "ping:FFFFFFFF00000000000000000000000000000000:pingInfo"
// 		if receivedMessage != expectedMessage {
// 			t.Errorf("Expected message %s, but got %s", expectedMessage, receivedMessage)
// 		}

// 		// Log the remote address to check if it's correct
// 		t.Logf("Received message from %s", remoteAddr.String())
// 	}()

// 	// Allow some time for the listener to be ready
// 	time.Sleep(100 * time.Millisecond)

// 	// Send a ping message to the listener
// 	pingMessage := "ping:FFFFFFFF00000000000000000000000000000000:pingInfo"
// 	network.SendMessage(&contact, pingMessage)

// 	// Allow some time for the message to be processed
// 	time.Sleep(100 * time.Millisecond)
// }

// // Test the Listen function to ensure it properly handles incoming messages
// func TestListen(t *testing.T) {
// 	// Create the receiving UDP server (listener)
// 	addr := "127.0.0.1:0" // Use :0 to let the OS choose a free port
// 	udpAddr, err := net.ResolveUDPAddr("udp", addr)
// 	if err != nil {
// 		t.Fatalf("Failed to resolve UDP address: %v", err)
// 	}

// 	conn, err := net.ListenUDP("udp", udpAddr)
// 	if err != nil {
// 		t.Fatalf("Failed to start UDP listener: %v", err)
// 	}
// 	defer conn.Close()

// 	// Set up the network with the real UDP connection
// 	messageCh := make(chan Message)
// 	network := &Network{
// 		MessageCh: messageCh,
// 		Conn:      conn,
// 	}

// 	// Create a contact for testing
// 	contact := NewContact(NewRandomKademliaID(), conn.LocalAddr().String())

// 	// Set up a goroutine to simulate the Listen function
// 	go func() {
// 		err := network.Listen(contact)
// 		if err != nil {
// 			t.Errorf("Error in Listen: %v", err)
// 		}
// 	}()

// 	// Allow some time for the listener to be ready
// 	time.Sleep(100 * time.Millisecond)

// 	// Send a mock message to the listener
// 	mockMessage := "ping:FFFFFFFF00000000000000000000000000000000:pingInfo"
// 	sendAddr, err := net.ResolveUDPAddr("udp", conn.LocalAddr().String())
// 	if err != nil {
// 		t.Fatalf("Failed to resolve UDP address for sending: %v", err)
// 	}
// 	sendConn, err := net.DialUDP("udp", nil, sendAddr)
// 	if err != nil {
// 		t.Fatalf("Failed to dial UDP: %v", err)
// 	}
// 	defer sendConn.Close()

// 	// Send the message
// 	_, err = sendConn.Write([]byte(mockMessage))
// 	if err != nil {
// 		t.Fatalf("Failed to send message: %v", err)
// 	}

// 	// Allow some time for the message to be processed
// 	time.Sleep(100 * time.Millisecond)

// 	// Check if the message was sent to the MessageCh channel
// 	select {
// 	case receivedMessage := <-network.MessageCh:
// 		if receivedMessage.Command != "ping" {
// 			t.Errorf("Expected 'ping', but got '%s'", receivedMessage.Command)
// 		}
// 		if receivedMessage.SenderID != "FFFFFFFF00000000000000000000000000000000" {
// 			t.Errorf("Expected senderID 'FFFFFFFF00000000000000000000000000000000', but got '%s'", receivedMessage.SenderID)
// 		}
// 		if receivedMessage.CommandInfo != "pingInfo" {
// 			t.Errorf("Expected 'pingInfo', but got '%s'", receivedMessage.CommandInfo)
// 		}
// 	default:
// 		t.Error("Expected message in channel but received none")
// 	}
// }
