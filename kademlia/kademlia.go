package kademlia

import "log"

type Kademlia struct {
	//MeContact Contact 		//take from routing table
	//NodeID       string
	//IP           string
	//Port         int
	routingtable RoutingTable
	network      Network
}

func (kademlia *Kademlia) Ping(target Contact) error {
	err := kademlia.network.Ping(target)
	if err != nil {
		log.Printf("error in the kademlia ping function", err)
		return err
	}
	return nil
}

func (kademlia *Kademlia) Listen() error {
	//network Listen
	err := kademlia.network.Listen(kademlia.routingtable.me)
	if err != nil {
		log.Printf("error in the kademlia listen function", err)
		return err
	}
	return nil
}
func (kademlia *Kademlia) LookupContact(target *Contact) {
	//TODO
	/*candidates=routingtable.me.ContactCandidates
	closestNode := kademlia.routingtable.me
	candidates.append(routingtable.FindClosestContacts(target.ID, 3))
	if (candidates[0].Less(closestNode)){
		closestNode = candidates[0]
	} */

}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
