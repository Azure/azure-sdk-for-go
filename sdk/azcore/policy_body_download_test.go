// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func newEmptyRequest(pol ...Policy) *Request {
	return &Request{
		Request: &http.Request{
			Header: http.Header{},
		},
		policies: pol,
	}
}

func newMockResponsePolicy(resp *http.Response) Policy {
	return PolicyFunc(func(ctx context.Context, req *Request) (*Response, error) {
		return &Response{Response: resp}, nil
	})
}

func newMockResponseWithBody(s string) *http.Response {
	return &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(s)),
	}
}

func TestDownloadBody(t *testing.T) {
	const message = "downloaded"
	req := newEmptyRequest(newMockResponsePolicy(newMockResponseWithBody(message)))
	p := newBodyDownloadPolicy()
	resp, err := p.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(resp.Payload) != message {
		t.Fatalf("unexpected response: %s", string(resp.Payload))
	}
}

func TestSkipBodyDownload(t *testing.T) {
	const message = "not downloaded"
	req := newEmptyRequest(newMockResponsePolicy(newMockResponseWithBody(message)))
	req.SkipBodyDownload()
	p := newBodyDownloadPolicy()
	resp, err := p.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Payload) > 0 {
		t.Fatalf("unexpected download: %s", string(resp.Payload))
	}
}
