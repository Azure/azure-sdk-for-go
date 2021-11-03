//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestPolicyLoggingSuccess(t *testing.T) {
	rawlog := map[log.Event]string{}
	log.SetListener(func(cls log.Event, s string) {
		rawlog[cls] = s
	})
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := pipeline.NewPipeline(srv, NewLogPolicy(nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	qp := req.Raw().URL.Query()
	qp.Set("api-version", "12345")
	qp.Set("sig", "redact_me")
	req.Raw().URL.RawQuery = qp.Encode()
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if logReq, ok := rawlog[log.EventRequest]; ok {
		// Request ==> OUTGOING REQUEST (Try=1)
		// 	GET http://127.0.0.1:49475?one=fish&sig=REDACTED
		// 	(no headers)
		if !strings.Contains(logReq, "(no headers)") {
			t.Fatal("missing (no headers)")
		}
		if !strings.Contains(logReq, "api-version=12345") {
			t.Fatal("didn't find api-version query param")
		}
		if strings.Contains(logReq, "sig=redact_me") {
			t.Fatal("sig query param wasn't redacted")
		}
	} else {
		t.Fatal("missing LogRequest")
	}
	if logResp, ok := rawlog[log.EventResponse]; ok {
		// Response ==> REQUEST/RESPONSE (Try=1/1.0034ms, OpTime=1.0034ms) -- RESPONSE SUCCESSFULLY RECEIVED
		// 	GET http://127.0.0.1:49475?one=fish&sig=REDACTED
		// 	(no headers)
		// 	--------------------------------------------------------------------------------
		// 	RESPONSE Status: 200 OK
		// 	Content-Length: [0]
		// 	Date: [Fri, 22 Nov 2019 23:48:02 GMT]
		if !strings.Contains(logResp, "RESPONSE Status: 200 OK") {
			t.Fatal("missing response status")
		}
	} else {
		t.Fatal("missing LogResponse")
	}
}

func TestPolicyLoggingError(t *testing.T) {
	rawlog := map[log.Event]string{}
	log.SetListener(func(cls log.Event, s string) {
		rawlog[cls] = s
	})
	srv, close := mock.NewServer()
	defer close()
	srv.SetError(errors.New("bogus error"))
	pl := pipeline.NewPipeline(srv, NewLogPolicy(nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Raw().Header.Add("header", "one")
	req.Raw().Header.Add("Authorization", "redact")
	resp, err := pl.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("unexpected respose")
	}
	if logReq, ok := rawlog[log.EventRequest]; ok {
		// Request ==> OUTGOING REQUEST (Try=1)
		// 	GET http://127.0.0.1:50057
		// 	Authorization: REDACTED
		// 	Header: [one]
		if !strings.Contains(logReq, "Authorization: REDACTED") {
			t.Fatal("missing redacted authorization header")
		}
	} else {
		t.Fatal("missing LogRequest")
	}
	if logResponse, ok := rawlog[log.EventResponse]; ok {
		// Response ==> REQUEST/RESPONSE (Try=1/0s, OpTime=0s) -- REQUEST ERROR
		// 	GET http://127.0.0.1:50057
		// 	Authorization: REDACTED
		// 	Header: [one]
		// 	--------------------------------------------------------------------------------
		// 	ERROR:
		// 	bogus error
		// 	 ...stack track...
		if !strings.Contains(logResponse, "Authorization: REDACTED") {
			t.Fatal("missing redacted authorization header")
		}
		if !strings.Contains(logResponse, "bogus error") {
			t.Fatal("missing error message")
		}
	} else {
		t.Fatal("missing LogResponse")
	}
}

func TestShouldLogBody(t *testing.T) {
	b := bytes.NewBuffer(make([]byte, 64))
	if shouldLogBody(b, "application/octet-stream") {
		t.Fatal("shouldn't log for application/octet-stream")
	} else if b.Len() == 0 {
		t.Fatal("skip logging should write skip message to buffer")
	}
	b.Reset()
	if !shouldLogBody(b, "application/json") {
		t.Fatal("should log for application/json")
	} else if b.Len() != 0 {
		t.Fatal("logging shouldn't write message")
	}
	if !shouldLogBody(b, "application/xml") {
		t.Fatal("should log for application/xml")
	} else if b.Len() != 0 {
		t.Fatal("logging shouldn't write message")
	}
	if !shouldLogBody(b, "text/plain") {
		t.Fatal("should log for text/plain")
	} else if b.Len() != 0 {
		t.Fatal("logging shouldn't write message")
	}
}
