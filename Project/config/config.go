package config

//Package for common variables to be defined across the application
var Port string
var ID int = 1 //the order module is written so that this is the first elevator. (i.e it has -1 in all indexes)
const NumFloors int = 4

const NumElevators int = 3

var BcastPort int = 12345
var PeersPort int = 54321

var SimPort string = "15657"
