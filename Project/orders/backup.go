package orders

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"../config"
)

// BackupFile containts all the necessary parameters for the backupsystem
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

// RecoverCabOrders reads the file and updates a CabOrder-list
func (bf *BackupFile) RecoverCabOrders(cabOrders *[config.NumFloors]bool) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	file, err := os.Open(bf.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineParsed := strings.Split(line, " ")
		index := lineParsed[1]
		indexInt, _ := strconv.Atoi(index)
		status := lineParsed[4]
		cabOrders[indexInt], _ = strconv.ParseBool(status)
		fmt.Println("index: " + index + " | status: " + status)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// SaveCabOrders takes in a CabOrder-list and writes them to the backupfile
func (bf *BackupFile) SaveCabOrders(cabOrders [config.NumFloors]bool) {

	bf.ClearFile()

	for i := range cabOrders {
		indexString := strconv.Itoa(i)
		statusString := strconv.FormatBool(cabOrders[i])
		line := "Floor: " + indexString + " | CabOrderStatus: " + statusString + "\n"
		bf.WriteToFile(line)
	}
}

// WriteToFile writes a line to file.
func (bf *BackupFile) WriteToFile(line string) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

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
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	f, err := os.Create(bf.Path)
	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}
	defer f.Close()
}

// DeleteFile deletes the file
func (bf *BackupFile) DeleteFile() {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	err := os.Remove(bf.Path)
	if err != nil {
		log.Fatal(err)
	}
}
