// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package logger

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestLoggingDefault(t *testing.T) {
	// ensure logging with nil listener doesn't fail
	azcore.SetListener(nil)
	Log().Write(azcore.LogRequest, "this should work just fine")

	log := map[azcore.LogClassification]string{}
	azcore.SetListener(func(cls azcore.LogClassification, msg string) {
		log[cls] = msg
	})
	const req = "this is a request"
	Log().Write(azcore.LogRequest, req)
	const resp = "this is a response: %d"
	Log().Writef(azcore.LogResponse, resp, http.StatusOK)
	if l := len(log); l != 2 {
		t.Fatalf("unexpected log entry count: %d", l)
	}
	if log[azcore.LogRequest] != req {
		t.Fatalf("unexpected log request: %s", log[azcore.LogRequest])
	}
	if log[azcore.LogResponse] != fmt.Sprintf(resp, http.StatusOK) {
		t.Fatalf("unexpected log response: %s", log[azcore.LogResponse])
	}
}

func TestLoggingClassification(t *testing.T) {
	log := map[azcore.LogClassification]string{}
	azcore.SetListener(func(cls azcore.LogClassification, msg string) {
		log[cls] = msg
	})
	azcore.SetClassifications(azcore.LogRequest)
	defer Log().resetClassifications()
	Log().Write(azcore.LogResponse, "this shouldn't be in the log")
	if s, ok := log[azcore.LogResponse]; ok {
		t.Fatalf("unexpected log entry %s", s)
	}
	const req = "this is a request"
	Log().Write(azcore.LogRequest, req)
	if log[azcore.LogRequest] != req {
		t.Fatalf("unexpected log entry: %s", log[azcore.LogRequest])
	}
}
