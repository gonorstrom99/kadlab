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
	/*shortlist[]=FindClosestContacts(target.ID, 3)
	if shortlist[0]<closestNode{
		closestNode=shortlist[0]
	}
	else{

	}
	
	if (closestNode == target){
		return closestNode
	}
	else {
		for i=0, (i<shortlist[].len()), i++{
			LookupContact(target, closestNode)
		}
	}*/
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}