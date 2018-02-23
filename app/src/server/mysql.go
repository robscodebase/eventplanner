package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

func registerDB() (*sql.DB, error) {
	dbLog(fmt.Sprintf("mysql.go: registerDB()"))
	db, err = sql.Open("mysql", dbLogIn)
	dbLog(fmt.Sprintf("mysql.go: registerDB(): sql.Open: %v", db))
	if err != nil {
		return db, fmt.Errorf("mysql.go: registerDB(): sql.Open db: %v: err: %v", db, err)
	}
	return db, nil
}

var tableCreator = []string{
	`CREATE DATABASE IF NOT EXISTS eventplanner;`,
	`USE eventplanner;`,
	`CREATE TABLE IF NOT EXISTS events (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		name VARCHAR(255) NULL,
		starttime VARCHAR(255) NULL,
		endtime VARCHAR(255) NULL,
		description VARCHAR(255) NULL,
		createdby VARCHAR(255) NULL,
		PRIMARY KEY (ID)
	);`,
}

func isDB(db *sql.DB) error {
	dbLog(fmt.Sprintf("mysql.go: isDB()"))
	dbLog(fmt.Sprintf("mysql.go: isDB(): open db"))
	//db, err := sql.Open("mysql", dbLogIn)
	//if err != nil {
	//return fmt.Errorf("mysql.go: isDB(): sql.Open db: %v: error: %v", db, err)
	//}
	//defer db.Close()

	// Ping the database with db.Ping().
	dbLog(fmt.Sprintf("mysql.go: isDB(): ping db"))
	if db.Ping() == driver.ErrBadConn {
		return fmt.Errorf("mysql.go: isDB() db.Ping() error: could not ping database.")
	}

	// Try to use the database if there is an error 1049 create the database.
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
	return nil
}

func createDB(db *sql.DB) error {
	dbLog(fmt.Sprintf("mysql.go: createDB(): var db: %v", db))
	for _, sqlCommand := range tableCreator {
		dbLog(fmt.Sprintf("mysql.go: createDB(): inside for loop before db.Exec sqlCommand: %v", sqlCommand))
		result, err := db.Exec(sqlCommand)
		dbLog(fmt.Sprintf("mysql.go: createDB(): inside for loop after result: %v", result))
		if err != nil {
			return fmt.Errorf("mysql.go: createDB(): db.Exec: problem exec command: error: %v", err)
		}
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanEvent(scanRow rowScanner) (*Event, error) {
	var (
		id          int64
		name        sql.NullString
		starttime   sql.NullString
		endtime     sql.NullString
		description sql.NullString
		createdby   sql.NullString
	)

	if err := scanRow.Scan(&id, &name, &starttime, &endtime, &description, &createdby); err != nil {
		return nil, fmt.Errorf("mysql.go: scanEvent(): scanRow.Scan(): error: %v", err)
	}

	event := &Event{
		ID:          id,
		Name:        name.String,
		StartTime:   starttime.String,
		EndTime:     endtime.String,
		Description: description.String,
		CreatedBy:   createdby.String,
	}
	return event, nil
}

func viewDBEvents(db *sql.DB) { // ([]*Event, error) {
	dbLog(fmt.Sprintf("mysql.go: viewDBEvents(): var db: %v", db))
	rows, err := db.Query(`SELECT * FROM events`)
	if err != nil {
		log.Println(err)
		//return nil, fmt.Errorf("mysql.go: viewEvents(): rows = %v: error: %v", rows, err)
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			log.Println(err)
			//return nil, fmt.Errorf("mysql.go: viewDBEvents(): scanEvent(): printing row: error: %v", err)
		}

		events = append(events, event)
	}
	dbLog(fmt.Sprintf("mysql.go: viewDBEvents(): events from rows: %v", events))
}

func addDBEvent(db *sql.DB) {
	event := &Event{
		ID:          int64(1234),
		Name:        "First Event",
		StartTime:   "monday",
		EndTime:     "tuesday",
		Description: "description 1",
		CreatedBy:   "createdby1",
	}

	rowResult, err := writeResult(db.insert, event.ID, event.Name, event.StartTime, event.EndTime, event.Description, event.CreatedBy)
	if err != nil {
		fmt.Println(err)
	}

	dbLog(fmt.Sprintf("mysql.go: addEvent(): add finished: %v", rowResult))
}

func writeResult(sqlstmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	result, err := sqlstmt.Exec(args...)
	if err != nil {
		return result, fmt.Errorf("mysql.go: writeResult(): could not execute statement: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return result, fmt.Errorf("mysql.go: writeResult(): could not get rowsAffected: error: %v", err)
	} else if rowsAffected != 1 {
		return result, fmt.Errorf("mysql.go: writeResult(): expected 1 row affected, got %d", rowsAffected)
	}
	return result, nil
}
