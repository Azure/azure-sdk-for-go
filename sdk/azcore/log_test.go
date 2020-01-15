// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import "testing"

func TestLoggingDefault(t *testing.T) {
	// ensure logging with nil listener doesn't fail
	Log().SetListener(nil)
	Log().Write(LogError, "this should work just fine")

	log := map[LogClassification]string{}
	Log().SetListener(func(cls LogClassification, msg string) {
		log[cls] = msg
	})
	const req = "this is a request"
	Log().Write(LogRequest, req)
	const resp = "this is a response"
	Log().Write(LogResponse, resp)
	if l := len(log); l != 2 {
		t.Fatalf("unexpected log entry count: %d", l)
	}
	if log[LogRequest] != req {
		t.Fatalf("unexpected log request: %s", log[LogRequest])
	}
	if log[LogResponse] != resp {
		t.Fatalf("unexpected log response: %s", log[LogResponse])
	}
}

func TestLoggingClassification(t *testing.T) {
	log := map[LogClassification]string{}
	Log().SetListener(func(cls LogClassification, msg string) {
		log[cls] = msg
	})
	Log().SetClassifications(LogError)
	defer Log().resetClassifications()
	Log().Write(LogSlowResponse, "this shouldn't be in the log")
	if s, ok := log[LogSlowResponse]; ok {
		t.Fatalf("unexpected log entry %s", s)
	}
	const err = "this is an error"
	Log().Write(LogError, err)
	if log[LogError] != err {
		t.Fatalf("unexpected log entry: %s", log[LogError])
	}
}
