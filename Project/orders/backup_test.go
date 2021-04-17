package orders

import "testing"

func TestBackup(t *testing.T) {

	testSlice := [4]bool{false, false, true, true}

	backup := NewBackup("backup.txt")
	//backup.WriteToFile("Sloppy Seconds\n")
	//backup.WriteToFile("Hey Yo!\n")
	//backup.WriteToFile("Meh\n")
	backup.SaveCabOrders(testSlice)
	backup.RecoverCabOrders(&testSlice)
	//backup.DeleteFile()
	//backup.ClearFile()
}
