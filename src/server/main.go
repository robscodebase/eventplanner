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
	// For loop tries every 10 seconds 6 times before failure.
	// Check if the database exists isDB().
	for retries := 0; retries < 70; retries++ {
		db, err = registerDB()
		if err != nil {
			dbLog(fmt.Sprintf("main.go: call to registerDB(): waiting for db to be ready: retry: %v", retries))
			time.Sleep(time.Second * 10)
			if retries > 69 {
				log.Panicf("main.go: call to registerDB(): could not open db: db: %v: err: %v", db, err)
			}
			sLog(fmt.Sprintf("main.go: main(): registerDB(): db: %v", db))
		} else {
			dbLog(fmt.Sprintf("main.go: call to registerDB(): success: no of retries: %v", retries))
			retries = 71
		}

	}

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

	r.Methods("POST").Path("/login").
		Handler(errorCheck(loginHandler))

	r.Methods("GET").Path("/logout").
		Handler(errorCheck(logoutHandler))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("/go/src/eventplanner/src/server/templates")))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	log.Print("Listening on port 8081")
	log.Print(http.ListenAndServe(":8081", r))
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
	p := &PageData{PageName: "View Events"}
	sLog("main.go: main(): runHandlers(): viewEventsHandler() call to handler.")
	return viewEvents.runTemplate(w, r, p)
}

// addEventHandler() serves the HTML page for add-event.html.
func addEventHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	var user *User
	user, err = verifySession(db, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	log.Println("user", user)
	p := &PageData{PageName: "Add Event"}
	sLog("main.go: main(): runHandlers(): addEventsHandler(). call to handler")
	return addEvent.runTemplate(w, r, p)
}

// editEventHandler() serves the HTML page for edit-event.html.
func editEventHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	var user *User
	user, err = verifySession(db, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	log.Println("user", user)
	p := &PageData{PageName: "Edit Event"}
	sLog("main.go: main(): runHandlers(): editEventsHandler(). call to handler")
	return editEvent.runTemplate(w, r, p)
}

// logoutHandler sets the cookie to die immediately and
// clears memory of all username and session data.
func logoutHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	cookie := &http.Cookie{
		Name:   "golang-event-planner",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusFound)
	return nil
}

// ServeHTTP ensures there are no errors before serving the HTML data.
func (errCheck errorCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if errcheck := errCheck(w, r); errcheck != nil {
		log.Printf("main.go: ServeHTTP(): error: status code: %d, message: %s, error: %#v",
			errcheck.Code, errcheck.Message, errcheck.Error)
		http.Error(w, errcheck.Message, errcheck.Code)
	}
}
