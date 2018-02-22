package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

func registerDB() (*sql.DB, error) {
	dbLog(fmt.Sprintf("mysql.go: registerDB()"))
	db, err = sql.Open("mysql", dbLogIn)
	dbLog(fmt.Sprintf("mysql.go: registerDB(): sql.Open: %v", db))
	if err != nil {
		return db, fmt.Errorf("mysql.go: registerDB(): sql.Open db: %v: err: %v", db, err)
	}
	defer db.Close()
	return db, nil
}

var tableCreator = []string{
	`CREATE DATABASE IF NOT EXISTS eventplanner DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE eventplanner;`,
	`CREATE TABLE IF NOT EXISTS events (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		name VARCHAR(255) NULL,
		starttime VARCHAR(255) NULL,
		endtime VARCHAR(255) NULL,
		description VARCHAR(255) NULL,
		createdby VARCHAR(255) NULL,
		PRIMARY KEY (ID)
	)`,
}

func isDB(db *sql.DB) error {
	db, err := sql.Open("mysql", dbLogIn)
	if err != nil {
		return fmt.Errorf("mysql.go: isDB(): sql.Open db: %v: error: %v", db, err)
	}
	defer db.Close()

	// Ping the database.
	if db.Ping() == driver.ErrBadConn {
		return fmt.Errorf("mysql.go: isDB() db.Ping() error: could not ping database.")
	}

	// Try to use the database if there is an error 1049 create the database.
	if _, err := db.Exec("USE eventplanner"); err != nil {
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {
			return createDB(db)
		}
	}

	// Try to describe the table if there is an error 1146 create the table.
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
	dbLog(fmt.Sprintf("mysql.go: createDB() var db: %v", db))
	for _, sqlCommand := range tableCreator {
		dbLog(fmt.Sprintf("mysql.go: createDB(): inside for loop before db.Exec sqlCommand: %v", sqlCommand))
		result, err := db.Exec(sqlCommand)
		dbLog(fmt.Sprintf("mysql.go: createDB(): inside for loop after db.Exec: %v: result: %v", sqlCommand, result))
		if err != nil {
			return fmt.Errorf("mysql.go: createDB(): db.Exec: problem with command: %v: error: %v", sqlCommand, err)
		}
	}
	return nil
}
