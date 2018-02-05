package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)

type Login bool

func main() {
    registerHandlers()
}

func registerHandlers() {
    r := mux.NewRouter()
    r.Handle("/", http.RedirectHandler("/event-planner", http.StatusFound))
    r.Methods("GET").Path("/").
        Handler(errorCheck(home))
    http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
}

func home(w http.ResponseWriter, r *http.Request) *errorMessage {
    books, err := bookshelf.DB.ListBooks()
    if err != nil {
        return formatError(err, "could not list books: %v", err)
    }
    return compileTemplate.Execute(w, r, books)
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
