// Copyright (c) 2018 Robert Reyna. All rights reserved.
// License BSD 3-Clause https://github.com/robscodebase/eventplanner/blob/master/LICENSE.md
package main

import (
	"log"
	"testing"
)

func TestAdd2(t *testing.T) {
	replyFromAdd2 := add2(2, 2)
	if replyFromAdd2 != 4 {
		log.Fatalf("main_test.go: TestAdd2(): replyFromAdd2 call to add(2,2) fail: want 4; got %v", replyFromAdd2)
	}
}
