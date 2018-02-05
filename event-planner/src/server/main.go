package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"html/template"
    "io/ioutil"
)
var Login bool

func main() {
	runHandlers()
}
var (
    homePage  = compileTemplate("home.html")
)

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
    return homePage.runTemplate(w, r, nil)
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

type eventPlannerTemplate struct {
    eventPlannerTemplate *template.Template
}

func compileTemplate(templateName string) *eventPlannerTemplate {
    // Add the main template file.
    main := template.Must(template.ParseFiles("/go/src/event-planner/src/server/templates/main.html"))

    // Add the header.
    header := readFile("header.html")
    // Add the body.
    body := readFile(templateName)
    // Add the footer.
    footer := readFile("footer.html")

    // Combine header, body, and footer and return template.
    template.Must(main.New("header").Parse(string(header)))
    template.Must(main.New("body").Parse(string(body)))
    template.Must(main.New("footer").Parse(string(footer)))
    return &eventPlannerTemplate{main.Lookup("main.html")}
}

func readFile(fileName string) []byte {
    template, err := ioutil.ReadFile("/go/src/event-planner/src/server/templates/" + fileName)
    if err != nil {
        panic(fmt.Errorf("compileTemplate() could not read footer: %v", err))
    }
    return template
}
func (template *eventPlannerTemplate) runTemplate(w http.ResponseWriter, r *http.Request, input interface{}) *errorMessage {
    session := struct {
        Input        interface{}
        AuthEnabled bool
        UserData     *UserData
    }{
        Input:        input,
        AuthEnabled: true,
    }

    if Login {
         session.UserData = sessionData(r)
    }

    if err := template.eventPlannerTemplate.Execute(w, session); err != nil {
        return formatError(err, "runTemplate() could not execute template: %v")
    }
    return nil
}

func sessionData(r *http.Request) *UserData {
    return &UserData{
        ID:          "SECURE ID",
        UserName: "ADMIN",
    }
}

type UserData struct {
    ID, UserName string
}
