package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var Login bool
var filePathBase string

func main() {
	runHandlers()
}

var (
	viewEvents = compileTemplate("view-events.html")
	addEvents  = compileTemplate("view-events.html")
	editEvents = compileTemplate("view-events.html")
	login      = compileTemplate("view-events.html")
)

func runHandlers() {
	log.Println("main.go: main(): runHandlers(): running handlers.")
	r := mux.NewRouter()
	r.Handle("/", http.RedirectHandler("/home", http.StatusFound))
	r.Methods("GET").Path("/login").
		Handler(errorCheck(viewEventsHandler))
	r.Methods("GET").Path("/view-events").
		Handler(errorCheck(viewEventsHandler))
	r.Methods("GET").Path("/add-events").
		Handler(errorCheck(viewEventsHandler))
	r.Methods("GET").Path("/edit-events").
		Handler(errorCheck(viewEventsHandler))
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	log.Print("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func loginHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	log.Println("main.go: main(): runHandlers(): loginHandler().")
	return viewEvents.runTemplate(w, r, nil)
}
func viewEventsHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	log.Println("main.go: main(): runHandlers(): viewEventsHandler().")
	return viewEvents.runTemplate(w, r, nil)
}
func addEventsHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	log.Println("main.go: main(): runHandlers(): addEventsHandler().")
	return addEvents.runTemplate(w, r, nil)
}
func editEventsHandler(w http.ResponseWriter, r *http.Request) *errorMessage {
	log.Println("main.go: main(): runHandlers(): editEventsHandler().")
	return editEvents.runTemplate(w, r, nil)
}

type errorCheck func(http.ResponseWriter, *http.Request) *errorMessage

type errorMessage struct {
	Error   error
	Message string
	Code    int
}

func (errCheck errorCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if errcheck := errCheck(w, r); errcheck != nil {
		log.Printf("serveHTTP error: status code: %d, message: %s, error: %#v",
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

func sessionData(r *http.Request) *UserData {
	return &UserData{
		ID:       "SECURE ID",
		UserName: "ADMIN",
	}
}

type UserData struct {
	ID, UserName string
}
