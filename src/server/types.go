// types.go contains most of the structs and types encountered in the program including Event, User, and Page data.
package main

import (
	"net/http"
)

type PageData struct {
	Events []*Event
	User
	PageName string
	Message  string
}

type Event struct {
	ID          int64
	Name        string
	StartTime   string
	EndTime     string
	Description string
	UserID      int64
}

type User struct {
	ID            int64
	Username      string
	Secret        []byte
	CookieSession string
}

type errorCheck func(http.ResponseWriter, *http.Request) *errorMessage

type errorMessage struct {
	Error   error
	Message string
	Code    int
}
