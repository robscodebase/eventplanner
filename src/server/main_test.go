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

func TestRunHandlers(t *testing.T) {
	srv := httptest.NewServer(runHandlers())
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/login", srv.URL))
	if err != nil {
		log.Fatalf("main_test.go: TestRunHandlers(): http.Get() error: %v", err)
	}

}
