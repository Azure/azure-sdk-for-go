// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package logger

import (
	"fmt"
	"net/http"
	"testing"
)

func TestLoggingdefault(t *testing.T) {
	// ensure logging with nil listener doesn't fail
	Log().SetListener(nil)
	Log().Write(LogRequest, "this should work just fine")

	log := map[LogClassification]string{}
	Log().SetListener(func(cls LogClassification, msg string) {
		log[cls] = msg
	})
	const req = "this is a request"
	Log().Write(LogRequest, req)
	const resp = "this is a response: %d"
	Log().Writef(LogResponse, resp, http.StatusOK)
	if l := len(log); l != 2 {
		t.Fatalf("unexpected log entry count: %d", l)
	}
	if log[LogRequest] != req {
		t.Fatalf("unexpected log request: %s", log[LogRequest])
	}
	if log[LogResponse] != fmt.Sprintf(resp, http.StatusOK) {
		t.Fatalf("unexpected log response: %s", log[LogResponse])
	}
}

func TestLoggingClassification(t *testing.T) {
	log := map[LogClassification]string{}
	Log().SetListener(func(cls LogClassification, msg string) {
		log[cls] = msg
	})
	Log().SetClassifications(LogRequest)
	defer Log().resetClassifications()
	Log().Write(LogResponse, "this shouldn't be in the log")
	if s, ok := log[LogResponse]; ok {
		t.Fatalf("unexpected log entry %s", s)
	}
	const req = "this is a request"
	Log().Write(LogRequest, req)
	if log[LogRequest] != req {
		t.Fatalf("unexpected log entry: %s", log[LogRequest])
	}
}
