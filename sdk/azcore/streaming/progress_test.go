//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package streaming_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestProgressReporting(t *testing.T) {
	const contentSize = 4096
	content := make([]byte, contentSize)
	for i := 0; i < contentSize; i++ {
		content[i] = byte(i % 255)
	}
	body := bytes.NewReader(content)
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody(content))
	pl := exported.NewPipeline(srv)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	runtime.SkipBodyDownload(req)
	var bytesSent int64
	reqRpt := streaming.NewRequestProgress(streaming.NopCloser(body), func(bytesTransferred int64) {
		bytesSent = bytesTransferred
	})
	if err := req.SetBody(reqRpt, "application/octet-stream"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var bytesReceived int64
	respRpt := streaming.NewResponseProgress(resp.Body, func(bytesTransferred int64) {
		bytesReceived = bytesTransferred
	})
	defer respRpt.Close()
	b, err := io.ReadAll(respRpt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bytesSent != contentSize {
		t.Fatalf("wrong bytes sent: %d", bytesSent)
	}
	if bytesReceived != contentSize {
		t.Fatalf("wrong bytes received: %d", bytesReceived)
	}
	if !reflect.DeepEqual(content, b) {
		t.Fatal("request and response bodies don't match")
	}
}

// Ensure there is a seek to 0
// do some reading, call a seek, then make sure reads are from the beginning
func TestProgressReportingSeek(t *testing.T) {
	const contentSize = 4096
	content := make([]byte, contentSize)
	for i := 0; i < contentSize; i++ {
		content[i] = byte(i % 255)
	}
	body := bytes.NewReader(content)
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody(content))
	pl := exported.NewPipeline(srv)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	runtime.SkipBodyDownload(req)
	var bytesSent int64
	reqRpt := streaming.NewRequestProgress(streaming.NopCloser(body), func(bytesTransferred int64) {
		bytesSent = bytesTransferred
	})
	if err := req.SetBody(reqRpt, "application/octet-stream"); err != nil {
		t.Fatal(err)
	}
	_, err = pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bytesSent == 0 {
		t.Fatalf("bytesSent unexpectedly 0")
	}

	_, err = reqRpt.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, err := reqRpt.Read(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != contentSize {
		t.Fatalf("Seek did not reset Reader")
	}
}
