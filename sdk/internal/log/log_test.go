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
	Write(RequestEvent, "this should work just fine")

	log := map[Event]string{}
	SetListener(func(cls Event, msg string) {
		log[cls] = msg
	})
	const req = "this is a request"
	Write(RequestEvent, req)
	const resp = "this is a response: %d"
	Writef(ResponseEvent, resp, http.StatusOK)
	if l := len(log); l != 2 {
		t.Fatalf("unexpected log entry count: %d", l)
	}
	if log[RequestEvent] != req {
		t.Fatalf("unexpected log request: %s", log[RequestEvent])
	}
	if log[ResponseEvent] != fmt.Sprintf(resp, http.StatusOK) {
		t.Fatalf("unexpected log response: %s", log[ResponseEvent])
	}
}

func TestLoggingEvent(t *testing.T) {
	log := map[Event]string{}
	SetListener(func(cls Event, msg string) {
		log[cls] = msg
	})
	SetEvents(RequestEvent)
	defer TestResetEvents()
	Write(ResponseEvent, "this shouldn't be in the log")
	if s, ok := log[ResponseEvent]; ok {
		t.Fatalf("unexpected log entry %s", s)
	}
	const req = "this is a request"
	Write(RequestEvent, req)
	if log[RequestEvent] != req {
		t.Fatalf("unexpected log entry: %s", log[RequestEvent])
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
