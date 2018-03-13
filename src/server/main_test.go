// Copyright (c) 2018 Robert Reyna. All rights reserved.
// License BSD 3-Clause https://github.com/robscodebase/eventplanner/blob/master/LICENSE.md
package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var server *httptest.Server

// TestDBMaker() opens a sql connection
// and runs dbMaker() which runs isDB(),
// confirming the database tables are in place.
func TestDBMaker(t *testing.T) {
	// testingSession = true bypasses the verifySession() function
	// to allow for testing of the handlers using the demo user data.
	testingSession = true
	var err error
	db, err = sql.Open("mysql", "root:insecure@(172.17.0.2:3306)/mysql")
	if err != nil {
		t.Fatalf("main_test.go: dbMaker(): sql.Open(): error: %v", err)
	}
	db, err = dbMaker(db, "isDB", "main_test.go: dbMaker(): call to dbMaker(): isDB():")
	if err != nil {
		t.Fatalf("main_test.go: dbMaker(): dbMaker(): could not make db: err: %v", err)
	}
}

// TestCreateDemoDB() creates the demo database and populates it with events for testing handlers.
func TestCreateDemoDB(t *testing.T) {
	createdEvents, user, err := createDemoDB(db)
	if err != nil {
		t.Fatalf("main_test.go: TestCreateDemoDB(): err: createDemoDB(): error: %v", err)
	}
	if user == "" {
		t.Fatalf("main_test.go: TestCreateDemoDB(): user: createDemoDB(): user should be demo: user: %v", user)
	}
	if createdEvents == 0 {
		t.Fatalf("main_test.go: TestCreateDemoDB(): createdEvents: createDemoDB(): createdEvents: %v", createdEvents)
	}
}

// funcTestRunHandlers checks that the handler is
// operational. It checks if the title matches.
func TestRunHandlers(t *testing.T) {
	var server = httptest.NewServer(runHandlers())
	defer server.Close()
	testHandlerTable := []struct {
		pageName  string
		pageTitle string
	}{
		{pageName: "login", pageTitle: "<title>Event Planner - View Events</title>"},
		{pageName: "add-event", pageTitle: "<title>Event Planner - Add Event</title>"},
		{pageName: "edit-event/1", pageTitle: "<title>Event Planner - Edit Event</title>"},
		{pageName: "view-events", pageTitle: "<title>Event Planner - View Events</title>"},
		{pageName: "register", pageTitle: "<title>Event Planner - Register</title>"},
		{pageName: "logout", pageTitle: "<title>Event Planner - View Events</title>"},
	}
	for _, v := range testHandlerTable {
		response := testGetHTTP(server, v.pageName)
		responseString := string(bytes.TrimSpace(testReadBody(response)))
		if responseString == "404 page not found" {
			log.Printf("main_test.go: TestRunHandlers(): testGetHTTP(): pageName: %v: want; body: got; %v", v.pageName, responseString)
		} else if strings.Contains(responseString, v.pageTitle) != true {
			log.Printf("main_test.go: TestRunHandlers(): pageName: %v, responseString error: want %v; got %v", v.pageName, v.pageTitle, responseString)
		}
	}
}

// testReadBody() takes a response and returns
// the body in string format.
func testReadBody(response *http.Response) []byte {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("main_test.go: testReadBody(): ioutil.ReadAll: error: %v", err)
	}
	return body
}

// testGetHTTP() takes a server and a url name
// and returns an *http.Response.
func testGetHTTP(server *httptest.Server, request string) *http.Response {
	response, err := http.Get(fmt.Sprintf("%s/%s", server.URL, request))
	if err != nil {
		log.Fatalf("main_test.go: TestGetHTTP(): http.Get() error: %v", err)
	}
	return response
}

// testPostHTTP() takes a server and url
// and returns an *http.Response.
func testPostHTTP(server *httptest.Server, request string, data url.Values) *http.Response {
	response, err := http.PostForm(fmt.Sprintf("%s/%s", server.URL, request), data)
	if err != nil {
		log.Fatalf("main_test.go: TestPostHTTP(): http.Post() error: %v", err)
	}
	return response
}

// TestCloseServer() closes the server at the end of testing.
func closeServer(server *http.Server) {
	server.Close()
}
