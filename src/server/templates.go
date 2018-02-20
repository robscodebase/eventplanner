package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't open working directory. err %v:", err)
	}
	filePathBase = filepath.Base(workingDirectory)
}

type eventPlannerTemplate struct {
	eventPlannerTemplate *template.Template
}

func compileTemplate(templateName string) *eventPlannerTemplate {
	if templateName == "login" {
		login := template.Must(template.ParseFiles(filePathBase + "templates/login.html"))
		return &eventPlannerTemplate{login.Lookup("login.html")}
	}
	if templateName == "register" {
		register := template.Must(template.ParseFiles(filePathBase + "templates/register.html"))
		return &eventPlannerTemplate{register.Lookup("register.html")}
	}
	// Add the main template file.
	main := template.Must(template.ParseFiles(filePathBase + "templates/main.html"))

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
	template, err := ioutil.ReadFile(filePathBase + "templates/" + fileName)
	if err != nil {
		panic(fmt.Errorf("compileTemplate() could not read footer: %v", err))
	}
	return template
}

func (template *eventPlannerTemplate) runTemplate(w http.ResponseWriter, r *http.Request, input interface{}) *errorMessage {
	session := struct {
		Input       interface{}
		AuthEnabled bool
		UserData    *UserData
	}{
		Input:       input,
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
