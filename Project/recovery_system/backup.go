package recovery_system

// Backup containts all the necessary parameters for the backupsystem
type Backup struct {
	// TODO: Create all necessary struct parameters
	// E.g. File?
}

// NewBackup returns a new backupstruct
func NewBackup() *Backup {
	return &Backup{
		//TODO:Initialize with proper parameters
	}
}

// InitializeBackup initializes a new file with a given filename
func (backup *Backup) InitializeBackup(filename string) {
	// TODO: Set initializing parameters
	// Potentially unite with the NewBackup function
}

// WriteToFile writes a line to file
func (backup *Backup) WriteToFile(line string) {

	// TODO: Implement functionality for writing line to file
}

// RecoverFromBackup reads the file and includes it in order
func (backup *Backup) RecoverFromBackup() {

	// TODO: Implement functionality for returning everything from the recovery system
}

/*
* OTHER HELPER FUNCTIONS
 */

// ClearFile clears the entire file
func (backup *Backup) ClearFile() {

	// TODO: Clear content of file, but leave it intact?
}

// DeleteFile deletes the file
func (backup *Backup) DeleteFile() {

	// TODO: Delete file
}
