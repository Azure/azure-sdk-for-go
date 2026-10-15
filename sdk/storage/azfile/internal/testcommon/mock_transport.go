// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testcommon

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// MockResponse is a canned HTTP response returned by MockTransport.
type MockResponse struct {
	StatusCode int
	// Header, if set, is added to the response headers.
	Header http.Header
	// ErrorCode, if set, is returned in the x-ms-error-code header.
	ErrorCode string
	Body      string
}

// MockTransport is a policy.Transporter that returns Responses in order and records the requests it receives.
type MockTransport struct {
	Responses []MockResponse
	Requests  []*http.Request
}

// NewMockTransport creates a MockTransport that returns the given responses in order.
func NewMockTransport(responses ...MockResponse) *MockTransport {
	return &MockTransport{Responses: responses}
}

// Do implements the policy.Transporter interface.
func (m *MockTransport) Do(req *http.Request) (*http.Response, error) {
	if len(m.Requests) >= len(m.Responses) {
		return nil, fmt.Errorf("unexpected request %d to %s: only %d responses configured", len(m.Requests)+1, req.URL, len(m.Responses))
	}
	r := m.Responses[len(m.Requests)]
	m.Requests = append(m.Requests, req)
	header := http.Header{"Content-Type": []string{"application/xml"}}
	for k, v := range r.Header {
		header[k] = v
	}
	if r.ErrorCode != "" {
		header.Set("x-ms-error-code", r.ErrorCode)
	}
	return &http.Response{
		StatusCode: r.StatusCode,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(r.Body)),
		Request:    req,
	}, nil
}

// MockClientOptions returns azcore.ClientOptions that send requests to transport with retries disabled.
func MockClientOptions(transport *MockTransport) azcore.ClientOptions {
	return azcore.ClientOptions{
		Transport: transport,
		Retry:     policy.RetryOptions{MaxRetries: -1},
	}
}
