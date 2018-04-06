// Copyright (c) 2018 Robert Reyna. All rights reserved.
// License BSD 3-Clause https://github.com/robscodebase/eventplanner/blob/master/LICENSE.md
// templates.go contains functions which read html templates
// from the templates folder, return the template to main,
// and populate data using template variables.
package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

// Set localfilePath for running tests during development.
var localfilePath = "/home/robert/gocode/src/robert/eventplanner/src/server"
var dockerfilePath = "/go/src/eventplanner/src/server"
var filePath = dockerfilePath

type eventPlannerTemplate struct {
	eventPlannerTemplate *template.Template
}

// compileTemplate() builds the html template adding the header and footer
// using the name provided by the caller.
func compileTemplate(templateName string) *eventPlannerTemplate {
	sLog(fmt.Sprintf("templates.go: compileTemplate(): template name: %v", templateName))

	// Add the main template file.
	main := template.Must(template.ParseFiles(filePath + "/templates/main.html"))

	// Add the body.
	body := readFile(templateName)
	if templateName == "login.html" {
		// Head and footer are emptry for login.html.
		template.Must(main.New("header").Parse(string("")))
		template.Must(main.New("body").Parse(string(body)))
		template.Must(main.New("footer").Parse(string("")))
		return &eventPlannerTemplate{main.Lookup("main.html")}
	}
	if templateName == "register.html" {
		// Head and footer are empty for login.html.
		template.Must(main.New("header").Parse(string("")))
		template.Must(main.New("body").Parse(string(body)))
		template.Must(main.New("footer").Parse(string("")))
		return &eventPlannerTemplate{main.Lookup("main.html")}
	}

	// Add the header.
	header := readFile("header.html")

	// Add the footer.
	footer := readFile("footer.html")

	// Combine header, body, and footer and return template.
	template.Must(main.New("header").Parse(string(header)))
	template.Must(main.New("body").Parse(string(body)))
	template.Must(main.New("footer").Parse(string(footer)))
	sLog(fmt.Sprintf("templates.go: compileTemplate(): template parse successful return: %v", templateName))
	return &eventPlannerTemplate{main.Lookup("main.html")}
}

// readFile() takes a file name and returns
// a []byte.  readFile() panics on error.
func readFile(fileName string) []byte {
	sLog(fmt.Sprintf("templates.go: readFile(): file name: %v", fileName))
	template, err := ioutil.ReadFile(filePath + "/templates/" + fileName)
	if err != nil {
		panic(fmt.Errorf("templates.go: readFile(): could not read file: %v: %v", fileName, err))
	}
	return template
}

// runTemplate() combines template variables and the html template for final delivery to client.
func (template *eventPlannerTemplate) runTemplate(w http.ResponseWriter, r *http.Request, input interface{}) *errorMessage {
	sLog(fmt.Sprintf("templates.go: runTemplate(): input: %v", input))
	if err := template.eventPlannerTemplate.Execute(w, struct{ Input interface{} }{Input: input}); err != nil {
		return &errorMessage{Error: err, Message: fmt.Sprintf("templates.go: runTemplate(): could not execute template: %v", err)}
	}
	return nil
}
