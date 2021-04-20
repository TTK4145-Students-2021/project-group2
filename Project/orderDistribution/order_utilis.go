package orderDistribution

import(
	"../messages"
	"../config"
	"fmt"
	"time"
	"errors"
)

/*
*=============================================================================
 * @Description:
 * Contains utility functions for use in the rest of the orderDistribution module
 *
/*=============================================================================
*/


func orderListIdxToFloor(idx int) int {
	if idx < _numFloors-1 {
		return idx
	} else {
		return idx - _numFloors + 2
	}
}

func orderListIdxToBtnTyp(idx int) ButtonType {
	if idx < _numFloors-1 {
		return BT_HallUp
	} else {
		return BT_HallDown
	}
}

func floorToOrderListIdx(floor int, dir messages.ButtonType) int {
	if dir == messages.BT_HallUp {
		return floor
	} 
	if dir == messages.BT_HallDown {
		return floor + _numFloors - 2
	} else {
		return -1
	}
}

func minumumRowCol(list [][_numElevators]int) (int, int) {
	min := config.MaxInt
	minRowIndex := -1
	minColIndex := -1
	for rowIndex, row := range list {
		for colIndex, number := range row {
			if number < min {
				min = number
				minRowIndex = rowIndex
				minColIndex = colIndex
			}
		}
	}
	return minRowIndex, minColIndex
}

func ifXinSliceInt(x int, slice []int) bool {
	for _, j := range slice {
		if x == j {
			return true
		}
	}
	return false
}

func checkIfElevatorOffline(allElevatorStatuses *[_numElevators]ElevatorStatus) {
	for{
		for elevID := range allElevatorStatuses {
			if time.Since(allElevatorStatuses[elevID].LastAlive) > config.OfflineTimeout && elevID != _thisID {
				if allElevatorStatuses[elevID].IsOnline == true {
					fmt.Printf("[-] Elevator offline ID: %02d\n", elevID)
				}
				allElevatorStatuses[elevID].IsOnline = false
			} else {
				if allElevatorStatuses[elevID].IsOnline == false {
					fmt.Printf("[+] Elevator online ID: %02d\n", elevID)
				}
				allElevatorStatuses[elevID].IsOnline = true
			}
		}
		time.Sleep(_checkOfflineIntervall)
	}
}


func checkIfEngineFails(assignedFloor int, status *ElevatorStatus, send_status chan<- ElevatorStatus) {
	faultCounter := 0
	curPos := -1
	lastPos := -1
	for {
		curPos = status.Pos
		if curPos == assignedFloor {
			return
		}
		if curPos != lastPos {
			lastPos = curPos
		} else {
			faultCounter += 1
			if faultCounter > _engineFailureTimeOut {
				err := errors.New("engine failure")
				panic(err)
			}
		}
		time.Sleep(time.Second)
	}
}

func canAcceptOrder(newAssignedOrder ButtonEvent, elevatorStatus ElevatorStatus, TimeOfLastCompletion time.Time)  bool {
	
	return (newAssignedOrder.Floor != -1 && 
		elevatorStatus.DoorOpen && !elevatorStatus.IsObstructed && 
		time.Since(TimeOfLastCompletion) > _orderTimeOut && 
		elevatorStatus.Pos != -1)
}

// Functions for printing out the Elevator staus and OrderList
func printOrderList(list [_orderListLength]HallOrder) {
	for i := 0; i < len(list); i++ {
		fmt.Println("Floor: ", list[i].Floor, " |direction: ", list[i].Direction, " |Order: ", list[i].HasOrder, " |Version", list[i].VersionNum, "Costs", list[i].Costs)
	}
}

func printElevatorStatus(status ElevatorStatus) {
	fmt.Println("ID:", status.ID, " |Currentfloor:", status.Pos, "| Door open:", status.DoorOpen, "|isOnline:", status.IsOnline, "|isObstructed:", status.IsObstructed, "\nCabOrders:", status.CabOrders)

	printOrderList(status.OrderList)
	fmt.Println("-----------------------------------------------")
}