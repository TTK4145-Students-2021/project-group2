package orderDistribution


func (es *ElevatorStatus) updateDoor(value bool) {
	es.DoorOpen = value
}

func (es *ElevatorStatus) updateFloor(pos int) {
	es.Pos = pos
}


func (es *ElevatorStatus) sendOutLightsUpdate(sendLampUpdate chan<- messages.LampUpdate) {
	for {
		for _, order := range es.OrderList {
			light := messages.LampUpdate{
				Floor:  order.Floor,
				Button: messages.ButtonType(order.Direction),
				Turn:   order.HasOrder,
			}
			sendLampUpdate <- light
		}
		for i, order := range es.CabOrders {

			light := messages.LampUpdate{
				Floor:  i,
				Button: messages.BT_Cab,
				Turn:   order,
			}
			sendLampUpdate <- light
		}
		time.Sleep(lampUpdateIntervall)
	}
}
/*
func (es *ElevatorStatus) sendOutStatus(channel chan<- ElevatorStatus) {
	es.Timestamp = time.Now()
	channel <- es
}
*/

func (es *ElevatorStatus) checkIfElevatorStops(assignedFloor int, send_status chan<- ElevatorStatus) {
	timeSinceMovment := 0
	lastFloor := -1
	for {    //
		if es.Pos  == assignedFloor {
			return
		}
		if es.Pos != lastFloor {
			lastFloor = es.Pos 
			time.Sleep(time.Second * 1)
		} else {
			timeSinceMovment += 1
			if timeSinceMovment > 10 {
				es.IsAvailable = false
				go sendOutStatus(send_status, *es)
				err := errors.New("engine failure")
				panic(err)
				return
			}
		}
	}
}

func (es *ElevatorStauts) broadcast(channel chan<- ElevatorStatus) {
	for {
		go sendOutStatus(channel, *es)
		time.Sleep(config.BcastIntervall)
		es.Timestamp = time.Now()
		channel <- es
	}
}

//updateOrderListButton
func (es *ElevatorStatus) updateOrders(btnEvent messages.ButtonEvent_message) {

	if btnEvent.Button == messages.BT_Cab {
		es.CabOrders[btnEvent.Floor] = true

	} else {
		idx := floorToOrderListIdx(btnEvent.Floor, btnEvent.Button)
		if es.OrderList[idx].HasOrder == false {
			es.OrderList[idx].HasOrder = true
			es.OrderList[idx].VersionNum += 1
		}
	}
}

