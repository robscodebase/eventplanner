// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// W holds the configuration for a web test.
type WebTest struct {
	t      *testing.T
	host   string
	Client *http.Client
}

var webTest *WebTest

func NewWebTest(t *testing.T, host string) *WebTest {
	return &WebTest{
		t:      t,
		host:   host,
		Client: http.DefaultClient,
	}
}

func TestMain(m *testing.M) {
	serv := httptest.NewServer(nil)
	webTest = NewWebTest(nil, serv.Listener.Addr().String())
	runHandlers()

	os.Exit(m.Run())
}

func TestMainFunc(t *testing.T) {
	webTest := NewWebTest(t, "localhost:8080")
	m := BuildMain(t)
	defer m.Cleanup()
	m.Run(nil, func() {
		webTest.WaitForNet()
		bodyContains(t, webTest, "/", "Event Planner")
	})
}

func TestNoBooks(t *testing.T) {
	bodyContains(t, webTest, "/", "Event Planner")
}

func BuildMain(t *testing.T) *Runner {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp, err := ioutil.TempDir("", "runmain-"+filepath.Base(workingDirectory)+"-")
	if err != nil {
		t.Fatal(err)
	}

	r := &Runner{t: t, tmp: tmp}

	bin := filepath.Join(tmp, "a.out")
	cmd := exec.Command("go", "build", "-o", bin)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("go build: %v\n%s", err, out)
		return r
	}

	r.bin = bin
	return r
}

type Runner struct {
	t   *testing.T
	tmp string
	bin string
}

func (r *Runner) Built() bool {
	return r.bin != ""
}

func (r *Runner) Cleanup() {
	if err := os.RemoveAll(r.tmp); err != nil {
		r.t.Error(err)
	}
}

func (r *Runner) Run(env map[string]string, f func()) {
	if !r.Built() {
		r.t.Error("Tried to run when binary not built.")
		return
	}
	environ := os.Environ()
	for k, v := range env {
		environ = append(environ, k+"="+v)
	}

	cmd := exec.Command(r.bin)
	cmd.Env = environ

	if err := cmd.Start(); err != nil {
		r.t.Error(err)
		return
	}

	f()

	done := make(chan struct{})
	go func() {
		cmd.Wait()
		close(done)
	}()

	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		r.t.Error(err)
	}

	select {
	case <-time.After(5 * time.Second):
		r.t.Error("Timed out with SIGINT, trying SIGKILL.")
		if err := cmd.Process.Kill(); err != nil {
			r.t.Error(err)
		}
	case <-done:
	}
}

func (webTest *WebTest) WaitForNet() {
	const retryDelay = 100 * time.Millisecond
	deadline := time.Now().Add(30 * time.Second)

	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", webTest.host)
		if err != nil {
			time.Sleep(retryDelay)
			continue
		}
		conn.Close()
		return
	}

	webTest.t.Fatalf("Timed out wating for net %s", webTest.host)
}

func bodyContains(t *testing.T, webTest *WebTest, path, contains string) (ok bool) {
	body, _, err := webTest.GetBody(path)
	if err != nil {
		t.Error(err)
		return false
	}
	if !strings.Contains(body, contains) {
		t.Errorf("want %s to contain %s", body, contains)
		return false
	}
	return true
}

func (webTest *WebTest) GetBody(path string) (body string, resp *http.Response, err error) {
	resp, err = webTest.Get(path)
	if err != nil {
		return "", resp, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp, err
	}
	return string(b), resp, err
}

func (webTest *WebTest) Get(path string) (*http.Response, error) {
	return webTest.Client.Get("http://" + webTest.host + path)
}
