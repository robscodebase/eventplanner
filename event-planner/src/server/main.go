package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
)

func main() {
	runHandlers()
}

func runHandlers() {
	r := mux.NewRouter()
	r.Handle("/", http.RedirectHandler("/home", http.StatusFound))
	r.Methods("GET").Path("/home").
		Handler(errorCheck(home))
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func home(w http.ResponseWriter, r *http.Request) *errorMessage {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
    return nil
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
