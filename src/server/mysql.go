// Copyright (c) 2018 Robert Reyna. All rights reserved.
// License BSD 3-Clause https://github.com/robscodebase/eventplanner/blob/master/LICENSE.md
// mysql.go contains the functions for interacting with the database.
package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var createDBstmt = []string{
	`CREATE DATABASE IF NOT EXISTS eventplanner;`,
	`USE eventplanner;`,
	`CREATE TABLE IF NOT EXISTS events (
     id INT UNSIGNED NOT NULL AUTO_INCREMENT,
     name VARCHAR(255) NULL,
     starttime VARCHAR(255) NULL,
     endtime VARCHAR(255) NULL,
     description VARCHAR(255) NULL,
     userid INT UNSIGNED NOT NULL,
     PRIMARY KEY (id)
     );`,
	`CREATE TABLE IF NOT EXISTS users (
     id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	 username VARCHAR(255),
	 secret BINARY(255),
	 cookiesession VARCHAR(255),
	 PRIMARY KEY (id)
 	 );`,
}

// registerDB() opens the db and returns the db instance.
func registerDB() (*sql.DB, error) {
	var err error
	dbLog(fmt.Sprintf("mysql.go: registerDB()"))
	// retries give time for docker-compose and mysql to finish setup.
	db, err = sql.Open("mysql", dbLogIn)
	if err != nil {
		return nil, fmt.Errorf("mysql.go: registerDB(): sql.Open db: %v: err: %v", db, err)
	}
	dbLog(fmt.Sprintf("mysql.go: registerDB(): sql.Open success: db: %v", db))
	return db, nil
}

// isDB() pings, and checks that the db and table have been created.
func isDB(db *sql.DB) error {
	dbLog(fmt.Sprintf("mysql.go: isDB()"))
	// Ping the db with db.Ping().
	dbLog(fmt.Sprintf("mysql.go: isDB(): ping db"))
	if db.Ping() == driver.ErrBadConn {
		return fmt.Errorf("mysql.go: isDB() db.Ping() error: could not ping database.")
	}

	// Try to use the db if there is an error 1049 create the db.
	dbLog(fmt.Sprintf("mysql.go: isDB(): use db"))
	if _, err := db.Exec("USE eventplanner"); err != nil {
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {
			return createDB(db)
		}
	}

	// Try to describe the table if there is an error 1146 create the table.
	dbLog(fmt.Sprintf("mysql.go: isDB(): describe table"))
	if _, err := db.Exec("DESCRIBE events"); err != nil {
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			return createDB(db)
		}
		// Unknown error.
		return fmt.Errorf("mysql.go: isDB(): db.Exec: error: %v", err)
	}
	dbLog(fmt.Sprintf("mysql.go: isDB(): db is ready"))
	return nil
}

// createDB() executes the CREATE DATABASE and CREATE TABLE commands.
func createDB(db *sql.DB) error {
	dbLog(fmt.Sprintf("mysql.go: createDB(): var db: %v", db))
	for _, sqlCommand := range createDBstmt {
		_, err := db.Exec(sqlCommand)
		if err != nil {
			return fmt.Errorf("mysql.go: createDB(): db.Exec: problem exec command: error: %v", err)
		}
	}
	return nil
}

// createDemoDB() add the demo events and user to the db.
func createDemoDB(db *sql.DB) (int, string, error) {
	var eventsCreated int
	dbLog(fmt.Sprintf("mysql.go: createDemoDB(): var db: %v", db))
	// Range over each of the demoEvents and insert them in to the db.
	for _, demo := range demoEvents {
		// Prepare stmt.
		insertDemoEvent, err := db.Prepare("INSERT INTO events (name, starttime, endtime, description, userid) VALUES(?, ?, ?, ?, ?)")
		dbLog(fmt.Sprintf("mysql.go: createDemoDB(): insertDemoEvent: %v", insertDemoEvent))
		// Insert demo event into db.
		results, err := insertDemoEvent.Exec(demo.Name, demo.StartTime, demo.EndTime, demo.Description, demo.UserID)
		if err != nil {
			dbLog(fmt.Sprintf("mysql.go: createDemoDB(): problem creating demo db most likely entries already exist: %v", err))
			return eventsCreated, "", err
		}
		dbLog(fmt.Sprintf("mysql.go: createDemoDB(): insertDemoEvent success: results: %v", results))
		eventsCreated++
	}

	// Prepare insert stmt.
	insertDemoUser, err := db.Prepare("INSERT INTO users (username, secret) VALUES(?, ?)")
	dbLog(fmt.Sprintf("mysql.go: createDemoDB(): insertDemoUser: %v", insertDemoUser))
	// Encrypt password.
	secret, err := bcrypt.GenerateFromPassword(demoUser.Secret, bcrypt.DefaultCost)
	if err != nil {
		log.Panicf("mysql.go: createDemoDB(): error encrypting demo password secret: %v: error: %v", secret, err)
	}
	sLog(fmt.Sprintf("mysql.go: createDemoDB(): encrypted password: %v", secret))
	// Insert user into db.
	results, err := insertDemoUser.Exec(demoUser.Username, secret)
	if err != nil {
		dbLog(fmt.Sprintf("mysql.go: createDemoDB(): problem creating demo user most likely entries already exist: %v", err))
		return eventsCreated, "", err
	}
	dbLog(fmt.Sprintf("mysql.go: createDemoDB(): var results: %v", results))
	return eventsCreated, demoUser.Username, nil
}

// rowScanner() is implemented by sql.Row and sql.Rows
type eventScanner interface {
	Scan(scanTo ...interface{}) error
}

// scanEvent() reads an event from a sql.Row or sql.Rows
func scanEvent(eventScan eventScanner) (*Event, error) {
	dbLog(fmt.Sprintf("mysql.go: scanEvent(): eventScan: %v", eventScan))
	var (
		id          int64
		name        sql.NullString
		startTime   sql.NullString
		endTime     sql.NullString
		description sql.NullString
		userID      int64
	)
	if err := eventScan.Scan(&id, &name, &startTime, &endTime, &description, &userID); err != nil {
		return nil, fmt.Errorf("mysql.go: scanEvent(): eventScan.Scan(): error: %v", err)
	}
	event := &Event{
		ID:          id,
		Name:        name.String,
		StartTime:   startTime.String,
		EndTime:     endTime.String,
		Description: description.String,
		UserID:      userID,
	}
	dbLog(fmt.Sprintf("mysql.go: scanEvent(): event scan success: %v", event))
	return event, nil
}

// listEvents() returns all events created by the current user.
func listEvents(db *sql.DB, username string) ([]*Event, error) {
	dbLog(fmt.Sprintf("mysql.go: listEvents(): username: %v", username))
	rows, err := db.Query("SELECT * FROM events WHERE userid = ?", 1)
	//rows, err := db.Query("SELECT * FROM events")
	if err != nil {
		return nil, fmt.Errorf("mysql.go: listEvents(): db.Query(): error: %v", err)
	}
	dbLog(fmt.Sprintf("mysql.go: listEvents() db.Query() success: rows: %v", rows))
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		dbLog(fmt.Sprintf("mysql.go: rows.Next() inside loop: rows: %v", rows))
		event, err := scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("mysql.go: listEvents(): rows.Next(): error: %v", err)
		}
		dbLog(fmt.Sprintf("mysql.go: listEvents() rows.Next(): appending event: %v", event))
		events = append(events, event)
	}
	dbLog(fmt.Sprintf("mysql.go: rows.Next() success: events: %v", events))
	return events, nil
}

// listEvent() returns one event created by the current user.
func listEvent(db *sql.DB, userID, eventID int64) (*Event, error) {
	dbLog(fmt.Sprintf("mysql.go: listEvent() eventID: %v", eventID))
	event, err := scanEvent(db.QueryRow("SELECT * FROM events WHERE id = ? AND userid = ?", eventID, userID))
	if err != nil {
		return nil, fmt.Errorf("mysql.go: listEvent(): error getting event: %v: err: %v: eventID: %v", event, err, eventID)
	}
	return event, nil
}

// updateEvent() updates an event based on a userid and eventid.
func updateEvent(db *sql.DB, event *Event) error {
	dbLog(fmt.Sprintf("mysql.go: updateEvent() eventID: %v", event.ID))
	// Prepare insert stmt.
	updateDBEvent, err := db.Prepare("UPDATE events SET name=?, starttime=?, endtime=?, description=?  WHERE id=? AND userid=?")
	dbLog("mysql.go: updateEvent() db.Prepare() complete")
	// Insert demo event into db.
	dbLog(fmt.Sprintf("mysql.go: updateEvent() event.ID: %v, event.Name %v, event.StartTime %v, event.EndTime %v, event.Description %v, event.UserID %v",
		event.ID, event.Name, event.StartTime, event.EndTime, event.Description, event.UserID))
	results, err := updateDBEvent.Exec(event.Name, event.StartTime, event.EndTime, event.Description, event.ID, event.UserID)
	if err != nil {
		return fmt.Errorf("mysql.go: updateEvent(): problem updating db could be wrong userid: error: %v, userID: %v", err, event.UserID)
	}
	dbLog(fmt.Sprintf("mysql.go: updateEvent(): success: results: %v", results))
	return nil
}

// addEventDB() adds a new event to the db.
func addEventDB(db *sql.DB, event *Event) error {
	dbLog(fmt.Sprintf("mysql.go: addEvent() event: %v", event))
	// Prepare delete stmt.
	addDBEvent, err := db.Prepare("INSERT INTO events (name, starttime, endtime, description, userid) VALUES (?, ?, ?, ?, ?)")
	dbLog("mysql.go: addEvent(): db.Prepare(): complete")
	// Insert demo event into db.
	dbLog(fmt.Sprintf("mysql.go: addEvent(): userID: %v", event.UserID))
	results, err := addDBEvent.Exec(event.Name, event.StartTime, event.EndTime, event.Description, event.UserID)
	if err != nil {
		return fmt.Errorf("mysql.go: addEvent(): problem adding event: error: %v, userID: %v", err, event.UserID)
	}
	dbLog(fmt.Sprintf("mysql.go: addEvent(): success: results: %v", results))
	return nil
}

// deleteEvete() deletes an event based on an eventid and userid.
func deleteEvent(db *sql.DB, id, userID int64) error {
	dbLog(fmt.Sprintf("mysql.go: deleteEvent() id: %v, userID: %v", id, userID))
	// Prepare delete stmt.
	deleteDBEvent, err := db.Prepare("DELETE FROM events WHERE id=? AND userid=?")
	dbLog("mysql.go: updateEvent() db.Prepare() complete")
	// Insert demo event into db.
	dbLog(fmt.Sprintf("mysql.go: updateEvent(): id: %v, userID %v", id, userID))
	results, err := deleteDBEvent.Exec(id, userID)
	if err != nil {
		return fmt.Errorf("mysql.go: deleteEvent(): problem deleting db could be wrong userid: error: %v, userID: %v", err, userID)
	}
	dbLog(fmt.Sprintf("mysql.go: deleteEvent(): success: results: %v", results))
	return nil
}
