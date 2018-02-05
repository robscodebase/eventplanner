package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

type eventPlannerTemplate struct {
    eventPlannerTemplate *template.Template
}

func compileTemplate(templateName string) *eventPlannerTemplate {
    // Add the header.
    header := template.Must(template.ParseFiles("templates/header.html"))

    // Add the body.
    body, err := ioutil.ReadFile("templates" + templateName)
    if err != nil {
        panic(fmt.Errorf("compileTemplate() could not read body: %v", err))
    }

    // Add the footer.
    footer, err := ioutil.ReadFile("templates/footer.html")
    if err != nil {
        panic(fmt.Errorf("compileTemplate() could not read footer: %v", err))
    }

    // Combine header, body, and footer and return template.
    template.Must(header.New("body").Parse(string(body)))
    template.Must(header.New("footer").Parse(string(footer)))
    return &appTemplate{tmpl.Lookup("header.html")}
}

func (template *eventPlannerTemplate) runTemplate(w http.ResponseWriter, r *http.Request, input interface{}) *errorMessage {
    userInfo := struct {
        Input        interface{}
        AuthEnabled bool
        UserData     *UserData
    }{
        Input:        input,
        AuthEnabled: true,
    }

    if Login {
        user = userSessionData(r)
    }

    if err := template.eventPlannerTemplate.Execute(w, userInfo); err != nil {
        return formatError(err, "could not write template: %v")
    }
    return nil
}

func userSessionData(r *http.Request) *UserData {
    return &UserData{
        ID:          "SECURE ID",
        UserName: "ADMIN",
    }
}

type UserData struct {
    ID, DisplayName string
}
