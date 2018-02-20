package main

import (
	"net/http"
	"time"
)

type PageData struct {
	Event
	User
	PageName string
}

type Event struct {
	Name        string
	StartTime   time.Time
	EndTime     time.Time
	Description string
}

type User struct {
	ID        string
	UserName  string
	FirstName string
	LastName  string
	Email     string
}

type errorCheck func(http.ResponseWriter, *http.Request) *errorMessage

type errorMessage struct {
	Error   error
	Message string
	Code    int
}
