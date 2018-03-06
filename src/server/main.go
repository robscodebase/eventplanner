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
	"time"
)

var Login bool
var err error
var filePathBase string
var dbLogIn = "root:insecure@(mysql-event-planner:3306)/mysql"

var db *sql.DB

func main() {
	// Register credentials for database with registerDB().
	db, err = registerDB()
	if err != nil {
		log.Panicf("main.go: main(): call to registerDB(): error: %v", err)
	}
	sLog(fmt.Sprintf("main.go: main(): db.register(): db: %v", db))

	// For loop tries every 10 seconds 6 times before failure.
	// Check if the database exists isDB().
	for retries := 0; retries < 70; retries++ {
		err = isDB(db)
		if err != nil {
			dbLog(fmt.Sprintf("main.go: call to isDB(): waiting for db to be ready: retry: %v", retries))
			time.Sleep(time.Second * 10)
			if retries > 69 {
				log.Panicf("main.go: call to isDB(): could not open db: db: %v: err: %v", db, err)
			}
		} else {
			dbLog(fmt.Sprintf("main.go: call to isDB(): success: no of retries: %v", retries))
			retries = 71
		}

	}

	// Try to add events.
	//addDBEvent(db)
	// Try to view events.
	viewDBEvents(db)

	// Create demo database entries.
	createDemoDB(db)

	// Activate routing handlers with runHandlers()
	runHandlers()
}

// Template page variables viewEvents, addEvent, editEvent
// login and register link the html module for that pages
// body and returns a complete html page with header and footer.
var (
	viewEvents = compileTemplate("view-events.html")
	addEvent   = compileTemplate("add-event.html")
	editEvent  = compileTemplate("edit-event.html")
	login      = compileTemplate("login.html")
	register   = compileTemplate("register.html")
)

// runHandlers() activates routing handlers for each page
// and actions completed on each page and form.
func runHandlers() {
	sLog("main.go: main(): runHandlers(): running handlers.")
	r := mux.NewRouter()
	r.Handle("/", http.RedirectHandler("/edit-events", http.StatusFound))

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

	r.Methods("POST").Path("/register").
		Handler(errorCheck(registerHandler))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("/go/src/eventplanner/src/server/templates")))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	log.Print("Listening on port 8081")
	log.Print(http.ListenAndServe(":8081", r))
}

// loginHandler() serves the HTML page for login.html.
func loginHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog("main.go: main(): runHandlers(): loginHandler() call to handler.")
	return login.runTemplate(w, r, nil)
}

// registerHandler() serves the HTML page for register.html.
func registerHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog(fmt.Sprintf("main.go: registerHandler(): r: %v", r))
	var user *User
	var message string
	var err error

	// Check for an existing session.
	user, err = verifySession(db, r)
	if err != nil {
	}

	// If the user is empty continue to registration.
	if user == (&User{}) {
		if r.Method == "POST" {
			message, user, err = registerUser(db, w, r)
			if err != nil {
				return &errorMessage{Error: err, Message: fmt.Sprintf("main.go: registerHandler(): registerUser() message: %v, user: %v", message, user)}
			}
			sLog(fmt.Sprintf("main.go: registerHandler(): message: %v, user: %v", message, user))
		}
		return register.runTemplate(w, r, nil)
	}

	// If the user is not empty redirect to the view-events page.
	http.Redirect(w, r, "/view-events", http.StatusFound)
	return nil
}

// viewHandler() serves the HTML page for view-events.html.
func viewEventsHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	p := &PageData{PageName: "View Events"}
	sLog("main.go: main(): runHandlers(): viewEventsHandler() call to handler.")
	return viewEvents.runTemplate(w, r, p)
}

// addEventHandler() serves the HTML page for add-event.html.
func addEventHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	p := &PageData{PageName: "Add Event"}
	sLog("main.go: main(): runHandlers(): addEventsHandler(). call to handler")
	return addEvent.runTemplate(w, r, p)
}

// editEventHandler() serves the HTML page for edit-event.html.
func editEventHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	p := &PageData{PageName: "Edit Event"}
	sLog("main.go: main(): runHandlers(): editEventsHandler(). call to handler")
	return editEvent.runTemplate(w, r, p)
}

// ServeHTTP ensures there are no errors before serving the HTML data.
func (errCheck errorCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if errcheck := errCheck(w, r); errcheck != nil {
		log.Printf("main.go: ServeHTTP(): error: status code: %d, message: %s, error: %#v",
			errcheck.Code, errcheck.Message, errcheck.Error)
		http.Error(w, errcheck.Message, errcheck.Code)
	}
}
