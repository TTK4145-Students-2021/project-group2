package orderDistribution

import(
	"../messages"
)

type AllElevatorStatuses [numElevators]messages.ElevatorStatus //numElevators går ikke her

func (aes *AllElevatorStatuses) initAllStatuses() {

	CabOrders := [numFloors]bool{false}
	OrderList := &OrderList{}
	OrderList.init()

	for elevID := 0; elevID < (numElevators); i++ {
		Status := ElevatorStatus{
			ID:           elevID,
			Pos:          1,
			OrderList:    OrderList,
			IsOnline:     false, //NB isOline is deafult false
			DoorOpen:     false,
			CabOrders:    CabOrders,
			IsAvailable:  false, //NB isAvaliable is deafult false
			IsObstructed: false,
		}
		aes[i] = Status
	}
}

//Updates the information list with the incoming update
func (aes *AllElevatorStatuses) update(incomingStatus ElevatorStatus) {
	aes[incomingStatus.ID] = incomingStatus
}

func (aes *Eleva) initStatus(chans OrderChannels){
	
	allElevatorSatuses[curID].IsOnline = true
	allElevatorSatuses[curID].IsAvailable = true

	chans.Get_init_position <- true
	initFloor := <-chans.New_floor
	initDoor := <-chans.Door_status

	allElevatorSatuses[curID].Pos = initFloor
	allElevatorSatuses[curID].DoorOpen = initDoor
}

func (aes *AllElevatorStatuses) checkIfElevatorOffline() {
	for{
		for elevID := range aes {
			if time.Since(aes[elevID].Timestamp) > config.OfflineTimeout && elevID != curID {
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
		time.Sleep(checkOfflineIntervall)
	}
}

func (aes *AllElevatorStatuses)updateOrderListCompleted(assignedOrder ButtonEvent, orderTimeOut *time.Time) {
	curFloor := list[curID].Pos

	//TODO CLEAN CODE ANNNND LOGIC

	// checks if the elevator has a cab call for the current floor
	if aes[curID].CabOrders[curFloor] == true && aes[curID].DoorOpen == true && assignedOrder.Floor==curFloor {
		aes[curID].CabOrders[curFloor] = false
		*orderTimeOut = time.Now()
		return
	}

	orderIdx := floorToOrderListIdx(assignedOrder.Floor, assignedOrder.Button)
	if orderIdx == -1{
		return
	}

	if aes[curID].OrderList[orderIdx].HasOrder && aes[curID].OrderList[orderIdx].Floor == curFloor {
		aes[curID].OrderList[orderIdx].HasOrder = false
		aes[curID].OrderList[orderIdx].VersionNum += 1
		*orderTimeOut = time.Now()
		// fmt.Println("Hall order compleete:", orderListIdxToFloor(orderIdx))
	}
}