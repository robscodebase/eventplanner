package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func (db database) registerDB() error {
	dbLog("mysql.go: registerDB()")
	db.db, err = sql.Open("mysql", "mysqldb:insecure@(ipaddress:port)/mysqldb")
	if err != nil {
		return fmt.Errorf("mysql.go: registerDB(): sqlOpen err: %v", err)
	}
	return nil
}
