// mysql.go contains the functions for interacting with the database.
package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
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
     createdby VARCHAR(255) NULL,
     PRIMARY KEY (ID)
     );`,
	`CREATE TABLE IF NOT EXISTS users (
	 username VARCHAR(255),
	 secret BINARY(255),
	 cookieSession VARCHAR(255),
	 PRIMARY KEY (username)
 	 );`,
}

// registerDB() opens the db and returns the db instance.
func registerDB() (*sql.DB, error) {
	dbLog(fmt.Sprintf("mysql.go: registerDB()"))
	// retries give time for docker-compose and mysql to finish setup.
	for retries := 0; retries < 70; retries++ {
		db, err = sql.Open("mysql", dbLogIn)
		if err != nil {
			dbLog(fmt.Sprintf("mysql.go: registerDB(): sql.Open error: retry count: %v", retries))
			time.Sleep(time.Second * 10)
			if retries > 69 {
				return db, fmt.Errorf("mysql.go: registerDB(): sql.Open db: %v: err: %v", db, err)
			}
		}
	}
	dbLog(fmt.Sprintf("mysql.go: registerDB(): sql.Open success: db: %v", db))
	return db, nil
}

// isDB() pings, and checks that the db and table have been created.
func isDB(db *sql.DB) error {
	dbLog(fmt.Sprintf("mysql.go: isDB()"))
	dbLog(fmt.Sprintf("mysql.go: isDB(): open db"))
	//db, err := sql.Open("mysql", dbLogIn)
	//if err != nil {
	//return fmt.Errorf("mysql.go: isDB(): sql.Open db: %v: error: %v", db, err)
	//}
	//defer db.Close()

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

// viewDBEvents()_ returns events assigned to the current user.
func viewDBEvents(db *sql.DB) {
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

// createDemoDB() add the demo events and user to the db.
func createDemoDB(db *sql.DB) {
	dbLog(fmt.Sprintf("mysql.go: createDemoDB(): var db: %v", db))
	// Range over each of the demoEvents and insert them in to the db.
	for _, demo := range demoEvents {
		// Prepare stmt.
		insertDemoEvent, err := db.Prepare("INSERT INTO events (id, name, starttime, endtime, description, createdby) VALUES(?, ?, ?, ?, ?, ?)")
		dbLog(fmt.Sprintf("mysql.go: createDemoDB(): insertDemoEvent: %v", insertDemoEvent))
		// Insert demo event into db.
		results, err := insertDemoEvent.Exec(demo.ID, demo.Name, demo.StartTime, demo.EndTime, demo.Description, demo.Description)
		if err != nil {
			dbLog(fmt.Sprintf("mysql.go: createDemoDB(): problem creating demo db most likely entries already exist: %v", err))
			return
		}
		dbLog(fmt.Sprintf("mysql.go: createDemoDB(): insertDemoEvent success: results: %v", results))
	}

	// Prepare insert stmt.
	insertDemoUser, err := db.Prepare("INSERT INTO users (username, secret) VALUES(?, ?)")
	dbLog(fmt.Sprintf("mysql.go: createDemoDB(): insertDemoUser: %v", insertDemoUser))
	// Encrypt password.
	secret, err := bcrypt.GenerateFromPassword(demoUser.Secret, bcrypt.DefaultCost)
	if err != nil {
		log.Panicf("mysql.go: createDemoDB(): error encrypting demo password secret: %v: error: %v", secret, err)
	}
	sLog(fmt.Sprintf("auth.go: registerUser(): encrypt password: %v", secret))
	// Insert user into db.
	results, err := insertDemoUser.Exec(demoUser.Username, secret)
	if err != nil {
		dbLog(fmt.Sprintf("mysql.go: createDemoDB(): problem creating demo user most likely entries already exist: %v", err))
		return
	}
	dbLog(fmt.Sprintf("mysql.go: createDemoDB(): var results: %v", results))
}
