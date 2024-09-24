package kademlia

import (
	"testing"
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
