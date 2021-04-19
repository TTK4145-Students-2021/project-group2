package orderDistribution

import(
	"../messages"
	"../config"
	"fmt"
)

func orderListIdxToFloor(idx int) int {
	if idx < _numFloors-1 {
		return idx
	} else {
		return idx - _numFloors + 2
	}
}

func orderListIdxToBtnTyp(idx int) ButtonType{
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

// Functions for printing out the Elevator staus and OrderList
func printOrderList(list [_orderListLength]HallOrder) {
	for i := 0; i < len(list); i++ {
		fmt.Println("Floor: ", list[i].Floor, " |direction: ", list[i].Direction, " |Order: ", list[i].HasOrder, " |Version", list[i].VersionNum, "Costs", list[i].Costs)
	}
}

func printElevatorStatus(status ElevatorStatus) {
	fmt.Println("ID:", status.ID, " |Currentfloor:", status.Pos, "| Door open:", status.DoorOpen, "|avaliable:", status.IsAvailable, "|isOnline:", status.IsOnline, "|isObstructed:", status.IsObstructed, "\nCabOrders:", status.CabOrders)

	printOrderList(status.OrderList)
	fmt.Println("-----------------------------------------------")
}