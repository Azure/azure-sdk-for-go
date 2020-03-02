// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestPolicyLoggingSuccess(t *testing.T) {
	log := map[LogClassification]string{}
	Log().SetListener(func(cls LogClassification, s string) {
		log[cls] = s
	})
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := NewPipeline(srv, NewRequestLogPolicy(RequestLogOptions{}))
	req := NewRequest(http.MethodGet, srv.URL())
	qp := req.URL.Query()
	qp.Set("one", "fish")
	qp.Set("sig", "redact")
	req.URL.RawQuery = qp.Encode()
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if logReq, ok := log[LogRequest]; ok {
		// Request ==> OUTGOING REQUEST (Try=1)
		// 	GET http://127.0.0.1:49475?one=fish&sig=REDACTED
		// 	(no headers)
		if !strings.Contains(logReq, "sig=REDACTED") {
			t.Fatal("missing redacted sig query param")
		}
		if !strings.Contains(logReq, "(no headers)") {
			t.Fatal("missing (no headers)")
		}
	} else {
		t.Fatal("missing LogRequest")
	}
	if logResp, ok := log[LogResponse]; ok {
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
	log := map[LogClassification]string{}
	Log().SetListener(func(cls LogClassification, s string) {
		log[cls] = s
	})
	srv, close := mock.NewServer()
	defer close()
	srv.SetError(errors.New("bogus error"))
	pl := NewPipeline(srv, NewRequestLogPolicy(RequestLogOptions{}))
	req := NewRequest(http.MethodGet, srv.URL())
	req.Header.Add("header", "one")
	req.Header.Add("Authorization", "redact")
	resp, err := pl.Do(context.Background(), req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("unexpected respose")
	}
	if logReq, ok := log[LogRequest]; ok {
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
	if logError, ok := log[LogError]; ok {
		// Error ==> REQUEST/RESPONSE (Try=1/0s, OpTime=0s) -- REQUEST ERROR
		// 	GET http://127.0.0.1:50057
		// 	Authorization: REDACTED
		// 	Header: [one]
		// 	--------------------------------------------------------------------------------
		// 	ERROR:
		// 	bogus error
		// 	 ...stack track...
		if !strings.Contains(logError, "Authorization: REDACTED") {
			t.Fatal("missing redacted authorization header")
		}
		if !strings.Contains(logError, "bogus error") {
			t.Fatal("missing error message")
		}
	} else {
		t.Fatal("missing LogError")
	}
}

// TODO: add test for slow response
