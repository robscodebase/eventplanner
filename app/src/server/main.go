package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var Login bool
var err error
var filePathBase string

type database struct {
	db *sql.DB
}

func main() {
	var db database
	err = db.registerDB()
	if err != nil {
		log.Println("main.go: main(): call to registerDB(): error: ", err)
	}
	runHandlers()
}

var (
	viewEvents = compileTemplate("view-events.html")
	addEvent   = compileTemplate("add-event.html")
	editEvent  = compileTemplate("edit-event.html")
	login      = compileTemplate("login.html")
	register   = compileTemplate("register.html")
)

func runHandlers() {
	sLog("main.go: main(): runHandlers(): running handlers.")
	r := mux.NewRouter()
	r.Handle("/", http.RedirectHandler("/home", http.StatusFound))

	r.Methods("GET").Path("/register").
		Handler(errorCheck(registerHandler))

	r.Methods("GET").Path("/login").
		Handler(errorCheck(loginHandler))

	r.Methods("GET").Path("/view-events").
		Handler(errorCheck(viewEventsHandler))

	r.Methods("GET").Path("/add-event").
		Handler(errorCheck(addEventHandler))

	r.Methods("GET").Path("/edit-event").
		Handler(errorCheck(editEventHandler))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	log.Print("Listening on port 8081")
	log.Print(http.ListenAndServe(":8081", r))
}

func registerHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog("main.go: main(): runHandlers(): registerHandler(): call to handler.")
	return register.runTemplate(w, r, nil)
}
func loginHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog("main.go: main(): runHandlers(): loginHandler() call to handler.")
	return login.runTemplate(w, r, nil)
}
func viewEventsHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	p := &PageData{PageName: "View Events"}
	sLog("main.go: main(): runHandlers(): viewEventsHandler() call to handler.")
	return viewEvents.runTemplate(w, r, p)
}
func addEventHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	p := &PageData{PageName: "Add Event"}
	sLog("main.go: main(): runHandlers(): addEventsHandler(). call to handler")
	return addEvent.runTemplate(w, r, p)
}
func editEventHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	p := &PageData{PageName: "Edit Event"}
	sLog("main.go: main(): runHandlers(): editEventsHandler(). call to handler")
	return editEvent.runTemplate(w, r, p)
}

func (errCheck errorCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if errcheck := errCheck(w, r); errcheck != nil {
		log.Printf("main.go: ServeHTTP(): error: status code: %d, message: %s, error: %#v",
			errcheck.Code, errcheck.Message, errcheck.Error)
		http.Error(w, errcheck.Message, errcheck.Code)
	}
}

func formatError(err error, format string, v ...interface{}) *errorMessage {
	return &errorMessage{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}

func sessionData(r *http.Request) *User {
	return &User{
		ID:       "SECURE ID",
		UserName: "ADMIN",
	}
}
