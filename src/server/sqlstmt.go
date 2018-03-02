package main

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
