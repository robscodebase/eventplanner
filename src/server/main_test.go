// Copyright (c) 2018 Robert Reyna. All rights reserved.
// License BSD 3-Clause https://github.com/robscodebase/eventplanner/blob/master/LICENSE.md
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// funcTestRunHandlers checks that the handler is
// operational. It does not check for accuracy.
func TestRunHandlers(t *testing.T) {
	server := httptest.NewServer(runHandlers())
	defer server.Close()
	_ = testGetHTTP(server, "login")

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
