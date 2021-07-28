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

	log := map[logger.LogClassification]string{}
	SetListener(func(cls logger.LogClassification, msg string) {
		log[cls] = msg
	})
	const req = "this is a request"
	logger.Log().Write(logger.LogRequest, req)
	const resp = "this is a response: %d"
	logger.Log().Writef(logger.LogResponse, resp, http.StatusOK)
	if l := len(log); l != 2 {
		t.Fatalf("unexpected log entry count: %d", l)
	}
	if log[logger.LogRequest] != req {
		t.Fatalf("unexpected log request: %s", log[logger.LogRequest])
	}
	if log[logger.LogResponse] != fmt.Sprintf(resp, http.StatusOK) {
		t.Fatalf("unexpected log response: %s", log[logger.LogResponse])
	}
}

func TestLoggingClassification(t *testing.T) {
	log := map[logger.LogClassification]string{}
	SetListener(func(cls logger.LogClassification, msg string) {
		log[cls] = msg
	})
	SetClassifications(LogRequest)
	logger.Log().Write(logger.LogResponse, "this shouldn't be in the log")
	if s, ok := log[logger.LogResponse]; ok {
		t.Fatalf("unexpected log entry %s", s)
	}
	const req = "this is a request"
	logger.Log().Write(logger.LogRequest, req)
	if log[logger.LogRequest] != req {
		t.Fatalf("unexpected log entry: %s", log[logger.LogRequest])
	}
}
