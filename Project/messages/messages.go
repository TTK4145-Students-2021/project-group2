package messages

type ButtonType int

const (
	UNDEFINED   ButtonType = -1
	BT_HallUp   ButtonType = 0
	BT_HallDown ButtonType = 1
	BT_Cab      ButtonType = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type LampUpdate struct {
	Floor  int
	Button ButtonType
	Turn   bool
}
