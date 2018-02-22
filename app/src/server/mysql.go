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

func viewDBEvents(db *sql.DB) { // ([]*Event, error) {
	dbLog(fmt.Sprintf("mysql.go: viewDBEvents(): var db: %v", db))
	rows, err := db.Query(`SELECT * FROM events`)
	if err != nil {
		log.Println(err)
		//return nil, fmt.Errorf("mysql.go: viewEvents(): rows = %v: error: %v", rows, err)
	}
	defer rows.Close()

	for rows.Next() {
		log.Println("rows", rows)
	}

	dbLog(fmt.Sprintf("mysql.go: viewDBEvents(): var rows: %v", rows))
}
