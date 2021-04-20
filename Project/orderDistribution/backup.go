package orderDistribution

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
/*=============================================================================
 * @Description:
 * This module contains functionality for saving cab-orders for the ordersystem.
 * These can then be read upon initialization.
 *
/*=============================================================================
*/

type BackupFile struct {
	ID    int
	Path  string
	mutex sync.Mutex
}

// Creates a new file with the given filename if one does not exist
func NewBackup(filename string, elevatorID int) *BackupFile {
	path := "./ID" + strconv.Itoa(elevatorID) + "_" + filename
	file := &BackupFile{
		Path: path,
	}

	// If file does not exist -> create file
	if _, err := os.Stat(file.Path); err == nil {
		fmt.Printf("File exists\n")
	} else {
		fmt.Printf("File does not exist\n")
		f, err := os.Create(file.Path) 
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		defer f.Close() 
	}

	return file
}

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

func (bf *BackupFile) SaveCabOrders(cabOrders [config.NumFloors]bool) {

	bf.ClearFile()

	for i := range cabOrders {
		indexString := strconv.Itoa(i)
		statusString := strconv.FormatBool(cabOrders[i])
		line := "Floor: " + indexString + " | CabOrderStatus: " + statusString + "\n"
		bf.WriteToFile(line)
	}
}

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

func (bf *BackupFile) ClearFile() {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	f, err := os.Create(bf.Path)
	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}
	defer f.Close()
}

func (bf *BackupFile) DeleteFile() {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	err := os.Remove(bf.Path)
	if err != nil {
		log.Fatal(err)
	}
}
