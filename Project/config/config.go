package config

import (
	"os"
	"strconv"
)

//Package for common variables to be defined across the application
var Port string
var ID int
var NumFloors int = 4
var NumElevators int = 3

func InitConfig() {
	Port = os.Args[1]
	ID, _ = strconv.Atoi(os.Args[2])
}
