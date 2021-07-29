// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/logger"
)

func TestLoggingDefault(t *testing.T) {
	// ensure logging with nil listener doesn't fail
	SetListener(nil)
	logger.Log().Write(logger.LogRequest, "this should work just fine")

	log := map[LogClassification]string{}
	SetListener(func(cls LogClassification, msg string) {
		log[cls] = msg
	})
	const req = "this is a request"
	logger.Log().Write(logger.LogRequest, req)
	const resp = "this is a response: %d"
	logger.Log().Writef(logger.LogResponse, resp, http.StatusOK)
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
	SetListener(func(cls LogClassification, msg string) {
		log[cls] = msg
	})
	SetClassifications(LogRequest)
	defer resetClassifications()
	logger.Log().Write(logger.LogResponse, "this shouldn't be in the log")
	if s, ok := log[LogResponse]; ok {
		t.Fatalf("unexpected log entry %s", s)
	}
	const req = "this is a request"
	logger.Log().Write(logger.LogRequest, req)
	if log[LogRequest] != req {
		t.Fatalf("unexpected log entry: %s", log[LogRequest])
	}
}
