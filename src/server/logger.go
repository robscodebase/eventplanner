// logger.go contains variables and functions for turning logging on and off for certain
// actions.  Currently, server and db log control logging for respective functions.
package main

import (
	"log"
)

// If the toggles are set to true logging will be displayed
// in docker logs.
var serverLogToggle = true
var dbLogToggle = true

// sLog() controls the logging mechanism for functions that rely on the server.
func sLog(s string) {
	if serverLogToggle == true {
		log.Println(s)
	}
}

// dbLog() controls the logging mechanism for functions that rely on the database.
func dbLog(s string) {
	if dbLogToggle == true {
		log.Println(s)
	}
}
