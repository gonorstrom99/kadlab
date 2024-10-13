package kademlia

import (
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
	KademliaNode6 := CreateKademliaNode("127.0.0.1:8005")
	KademliaNode7 := CreateKademliaNode("127.0.0.1:8006")
	KademliaNode8 := CreateKademliaNode("127.0.0.1:8007")
	KademliaNode9 := CreateKademliaNode("127.0.0.1:8008")
	KademliaNode10 := CreateKademliaNode("127.0.0.1:8009")
	KademliaNode11 := CreateKademliaNode("127.0.0.1:8010")
	KademliaNode1.Start()
	KademliaNode2.Start()
	KademliaNode3.Start()
	KademliaNode4.Start()
	KademliaNode5.Start()
	KademliaNode6.Start()
	KademliaNode7.Start()
	KademliaNode8.Start()
	KademliaNode9.Start()
	KademliaNode10.Start()
	KademliaNode11.Start()

	time.Sleep(1 * time.Second)
	KademliaNode1.RoutingTable.AddContact(*KademliaNode2.RoutingTable.GetMe())
	KademliaNode2.RoutingTable.AddContact(*KademliaNode3.RoutingTable.GetMe())
	KademliaNode2.RoutingTable.AddContact(*KademliaNode4.RoutingTable.GetMe())
	KademliaNode2.RoutingTable.AddContact(*KademliaNode5.RoutingTable.GetMe())
	KademliaNode5.RoutingTable.AddContact(*KademliaNode6.RoutingTable.GetMe())
	KademliaNode6.RoutingTable.AddContact(*KademliaNode7.RoutingTable.GetMe())
	KademliaNode7.RoutingTable.AddContact(*KademliaNode8.RoutingTable.GetMe())
	KademliaNode8.RoutingTable.AddContact(*KademliaNode9.RoutingTable.GetMe())
	KademliaNode9.RoutingTable.AddContact(*KademliaNode10.RoutingTable.GetMe())
	KademliaNode10.RoutingTable.AddContact(*KademliaNode11.RoutingTable.GetMe())
	//lookupMessage := fmt.Sprintf("lookUpContact:%s:%s", KademliaNode1.Network.ID.String(), KademliaNode5.Network.ID.String())
	//KademliaNode1.Network.SendMessage(KademliaNode3.RoutingTable.GetMe(), lookupMessage)
	KademliaNode1.StartTask(KademliaNode11.RoutingTable.GetMe().ID, "LookupContact", "asd")
	time.Sleep(1 * time.Second)
	if !(KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode11.RoutingTable.GetMe())) {
		t.Fatalf("(File: kademlia_test, Test: TestIterativeLookupAndAdding) the iterative part of lookup does not work")
	}
}

// //TODO Dela upp networkjoining testet till två, ett som kollar att om
// nod 1 får en viss input så skickar den ut sina 3 noder, och ett som
// kollar om nod 5 får ett visst meddelande, att den lägger till noderna.
func TestNetworkJoining(t *testing.T) {

	KademliaNode1 := CreateKademliaNode("127.0.0.1:8011")
	KademliaNode2 := CreateKademliaNode("127.0.0.1:8012")
	KademliaNode3 := CreateKademliaNode("127.0.0.1:8013")
	KademliaNode4 := CreateKademliaNode("127.0.0.1:8014")
	KademliaNode5 := CreateKademliaNode("127.0.0.1:8015")
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
	KademliaNode5.RoutingTable.AddContact(*KademliaNode1.RoutingTable.GetMe())

	KademliaNode5.StartTask(KademliaNode5.RoutingTable.GetMe().ID, "LookupContact", "asd")

	// time.Sleep(1 * time.Second)
	// // if !(KademliaNode5.RoutingTable.IsContactInRoutingTable(KademliaNode2.RoutingTable.GetMe())) {
	// // 	t.Fatalf("(File: kademlia_test, Test: TestNetworkJoining) The network joining did not work")
	// // }
	// // if !(KademliaNode5.RoutingTable.IsContactInRoutingTable(KademliaNode3.RoutingTable.GetMe())) {
	// // 	t.Fatalf("(File: kademlia_test, Test: TestNetworkJoining) The network joining did not work")
	// // }
	// // if !(KademliaNode5.RoutingTable.IsContactInRoutingTable(KademliaNode4.RoutingTable.GetMe())) {
	// // 	t.Fatalf("(File: kademlia_test, Test: TestNetworkJoining) The network joining did not work")
	// // }
	// // if !(KademliaNode1.RoutingTable.IsContactInRoutingTable(KademliaNode5.RoutingTable.GetMe())) {
	// // 	t.Fatalf("(File: kademlia_test, Test: TestNetworkJoining) The network joining did not work")
	// // }

}

func TestFullBucket(t *testing.T) {

	KademliaNode1 := CreateKademliaNode("127.0.0.1:8016")
	KademliaNode2 := CreateKademliaNode("127.0.0.1:8017")
	KademliaNode3 := CreateKademliaNode("127.0.0.1:8018")
	KademliaNode4 := CreateKademliaNode("127.0.0.1:8019")
	KademliaNode5 := CreateKademliaNode("127.0.0.1:8020")
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
	KademliaNode5.RoutingTable.AddContact(*KademliaNode1.RoutingTable.GetMe())

	KademliaNode5.StartTask(KademliaNode5.RoutingTable.GetMe().ID, "LookupContact", "asd")

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
