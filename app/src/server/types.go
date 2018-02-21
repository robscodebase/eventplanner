package main

import (
	"net/http"
)

type PageData struct {
	Event
	User
	PageName string
}

type Event struct {
	ID          int64
	Name        string
	StartTime   string
	EndTime     string
	Description string
	CreatedBy   string
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
