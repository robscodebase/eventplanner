// Copyright (c) 2018 Robert Reyna. All rights reserved.
// License BSD 3-Clause https://github.com/robscodebase/eventplanner/blob/master/LICENSE.md

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
	"strconv"
	"time"
)

// dbLogin contains the credentials and connection data
// for the mysql db.
// db is the global db variable.
var dbLogIn = "root:insecure@(mysql-event-planner:3306)/mysql"
var db *sql.DB

func main() {
	// Register the db and create db and tables.
	// dbMaker runs on a loop every ten seconds
	// up to 70 times waiting for docker-compose
	// and mysql to finish setup.
	dbMaker(db, "registerDB", "main.go: call to registerDB() from dbMaker():")
	dbMaker(db, "isDB", "main.go: call to isDB() from dbMaker():")

	// Create demo database entries.
	createDemoDB(db)

	// Activate routing handlers and serve http.
	log.Print("Listening on port 8081")
	log.Print(http.ListenAndServe(":8081", runHandlers()))
}

// dbMaker() takes a funcName either registerDB() or isDB()
// both of which are responsible for creating the user
// creditials, db and tables. To allow time for docker-compose
// and mysql to setup dbMaker() uses loops every ten seconds
// up to 70 times.
func dbMaker(db *sql.DB, funcName, message string) {
	var err error
	for retries := 0; retries < 70; retries++ {
		if funcName == "" {
			log.Fatal("no function specified: must use registerDB or isDB:")
		} else if funcName == "registerDB" {
			db, err = registerDB()
		} else {
			err = isDB(db)
		}
		if err != nil {
			dbLog(fmt.Sprintf("%v: waiting for db to be ready: retry: %v", funcName, retries))
			time.Sleep(time.Second * 10)
			if retries > 69 {
				log.Panicf("%v: could not open db: db: %v: err: %v", funcName, db, err)
			}
			sLog(fmt.Sprintf("%v: db: %v", funcName, db))
		} else {
			dbLog(fmt.Sprintf("%v: success: no of retries: %v", funcName, retries))
			retries = 71
		}
	}
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
func runHandlers() http.Handler {
	sLog("main.go: main(): runHandlers(): running handlers.")
	r := mux.NewRouter()
	r.Handle("/", http.RedirectHandler("/login", http.StatusFound))

	// Get methods.
	r.Methods("GET").Path("/logout").
		Handler(errorCheck(logoutHandler))

	r.Methods("GET").Path("/register").
		Handler(errorCheck(registerHandler))

	r.Methods("GET").Path("/login").
		Handler(errorCheck(loginHandler))

	r.Methods("GET").Path("/view-events").
		Handler(errorCheck(viewEventsHandler))

	r.Methods("GET").Path("/add-event").
		Handler(errorCheck(addEventHandler))

	r.Methods("GET").Path("/edit-event/{id:[0-9]+}").
		Handler(errorCheck(editEventHandler))

	r.Methods("GET").Path("/delete-event/{id:[0-9]+}").
		Handler(errorCheck(deleteEventHandler))

	// Post Methods
	r.Methods("POST").Path("/add-event").
		Handler(errorCheck(addEventHandler))

	r.Methods("POST").Path("/update-event/{id:[0-9]+}").
		Handler(errorCheck(updateEventHandler))

	r.Methods("POST").Path("/delete-event/{id:[0-9]+}").
		Handler(errorCheck(deleteEventHandler))

	r.Methods("POST").Path("/register").
		Handler(errorCheck(registerHandler))

	r.Methods("POST").Path("/login").
		Handler(errorCheck(loginHandler))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("/go/src/eventplanner/src/server/templates")))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))

	return r
}

// loginHandler() serves the HTML page for login.html.
func loginHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog("main.go: loginHandler()")
	var user *User
	var message string
	var err error
	// Check for an existing session.
	user, err = verifySession(db, r)
	if user != nil {
		sLog(fmt.Sprintf("main.go: loginHandler(): verifySession(): user exists redirecting to view-events: user: %v", user))
		http.Redirect(w, r, "/view-events", http.StatusFound)
	}
	sLog(fmt.Sprintf("main.go: loginHandler(): after verifySession() user should be nil: %v", user))

	if r.Method == "POST" {
		if r.FormValue("username") == "" || r.FormValue("password") == "" {
			p := &PageData{Message: "username or password cannot be blank."}
			return login.runTemplate(w, r, p)
		}
		sLog(fmt.Sprintf("main.go: loginHandler(): login attempt: %v", r.FormValue("username")))
		message, user, err = userLogin(db, w, r)
		if err != nil {
			if message == "Wrong username or password" {
				p := &PageData{Message: message}
				return login.runTemplate(w, r, p)
			}
			return &errorMessage{Error: err, Message: fmt.Sprintf("main.go: loginHandler(): userLogin(): message: %v: user: %v", message, user)}
		}
		sLog(fmt.Sprintf("main.go: loginHandler(): message: %v, user: %v", message, user))
		http.Redirect(w, r, "/view-events", http.StatusFound)
	}
	p := &PageData{Message: "Enter your username."}
	return login.runTemplate(w, r, p)
}

// registerHandler() serves the HTML page for register.html.
func registerHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog(fmt.Sprintf("main.go: registerHandler(): r: %v", r))
	var user *User
	var message string
	var err error

	// Check for an existing session.
	user, err = verifySession(db, r)
	if user != nil {
		sLog(fmt.Sprintf("main.go: registerHandler(): user exists redirecting to view-events: user: %v", user))
		http.Redirect(w, r, "/view-events", http.StatusFound)
	}

	// If the user is empty continue to registration.
	if r.Method == "POST" {
		if r.FormValue("username") == "" || r.FormValue("password") == "" {
			p := &PageData{Message: "username or password cannot be blank."}
			return register.runTemplate(w, r, p)
		}
		message, user, err = registerUser(db, w, r)
		if err != nil {
			return &errorMessage{Error: err, Message: fmt.Sprintf("main.go: registerHandler(): registerUser() message: %v, user: %v", message, user)}
		}
		if message == "userExists" {
			p := &PageData{Message: "Username already exists. Please choose something else."}
			return register.runTemplate(w, r, p)
		}
		sLog(fmt.Sprintf("main.go: registerHandler(): message: %v, user: %v", message, user))
		http.Redirect(w, r, "/view-events", http.StatusFound)
	}

	p := &PageData{Message: "We will never share your information"}
	return register.runTemplate(w, r, p)
}

// viewHandler() serves the HTML page for view-events.html.
func viewEventsHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog("main.go: viewEventsHandler()")
	var user *User
	var events []*Event
	var err error
	// Check for an existing session.
	user, err = verifySession(db, r)
	if err != nil {
		log.Printf("main.go: viewEventsHandler(): error: %v: redirecting to login page", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	events, err = listEvents(db, user.Username)
	if err != nil {
		p := &PageData{PageName: "View Events", Message: fmt.Sprintf("No Events to view: %v", err)}
		return viewEvents.runTemplate(w, r, p)
	}
	log.Printf("main.go: viewEventsHandler(): events from listEvents: %v", events)
	log.Printf("main.go: viewEventsHandler(): user: %v", user)
	p := &PageData{Events: events, PageName: "View Events"}
	sLog("main.go: main():  viewEventsHandler() call to handler.")
	return viewEvents.runTemplate(w, r, p)
}

// addEventHandler() serves the HTML page for add-event.html.
func addEventHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog("main.go: main():  addEventHandler()")
	var err error
	var user *User
	user, err = verifySession(db, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	if r.Method == "POST" {
		event := &Event{
			Name:        r.FormValue("name"),
			StartTime:   r.FormValue("startTime"),
			EndTime:     r.FormValue("endTime"),
			Description: r.FormValue("description"),
			UserID:      user.ID,
		}
		err = addEventDB(db, event)
		if err != nil {
			sLog(fmt.Sprintf("main.go: main():  addEventHandler(): addEvent(): err: %v: redirecting to add-event page:", err))
			http.Redirect(w, r, "/add-event", http.StatusFound)
			return nil
		}
		http.Redirect(w, r, "/view-events", http.StatusFound)
		return nil
	}
	p := &PageData{PageName: "Add Event"}
	sLog("main.go: main():  addEventsHandler(): call to handler")
	return addEvent.runTemplate(w, r, p)
}

// editEventHandler() serves the HTML page for edit-event.html.
func editEventHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog("main.go: main():  editEventHandler()")
	var err error
	var user *User
	// Check for an existing session.
	user, err = verifySession(db, r)
	if err != nil {
		sLog(fmt.Sprintf("main.go: editEventHandler(): error: %v: redirecting to login page:", err))
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	sLog("main.go: editEventHandler(): db.QueryRow: event found:")
	eventID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		sLog(fmt.Sprintf("main.go: main():  editEventHandler(): strconv.ParseInt(): err: %v: redirecting to view-events page:", err))
		http.Redirect(w, r, "/view-events", http.StatusFound)
		return nil
	}
	sLog(fmt.Sprintf("main.go: editEventHandler(): user: %v: eventID: %v", user, eventID))
	event, err := listEvent(db, user.ID, eventID)
	if err != nil {
		sLog(fmt.Sprintf("main.go: editEventHandler(): call to listEvent(): error: %v: redirecting to view-events page:", err))
		http.Redirect(w, r, "/view-events", http.StatusFound)
		return nil
	}
	p := &PageData{PageName: "Edit Event", Event: event}
	return editEvent.runTemplate(w, r, p)
}

// updateEventHandler() updates an event if the userid and eventid match the db.
func updateEventHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog("main.go: main():  updateEventHandler():")
	var err error
	var user *User
	// Check for an existing session.
	user, err = verifySession(db, r)
	if err != nil {
		sLog(fmt.Sprintf("main.go: updateEventHandler(): error: %v: redirecting to login page:", err))
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	eventID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		sLog(fmt.Sprintf("main.go: main():  updateEventHandler(): strconv.ParseInt(): err: %v: redirecting to view-events page:", err))
		http.Redirect(w, r, "/view-events", http.StatusFound)
		return nil
	}
	event := &Event{
		ID:          eventID,
		Name:        r.FormValue("name"),
		StartTime:   r.FormValue("startTime"),
		EndTime:     r.FormValue("endTime"),
		Description: r.FormValue("description"),
		UserID:      user.ID,
	}
	err = updateEvent(db, event)
	if err != nil {
		sLog(fmt.Sprintf("main.go: main():  updateEventHandler(): updateEvent(): err: %v: redirecting to view-events page:", err))
		http.Redirect(w, r, "/view-events", http.StatusFound)
		return nil
	}
	http.Redirect(w, r, "/edit-event/"+mux.Vars(r)["id"], http.StatusFound)
	return nil
}

// deleteEventHandler() deletes an event if the userid and eventid matches.
func deleteEventHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	sLog("main.go: main():  deleteEventHandler():")
	var err error
	var user *User
	// Check for an existing session.
	user, err = verifySession(db, r)
	if err != nil {
		sLog(fmt.Sprintf("main.go: deleteEventHandler(): error: %v: redirecting to login page:", err))
		http.Redirect(w, r, "/login", http.StatusFound)
		return nil
	}
	eventID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		sLog(fmt.Sprintf("main.go: main():  deleteEventHandler(): strconv.ParseInt(): err: %v: redirecting to view-events page:", err))
		http.Redirect(w, r, "/view-events", http.StatusFound)
		return nil
	}
	err = deleteEvent(db, eventID, user.ID)
	if err != nil {
		sLog(fmt.Sprintf("main.go: main():  updateEventHandler(): deleteEvent(): err: %v: redirecting to view-events page:", err))
		http.Redirect(w, r, "/view-events", http.StatusFound)
		return nil
	}
	http.Redirect(w, r, "/view-events", http.StatusFound)
	return nil
}

// logoutHandler() sets the cookie to die immediately.
func logoutHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	cookie := &http.Cookie{
		Name:   "golang-event-planner",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusFound)
	return nil
}

// ServeHTTP() ensures there are no errors before serving the HTML data.
func (errCheck errorCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if errcheck := errCheck(w, r); errcheck != nil {
		log.Printf("main.go: ServeHTTP(): error: status code: %d, message: %s, error: %#v",
			errcheck.Code, errcheck.Message, errcheck.Error)
		http.Error(w, errcheck.Message, errcheck.Code)
	}
}
