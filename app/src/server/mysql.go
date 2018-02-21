package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func (db *sql.DB) registerDB() error {
	dbLog("mysql.go: registerDB()")
	db, err = sql.Open("mysql", "mysqldb:insecure@(ipaddress:port)/mysqldb")
	if err != nil {
		return fmt.ErrorF("mysql.go: registerDB(): sqlOpen err: %v", err)
	}
}
