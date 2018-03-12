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
	"strings"
	"testing"
)

var testDB *sql.DB

func TestDBMaker(t *testing.T) {
	var err error
	testDB, err = sql.Open("mysql", "root:insecure@(mysql-event-planner:3306)/mysql")
	if err != nil {
		t.Fatalf("main_test.go: TestDBMaker(): sql.Open(): error: %v", err)
	}
	err = isDB(testDB)
	if err != nil {
		t.Fatalf("main_test.go: TestDBMaker(): isDB(): unable to verify and create db: error: %v", err)
	}
}

// funcTestRunHandlers checks that the handler is
// operational. It does not check for accuracy.
func TestRunHandlers(t *testing.T) {
	var server = httptest.NewServer(runHandlers())
	defer server.Close()
	testHandlerTable := []struct {
		pageName  string
		pageTitle string
	}{
	//{pageName: "login", pageTitle: "<title>Event Planner - Login</title>"},
	//{pageName: "add-event", pageTitle: "<title>Event Planner - Add</title>"},
	//{pageName: "edit-event/1", pageTitle: "<title>Event Planner - Edit</title>"},
	//{pageName: "view-events", pageTitle: "<title>Event Planner - View</title>"},
	//{pageName: "register", pageTitle: "<title>Event Planner - Register</title>"},
	//{pageName: "logout", pageTitle: "<title>Event Planner - Logout</title>"},
	}
	for _, v := range testHandlerTable {
		response := testGetHTTP(server, v.pageName)
		responseString := string(bytes.TrimSpace(testReadBody(response)))
		if responseString == "404 page not found" {
			log.Printf("main_test.go: TestRunHandlers(): testGetHTTP(): pageName: %v: want; body: got; %v", v.pageName, responseString)
		} else if strings.Contains(responseString, v.pageTitle) != true {
			log.Printf("main_test.go: TestRunHandlers(): responseString error: want %v; got %v", v.pageTitle, responseString)
		}
	}
}

// testReadBody() takes a response a returns
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
