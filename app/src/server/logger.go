package main

import (
	"log"
)

var serverLogToggle = true
var dbLogToggle = true

func sLog(s string) {
	if serverLogToggle == true {
		log.Println(s)
	}
}

func dbLog(s string) {
	if dbLogToggle == true {
		log.Println(s)
	}
}
