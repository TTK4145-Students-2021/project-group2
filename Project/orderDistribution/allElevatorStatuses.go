package orderDistribution

import(
	"fmt"
	"time"
)

type AllElevatorStatuses [_numElevators]ElevatorStatus //numElevators går ikke her

func (aes *AllElevatorStatuses) init() {

	CabOrders := [_numFloors]bool{false}
	
	for elevID := 0; elevID < (_numElevators); elevID++ {
		OrderList := &OrderList{}
		OrderList.init()
		Status := ElevatorStatus{
			ID:           elevID,
			Pos:          1,
			OrderList:    *OrderList, 
			IsOnline:     false, 
			DoorOpen:     false,
			CabOrders:    CabOrders,
			IsAvailable:  false, 
			IsObstructed: false,
			// TimeOfLastCompletion: time.Now().Add(-time.Hour),
			// LastAlive:	   time.Now(),
		}
		aes[elevID] = Status
	}
}

//Updates the information list with the incoming update
func (aes *AllElevatorStatuses) update(incomingStatus ElevatorStatus) {
	aes[incomingStatus.ID] = incomingStatus
}

func (aes *AllElevatorStatuses) initStatus(chans OrderChannels){
	
	aes[_thisID].IsOnline = true
	aes[_thisID].IsAvailable = true

	chans.RequestElevatorStatus <- true
	initFloor := <-chans.NewFloor
	initDoor := <-chans.DoorOpen

	aes[_thisID].Pos = initFloor
	aes[_thisID].DoorOpen = initDoor
}

func (aes *AllElevatorStatuses) checkIfElevatorOffline() {
	for{
		for elevID := range aes {
			if time.Since(aes[elevID].LastAlive) > _offlineTimeout && elevID != _thisID {
				if aes[elevID].IsOnline == true {
					fmt.Printf("[-] Elevator offline ID: %02d\n", elevID)
				}
				aes[elevID].IsOnline = false
			} else {
				if aes[elevID].IsOnline == false {
					fmt.Printf("[+] Elevator online ID: %02d\n", elevID)
				}
				aes[elevID].IsOnline = true
			}
		}
		time.Sleep(_checkOfflineIntervall)
	}
}
