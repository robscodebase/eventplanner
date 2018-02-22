package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func registerDB() (*sql.DB, error) {
	dbLog(fmt.Sprintf("mysql.go: registerDB()"))
	db, err = sql.Open("mysql", "root:insecure@(mysql:3306)/mysql")
	dbLog(fmt.Sprintf("mysql.go: registerDB(): sql.Open: %v", db))
	if err != nil {
		return db, fmt.Errorf("mysql.go: registerDB(): sql.Open db: %v: err: %v", db, err)
	}
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

func createDB(db *sql.DB) error {
	dbLog(fmt.Sprintf("mysql.go: createDB() var db: %v", db))
	for _, sqlCommand := range tableCreator {
		dbLog(fmt.Sprintf("mysql.go: createDB(): inside for loop sqlCommand: %v", sqlCommand))

		result, err := db.Exec(`SHOW DATABASES;`)

		dbLog(fmt.Sprintf("mysql.go: createDB(): inside for loop db.Exec: %v: result: %v", sqlCommand, result))
		if err != nil {
			return fmt.Errorf("mysql.go: createDB(): db.Exec: problem with command: %v: error: %v", sqlCommand, err)
		}
	}
	return nil
}
