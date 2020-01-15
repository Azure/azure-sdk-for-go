// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

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
	pl := NewPipeline(srv, NewTelemetryPolicy(TelemetryOptions{}))
	req := NewRequest(http.MethodGet, srv.URL())
	req.SkipBodyDownload()
	var bytesSent int64
	reqRpt := NewRequestBodyProgress(NopCloser(body), func(bytesTransferred int64) {
		bytesSent = bytesTransferred
	})
	req.SetBody(reqRpt)
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var bytesReceived int64
	respRpt := NewResponseBodyProgress(resp.Body, func(bytesTransferred int64) {
		bytesReceived = bytesTransferred
	})
	defer respRpt.Close()
	b, err := ioutil.ReadAll(respRpt)
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
