package orders

import "testing"

func TestBackup(t *testing.T) {

	backup := NewBackup("backup.txt")
	backup.RecoverFromBackup()
	backup.WriteToFile("Sloppy Seconds", false)
	backup.WriteToFile("Hey Yo!", false)
	backup.DeleteFile()
	//backup.ClearFile()
}
