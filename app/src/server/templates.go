package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

var filePath = "/go/src/eventplanner/src/server"

type eventPlannerTemplate struct {
	eventPlannerTemplate *template.Template
}

func compileTemplate(templateName string) *eventPlannerTemplate {
	if templateName == "login.html" {
		login := template.Must(template.ParseFiles(filePath + "/templates/" + templateName))
		return &eventPlannerTemplate{login.Lookup(templateName)}
	}
	if templateName == "register.html" {
		register := template.Must(template.ParseFiles(filePath + "/templates/" + templateName))
		return &eventPlannerTemplate{register.Lookup(templateName)}
	}
	// Add the main template file.
	main := template.Must(template.ParseFiles(filePath + "/templates/main.html"))

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
	template, err := ioutil.ReadFile(filePath + "/templates/" + fileName)
	if err != nil {
		panic(fmt.Errorf("templates.go: readFile(): could not read file: %v: %v", fileName, err))
	}
	return template
}

func (template *eventPlannerTemplate) runTemplate(w http.ResponseWriter, r *http.Request, input interface{}) *errorMessage {
	session := struct {
		Input       interface{}
		AuthEnabled bool
		User        *User
	}{
		Input:       input,
		AuthEnabled: true,
	}

	if Login {
		session.User = sessionData(r)
	}

	if err := template.eventPlannerTemplate.Execute(w, session); err != nil {
		return formatError(err, "templates.go: runTemplate(): could not execute template: %v")
	}
	return nil
}
