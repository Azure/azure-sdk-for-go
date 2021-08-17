//go:build go1.16
// +build go1.16

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
	log.Write(log.Request, "this should work just fine")

	testlog := map[Classification]string{}
	SetListener(func(cls Classification, msg string) {
		testlog[cls] = msg
	})
	const req = "this is a request"
	log.Write(log.Request, req)
	const resp = "this is a response: %d"
	log.Writef(log.Response, resp, http.StatusOK)
	if l := len(testlog); l != 2 {
		t.Fatalf("unexpected log entry count: %d", l)
	}
	if testlog[Request] != req {
		t.Fatalf("unexpected log request: %s", testlog[Request])
	}
	if testlog[Response] != fmt.Sprintf(resp, http.StatusOK) {
		t.Fatalf("unexpected log response: %s", testlog[Response])
	}
}

func TestLoggingClassification(t *testing.T) {
	testlog := map[Classification]string{}
	SetListener(func(cls Classification, msg string) {
		testlog[cls] = msg
	})
	SetClassifications(Request)
	defer resetClassifications()
	log.Write(log.Response, "this shouldn't be in the log")
	if s, ok := testlog[Response]; ok {
		t.Fatalf("unexpected log entry %s", s)
	}
	const req = "this is a request"
	log.Write(log.Request, req)
	if testlog[Request] != req {
		t.Fatalf("unexpected log entry: %s", testlog[Request])
	}
}
