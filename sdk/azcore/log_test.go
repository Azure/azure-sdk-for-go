// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

func TestLoggingDefault(t *testing.T) {
	// ensure logging with nil listener doesn't fail
	LogSetListener(nil)
	log.Write(log.Request, "this should work just fine")

	testlog := map[LogClassification]string{}
	LogSetListener(func(cls LogClassification, msg string) {
		testlog[cls] = msg
	})
	const req = "this is a request"
	log.Write(log.Request, req)
	const resp = "this is a response: %d"
	log.Writef(log.Response, resp, http.StatusOK)
	if l := len(testlog); l != 2 {
		t.Fatalf("unexpected log entry count: %d", l)
	}
	if testlog[LogRequest] != req {
		t.Fatalf("unexpected log request: %s", testlog[LogRequest])
	}
	if testlog[LogResponse] != fmt.Sprintf(resp, http.StatusOK) {
		t.Fatalf("unexpected log response: %s", testlog[LogResponse])
	}
}

func TestLoggingClassification(t *testing.T) {
	testlog := map[LogClassification]string{}
	LogSetListener(func(cls LogClassification, msg string) {
		testlog[cls] = msg
	})
	LogSetClassifications(LogRequest)
	defer resetClassifications()
	log.Write(log.Response, "this shouldn't be in the log")
	if s, ok := testlog[LogResponse]; ok {
		t.Fatalf("unexpected log entry %s", s)
	}
	const req = "this is a request"
	log.Write(log.Request, req)
	if testlog[LogRequest] != req {
		t.Fatalf("unexpected log entry: %s", testlog[LogRequest])
	}
}
