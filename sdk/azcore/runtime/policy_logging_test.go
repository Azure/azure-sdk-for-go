//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestPolicyLoggingSuccess(t *testing.T) {
	rawlog := map[log.Event]string{}
	log.SetListener(func(cls log.Event, s string) {
		rawlog[cls] = s
	})
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := exported.NewPipeline(srv, NewLogPolicy(nil))
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
	pl := exported.NewPipeline(srv, NewLogPolicy(nil))
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

func TestWithAllowedHeadersQueryParams(t *testing.T) {
	rawlog := map[log.Event]string{}
	log.SetListener(func(cls log.Event, s string) {
		rawlog[cls] = s
	})

	const (
		plAllowedHeader = "pipeline-allowed"
		plAllowedQP     = "pipeline-allowed-qp"
		clAllowedHeader = "client-allowed"
		clAllowedQP     = "client-allowed-qp"
		redactedHeader  = "redacted-header"
		redactedQP      = "redacted-qp"
	)

	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithHeader(plAllowedHeader, "received1"), mock.WithHeader(clAllowedHeader, "received2"), mock.WithHeader(redactedHeader, "cantseeme"))

	pl := NewPipeline("", "", PipelineOptions{
		AllowedHeaders:         []string{plAllowedHeader},
		AllowedQueryParameters: []string{plAllowedQP},
	}, &policy.ClientOptions{
		Logging: policy.LogOptions{
			AllowedHeaders:     []string{clAllowedHeader},
			AllowedQueryParams: []string{clAllowedQP},
		},
		Transport: srv,
	})

	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	req.Raw().Header.Set(plAllowedHeader, "sent1")
	req.Raw().Header.Set(clAllowedHeader, "sent2")
	req.Raw().Header.Set(redactedHeader, "cantseeme")
	qp := req.Raw().URL.Query()
	qp.Add(plAllowedQP, "sent1")
	qp.Add(clAllowedQP, "sent2")
	qp.Add(redactedQP, "cantseeme")
	req.Raw().URL.RawQuery = qp.Encode()

	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Len(t, rawlog, 3)
	require.Contains(t, rawlog[log.EventRequest], "?client-allowed-qp=sent2&pipeline-allowed-qp=sent1&redacted-qp=REDACTED")
	require.Regexp(t, `Client-Allowed: sent2\s+Pipeline-Allowed: sent1`, rawlog[log.EventRequest])
	require.Regexp(t, `Client-Allowed: sent2\s+Pipeline-Allowed: sent1`, rawlog[log.EventResponse])
}

func TestSkipWriteReqBody(t *testing.T) {
	req, err := exported.NewRequest(context.Background(), http.MethodGet, "https://contoso.com")
	require.NoError(t, err)

	buf := bytes.Buffer{}
	require.NoError(t, writeReqBody(req, &buf))
	require.Contains(t, buf.String(), "Request contained no body")
	buf.Reset()

	require.NoError(t, req.SetBody(exported.NopCloser(bytes.NewReader([]byte{0xf0, 0x0d})), "application/octet-stream"))
	require.NoError(t, writeReqBody(req, &buf))
	require.Contains(t, buf.String(), "Skip logging body for application/octet-stream")
}

func TestWriteReqBody(t *testing.T) {
	req, err := exported.NewRequest(context.Background(), http.MethodGet, "https://contoso.com")
	require.NoError(t, err)
	require.NoError(t, req.SetBody(exported.NopCloser(strings.NewReader(`{"foo":"bar"}`)), shared.ContentTypeAppJSON))

	buf := bytes.Buffer{}
	require.NoError(t, writeReqBody(req, &buf))
	require.Contains(t, buf.String(), `{"foo":"bar"}`)
}

type readSeekerFailer struct {
	failRead bool
	failSeek bool
}

func (r *readSeekerFailer) Read([]byte) (int, error) {
	if r.failRead {
		return 0, errors.New("read failed")
	}
	return 0, io.EOF
}

func (r *readSeekerFailer) Seek(int64, int) (int64, error) {
	if r.failSeek {
		return 0, errors.New("seek failed")
	}
	// return a positive value to fake that we have content
	return 16, nil
}

func TestWriteReqBodyReadError(t *testing.T) {
	req, err := exported.NewRequest(context.Background(), http.MethodGet, "https://contoso.com")
	require.NoError(t, err)
	rsf := &readSeekerFailer{}
	require.NoError(t, req.SetBody(exported.NopCloser(rsf), shared.ContentTypeAppJSON))

	buf := bytes.Buffer{}
	rsf.failRead = true
	require.Error(t, writeReqBody(req, &buf))
	require.Contains(t, buf.String(), "Failed to read request body: read failed")

	buf.Reset()
	rsf.failRead = false
	rsf.failSeek = true
	require.Error(t, writeReqBody(req, &buf))
	require.Zero(t, buf.Len())
}

func TestSkipWriteRespBody(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	buf := bytes.Buffer{}
	require.NoError(t, writeRespBody(resp, &buf))
	require.Contains(t, buf.String(), "Response contained no body")

	resp.Header.Set(shared.HeaderContentType, "application/octet-stream")
	buf.Reset()
	require.NoError(t, writeRespBody(resp, &buf))
	require.Contains(t, buf.String(), "Skip logging body for application/octet-stream")

	resp.Header.Set(shared.HeaderContentType, "application/json")
	resp.Body = io.NopCloser(strings.NewReader(""))
	buf.Reset()
	require.NoError(t, writeRespBody(resp, &buf))
	require.Contains(t, buf.String(), "Response contained no body")
}

func TestWriteRespBody(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	buf := bytes.Buffer{}

	resp.Header.Set(shared.HeaderContentType, "application/json")
	resp.Body = io.NopCloser(strings.NewReader(`{"foo":"bar"}`))
	require.NoError(t, writeRespBody(resp, &buf))
	require.Contains(t, buf.String(), `{"foo":"bar"}`)
}

func TestWriteRespBodyReadError(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	buf := bytes.Buffer{}

	resp.Header.Set(shared.HeaderContentType, "application/json")
	resp.Body = exported.NopCloser(&readSeekerFailer{failRead: true})
	require.Error(t, writeRespBody(resp, &buf))
	require.Contains(t, buf.String(), "Failed to read response body: read failed")
}
