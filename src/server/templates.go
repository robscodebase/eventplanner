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

var filePath = "/go/src/eventplanner/src/server"

type eventPlannerTemplate struct {
	eventPlannerTemplate *template.Template
}

// compileTemplate() builds the html template adding the header and footer
// using the name provided by the caller.
func compileTemplate(templateName string) *eventPlannerTemplate {
	slog(fmt.Sprintf("templates.go: compileTemplate(): template name: %v", templateName))
	// For template login.html and register.html a header and footer
	// is not needed.
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
	slog(fmt.Sprintf("templates.go: compileTemplate(): template pase successful return: %v", templateName))
	return &eventPlannerTemplate{main.Lookup("main.html")}
}

// readFile() takes a file name and returns
// a []byte.  readFile() panics on error.
func readFile(fileName string) []byte {
	slog(fmt.Sprintf("templates.go: readFile(): file name: %v", fileName))
	template, err := ioutil.ReadFile(filePath + "/templates/" + fileName)
	if err != nil {
		panic(fmt.Errorf("templates.go: readFile(): could not read file: %v: %v", fileName, err))
	}
	return template
}

// runTemplate() combines template variables and the html template for final delivery to client.
func (template *eventPlannerTemplate) runTemplate(w http.ResponseWriter, r *http.Request, input interface{}) *errorMessage {
	slog(fmt.Sprintf("templates.go: runTemplate(): file name: %v", input))
	if err := template.eventPlannerTemplate.Execute(w, input); err != nil {
		return fmt.Errorf(err, "templates.go: runTemplate(): could not execute template: %v")
	}
	return nil
}
