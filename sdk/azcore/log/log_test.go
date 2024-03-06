//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package log

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

func TestLoggingDefault(t *testing.T) {
	// ensure logging with nil listener doesn't fail
	SetListener(nil)
	log.Write(EventRequest, "this should work just fine")

	testlog := map[Event]string{}
	SetListener(func(cls Event, msg string) {
		testlog[cls] = msg
	})
	const req = "this is a request"
	log.Write(EventRequest, req)
	const resp = "this is a response: %d"
	log.Writef(EventResponse, resp, http.StatusOK)
	if l := len(testlog); l != 2 {
		t.Fatalf("unexpected log entry count: %d", l)
	}
	if testlog[EventRequest] != req {
		t.Fatalf("unexpected log request: %s", testlog[EventRequest])
	}
	if testlog[EventResponse] != fmt.Sprintf(resp, http.StatusOK) {
		t.Fatalf("unexpected log response: %s", testlog[EventResponse])
	}
}

func TestLoggingEvent(t *testing.T) {
	testlog := map[Event]string{}
	SetListener(func(cls Event, msg string) {
		testlog[cls] = msg
	})
	SetEvents(EventRequest)
	defer resetEvents()
	log.Write(EventResponse, "this shouldn't be in the log")
	if s, ok := testlog[EventResponse]; ok {
		t.Fatalf("unexpected log entry %s", s)
	}
	const req = "this is a request"
	log.Write(EventRequest, req)
	if testlog[EventRequest] != req {
		t.Fatalf("unexpected log entry: %s", testlog[EventRequest])
	}
}
