package orders

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Souces:
// https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file-using-go
// https://tutorialedge.net/golang/reading-writing-files-in-go/
// https://www.golangprograms.com/how-to-read-write-from-to-file-in-golang.html

// Backup containts all the necessary parameters for the backupsystem
type BackupFile struct {
	Path  string
	mutex sync.Mutex
}

// NewBackup initializes and returns a new backupfile.
// Creates a new file with the given filename if one does not exist
func NewBackup(filename string) *BackupFile {

	path := "./" + filename

	file := &BackupFile{
		Path: path,
	}

	// If file does not exist -> create file
	if _, err := os.Stat(file.Path); err == nil {
		fmt.Printf("File exists\n")
	} else {
		fmt.Printf("File does not exist\n")

		f, err := os.Create(file.Path) // Truncates if file already exists, be careful!
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		defer f.Close() // Make sure to close the file when you're done
	}

	return file
}

// RecoverFromBackup reads the file and updates a CabOrder-list
func (bf *BackupFile) RecoverFromBackup( /*[config.NumFloors]CabOrders*/ ) {

	// TODO: Implement functionality for returning everything from the recovery system

}

// WriteToBackup takes in a CabOrder-list and writes them to the backupfile
func (bf *BackupFile) WriteToBackup( /*[config.NumFloors]CabOrders*/ ) {

}

// WriteToFile writes a line to file. Option to add a newLine automatically
func (bf *BackupFile) WriteToFile(line string, newLine bool) {
	if newLine {
		line = line + "\n"
	}

	f, err := os.OpenFile(bf.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(line); err != nil {
		panic(err)
	}
}

// ClearFile clears the entire file
func (bf *BackupFile) ClearFile() {
	f, err := os.Create(bf.Path)
	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}
	defer f.Close()
}

// DeleteFile deletes the file
func (bf *BackupFile) DeleteFile() {
	err := os.Remove(bf.Path)
	if err != nil {
		log.Fatal(err)
	}
}
