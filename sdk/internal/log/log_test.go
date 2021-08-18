// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package log

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestLoggingdefault(t *testing.T) {
	// ensure logging with nil listener doesn't fail
	SetListener(nil)
	Write(Request, "this should work just fine")

	log := map[Classification]string{}
	SetListener(func(cls Classification, msg string) {
		log[cls] = msg
	})
	const req = "this is a request"
	Write(Request, req)
	const resp = "this is a response: %d"
	Writef(Response, resp, http.StatusOK)
	if l := len(log); l != 2 {
		t.Fatalf("unexpected log entry count: %d", l)
	}
	if log[Request] != req {
		t.Fatalf("unexpected log request: %s", log[Request])
	}
	if log[Response] != fmt.Sprintf(resp, http.StatusOK) {
		t.Fatalf("unexpected log response: %s", log[Response])
	}
}

func TestLoggingClassification(t *testing.T) {
	log := map[Classification]string{}
	SetListener(func(cls Classification, msg string) {
		log[cls] = msg
	})
	SetClassifications(Request)
	defer TestResetClassifications()
	Write(Response, "this shouldn't be in the log")
	if s, ok := log[Response]; ok {
		t.Fatalf("unexpected log entry %s", s)
	}
	const req = "this is a request"
	Write(Request, req)
	if log[Request] != req {
		t.Fatalf("unexpected log entry: %s", log[Request])
	}
}

func TestEnvironment(t *testing.T) {
	os.Setenv("AZURE_SDK_GO_LOGGING", "all")
	defer os.Unsetenv("AZURE_SDK_GO_LOGGING")
	initLogging()
	if log.lst == nil {
		t.Fatal("unexpected nil listener")
	}
}
