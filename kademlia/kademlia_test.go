package kademlia

import (
	"fmt"
	"testing"
	"time"
)

// Example function in kademlia.go
// func CheckContactStatus(contact *Contact) bool {
//     // some implementation
// }

func TestKademliaContactAdding(t *testing.T) {
	KademliaNode1 := CreateKademliaNode("127.0.0.1:8000")
	KademliaNode2 := CreateKademliaNode("127.0.0.1:8001")
	KademliaNode1.RoutingTable.AddContact(*KademliaNode2.RoutingTable.GetMe())
	KademliaNode1.RoutingTable.AddContact(*KademliaNode2.RoutingTable.GetMe())
	bucketIndex := KademliaNode1.RoutingTable.getBucketIndex(KademliaNode2.RoutingTable.me.ID)
	bucket := KademliaNode1.RoutingTable.getBucket(bucketIndex)
	if bucket.Len() > 1 {
		t.Fatalf("the same contact got added twice somehow")
	}
	if bucket.Len() < 1 {
		t.Fatalf("we tried adding a contact but it did not get added")
	}

	// Add more tests for other functions
}

// This is not an unit test
func TestIterativeLookupAndAdding(t *testing.T) {

	KademliaNode1 := CreateKademliaNode("127.0.0.1:8000")
	KademliaNode2 := CreateKademliaNode("127.0.0.1:8001")
	KademliaNode3 := CreateKademliaNode("127.0.0.1:8002")
	KademliaNode4 := CreateKademliaNode("127.0.0.1:8003")
	KademliaNode5 := CreateKademliaNode("127.0.0.1:8004")
	time.Sleep(1 * time.Second)
	KademliaNode1.Start()
	KademliaNode2.Start()
	KademliaNode3.Start()
	KademliaNode4.Start()
	KademliaNode5.Start()
	time.Sleep(1 * time.Second)

	KademliaNode1.RoutingTable.AddContact(*KademliaNode2.RoutingTable.GetMe())
	KademliaNode1.RoutingTable.AddContact(*KademliaNode3.RoutingTable.GetMe())
	KademliaNode3.RoutingTable.AddContact(*KademliaNode4.RoutingTable.GetMe())

	KademliaNode4.RoutingTable.AddContact(*KademliaNode5.RoutingTable.GetMe())
	lookupMessage := fmt.Sprintf("lookUpContact:%s:%s", KademliaNode1.Network.ID.String(), KademliaNode5.Network.ID.String())
	KademliaNode1.Network.SendMessage(KademliaNode3.RoutingTable.GetMe(), lookupMessage)

	time.Sleep(1 * time.Second)
	if !(KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode5.RoutingTable.GetMe())) {
		t.Fatalf("(File: kademlia_test, Test: TestIterativeLookupAndAdding) the iterative part of lookup does not work")
	}
}
func TestNetworkJoining(t *testing.T) {

	KademliaNode1 := CreateKademliaNode("127.0.0.1:8000")
	KademliaNode2 := CreateKademliaNode("127.0.0.1:8001")
	KademliaNode3 := CreateKademliaNode("127.0.0.1:8002")
	KademliaNode4 := CreateKademliaNode("127.0.0.1:8003")
	KademliaNode5 := CreateKademliaNode("127.0.0.1:8004")
	time.Sleep(1 * time.Second)
	KademliaNode1.Start()
	KademliaNode2.Start()
	KademliaNode3.Start()
	KademliaNode4.Start()
	KademliaNode5.Start()
	time.Sleep(1 * time.Second)

	KademliaNode1.RoutingTable.AddContact(*KademliaNode2.RoutingTable.GetMe())
	KademliaNode1.RoutingTable.AddContact(*KademliaNode3.RoutingTable.GetMe())
	KademliaNode1.RoutingTable.AddContact(*KademliaNode4.RoutingTable.GetMe())

	//KademliaNode5.RoutingTable.AddContact(*KademliaNode5.RoutingTable.GetMe())
	lookupMessage := fmt.Sprintf("lookUpContact:%s:%s", KademliaNode5.Network.ID.String(), KademliaNode5.Network.ID.String())
	KademliaNode5.Network.SendMessage(KademliaNode1.RoutingTable.GetMe(), lookupMessage)

	time.Sleep(1 * time.Second)
	if !(KademliaNode5.RoutingTable.IsContactInRoutingTable(KademliaNode2.RoutingTable.GetMe())) {
		t.Fatalf("(File: kademlia_test, Test: TestNetworkJoining) The network joining did not work")
	}
	if !(KademliaNode5.RoutingTable.IsContactInRoutingTable(KademliaNode3.RoutingTable.GetMe())) {
		t.Fatalf("(File: kademlia_test, Test: TestNetworkJoining) The network joining did not work")
	}
	if !(KademliaNode5.RoutingTable.IsContactInRoutingTable(KademliaNode4.RoutingTable.GetMe())) {
		t.Fatalf("(File: kademlia_test, Test: TestNetworkJoining) The network joining did not work")
	}
	if !(KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode5.RoutingTable.GetMe())) {
		t.Fatalf("(File: kademlia_test, Test: TestNetworkJoining) The network joining did not work")
	}

}
