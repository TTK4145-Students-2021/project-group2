package orderDistribution

	
func (ol *OrderList) update(incomingOrderList OrderList)

	for i,curOrder := range ol {
		/*
		// the first if sentences checks for faiure is the version controll
		if incomingOrderIsNewer()
		if incomingStatus.OrderList[i].VersionNum == aes[curID].OrderList[i].VersionNum && incomingStatus.OrderList[i].HasOrder != aes[curID].OrderList[i].HasOrder {
			fmt.Println("There is smoething wrong with the version control of the order system")
			aes[curID].OrderList[i].HasOrder = true
			aes[curID].OrderList[i].VersionNum += 1
		}
		*/
		// this if statement checks if the order should be updated and updates it
		incomingOrder := incomingOrderList[i]
		if incomingOrderIsNewer(incomingOrder.VersionNum, curOrder.VersionNum){
			curOrder.HasOrder = incomingOrder.HasOrder
			curOrder.VersionNum = incomingOrder.VersionNum
		}
	}
}
	
func (ol *OrderList) init() {
	ol := [orderListLength]HallOrder{}

	// initalizing HallUp Orders
	initialCosts := [numElevators]int{10000}

	for idx := 0; idx < orderListLength; idx ++{
		ol[idx] = HallOrder{
			HasOrder:   false,
			Floor:      orderListIdxToFloor(idx),
			Direction:  orderListIdxToBtnTyp(idx),
			VersionNum: 0,
			Costs:      initialCosts,
			TimeStamp:  time.Now(),
		}
	}
}


func (ol *OrderList) calculateAllCosts()  {
	for idx,_ := range ol{
		curOrder := &ol[idx]
		if ol[idx].HasOrder{
			for _, curElevator := range allElevatorSatuses{
				cost := costFunction(curElevator, curOrder.Floor, curOrder.TimeStamp)
				curOrder.Costs[curElevator.ID] = cost
			}
		} else {
			for idx, _ := range curOrder.Costs {
				curOrder.Costs[idx] = maxCost
			}
		}
		
	}
}
