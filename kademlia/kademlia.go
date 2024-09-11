package kademlia

// Struct holding nodeID and its routingtable
type Kademlia struct {
	ID           *KademliaID
	routingTable []*bucket
}

// LookupContact searches for the closest contacts to the target node.
func (kademlia *Kademlia) LookupContact(target *Contact) []*Contact {
	closestContacts := kademlia.routingTable.FindClosestContacts(target.ID, bucketSize)
	// Perform iterative search on these closest contacts
	for _, contact := range closestContacts {
		// Query each contact to find closer nodes (this is typically done via network communication)
		responseContacts := kademlia.sendFindNodeMessage(contact, target)
		closestContacts = mergeAndSortContacts(closestContacts, responseContacts, target.ID)
	}
	return closestContacts[:bucketSize] // Return the closest ones
}

// LookupData retrieves data associated with the given hash from the Kademlia network,
// querying nodes until the value is found or no closer nodes are available.
func (kademlia *Kademlia) LookupData(hash string) []byte {
	closestContacts := kademlia.routingTable.FindClosestContacts(HashToID(hash), bucketSize)
	// Query each contact for the data
	for _, contact := range closestContacts {
		data := kademlia.sendFindValueMessage(contact, hash)
		if data != nil {
			return data // Return data if found
		}
	}
	return nil // Data not found
}

// Store saves the given data into the Kademlia network,
// distributing it to nodes that are closest to the hash of the data.
func (kademlia *Kademlia) Store(data []byte) {
	hash := HashData(data)
	closestContacts := kademlia.routingTable.FindClosestContacts(HashToID(hash), bucketSize)
	// Send store message to closest contacts
	for _, contact := range closestContacts {
		kademlia.sendStoreMessage(contact, hash, data)
	}
}
