package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func (db *database) registerDB() error {
	dbLog(fmt.Sprintf("mysql.go: registerDB():"))
	db, err = sql.Open("mysql", "mysqldb:insecure@(127.17.0.2:3306)/")
	db.db = &sql.DB{db}
	dbLog(fmt.Sprintf("mysql.go: registerDB(): sql.Open: %v", db.db))
	if err != nil {
		return fmt.Errorf("mysql.go: registerDB(): sqlOpen err: %v", err)
	}
	return nil
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

func isTable() {}
func createDB(db *sql.DB) error {
	dbLog(fmt.Sprintf("mysql.go: createDB() var db: %v", db))
	for _, sqlCommand := range tableCreator {
		dbLog(fmt.Sprintf("mysql.go: createDB(): inside for loop sqlCommand: %v", sqlCommand))
		result, err := db.Exec(`SELECT User FROM mysql.user`)
		dbLog(fmt.Sprintf("mysql.go: createDB(): inside for loop db.Exec: %v: result: %v", sqlCommand, result))
		if err != nil {
			return fmt.Errorf("mysql.go: createDB(): db.Exec: problem with command: %v: error: %v", sqlCommand, err)
		}
	}
	return nil
}
