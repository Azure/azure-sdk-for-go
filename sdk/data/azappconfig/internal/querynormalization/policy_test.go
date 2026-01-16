//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package querynormalization

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestQueryNormalizationPolicy(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		queryParams map[string][]string
		expectedURL string
	}{
		{
			name: "Normalizes query parameters to lowercase and sorts alphabetically",
			path: "/kv",
			queryParams: map[string][]string{
				"api-version": {"2023-11-01"},
				"After":       {"abcdefg"},
				"tags":        {"tag3=value3", "tag2=value2", "tag1=value1"},
				"key":         {"*"},
				"label":       {"dev"},
				"$Select":     {"key"},
			},
			expectedURL: "/kv?%24select=key&after=abcdefg&api-version=2023-11-01&key=%2A&label=dev&tags=tag3%3Dvalue3&tags=tag2%3Dvalue2&tags=tag1%3Dvalue1",
		},
		{
			name: "Keeps original order of duplicate query parameters",
			path: "/kv",
			queryParams: map[string][]string{
				"tags":        {"tag2", "tag1"},
				"api-version": {"2023-11-01"},
			},
			expectedURL: "/kv?api-version=2023-11-01&tags=tag2&tags=tag1",
		},
		{
			name: "Handles complex query with mixed case parameters",
			path: "/kv",
			queryParams: map[string][]string{
				"Label":       {"prod"},
				"KEY":         {"mykey"},
				"Api-Version": {"2023-11-01"},
			},
			expectedURL: "/kv?api-version=2023-11-01&key=mykey&label=prod",
		},
		{
			name: "Handles special characters in parameter names",
			path: "/kv",
			queryParams: map[string][]string{
				"$filter":     {"value"},
				"api-version": {"2023-11-01"},
				"$select":     {"key"},
			},
			expectedURL: "/kv?%24filter=value&%24select=key&api-version=2023-11-01",
		},
		{
			name: "Handles multiple values for same parameter",
			path: "/kv",
			queryParams: map[string][]string{
				"tags":        {"tag1", "tag2", "tag3"},
				"key":         {"*"},
				"api-version": {"2023-11-01"},
			},
			expectedURL: "/kv?api-version=2023-11-01&key=%2A&tags=tag1&tags=tag2&tags=tag3",
		},
		{
			name: "Handles continuation token (after parameter)",
			path: "/kv",
			queryParams: map[string][]string{
				"key":         {"mykey"},
				"After":       {"token123"},
				"api-version": {"2023-11-01"},
			},
			expectedURL: "/kv?after=token123&api-version=2023-11-01&key=mykey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, close := mock.NewServer()
			defer close()

			var capturedURL string
			pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
				PerRetry: []policy.Policy{NewPolicy()},
			}, &policy.ClientOptions{
				Transport: &urlCapturingTransporter{
					real: srv,
					onRequest: func(req *http.Request) {
						capturedURL = req.URL.String()
					},
				},
			})

			// Build request with query parameters
			req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL()+tt.path)
			require.NoError(t, err)

			// Add query parameters
			q := req.Raw().URL.Query()
			for name, values := range tt.queryParams {
				for _, value := range values {
					q.Add(name, value)
				}
			}
			req.Raw().URL.RawQuery = q.Encode()

			srv.AppendResponse()
			resp, err := pl.Do(req)
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Verify the URL was normalized correctly
			require.Contains(t, capturedURL, tt.expectedURL, "URLs do not match")
		})
	}
}

func TestQueryNormalizationPolicyNoQueryParams(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()

	var capturedURL string
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{NewPolicy()},
	}, &policy.ClientOptions{
		Transport: &urlCapturingTransporter{
			real: srv,
			onRequest: func(req *http.Request) {
				capturedURL = req.URL.String()
			},
		},
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL()+"/kv")
	require.NoError(t, err)

	srv.AppendResponse()
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// URL should not have query parameters
	require.NotContains(t, capturedURL, "?")
}

func TestQueryNormalizationPolicyStability(t *testing.T) {
	// Test that the policy maintains stable sorting for duplicate parameter names
	srv, close := mock.NewServer()
	defer close()

	var capturedURL string
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{NewPolicy()},
	}, &policy.ClientOptions{
		Transport: &urlCapturingTransporter{
			real: srv,
			onRequest: func(req *http.Request) {
				capturedURL = req.URL.String()
			},
		},
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL()+"/kv")
	require.NoError(t, err)

	// Manually set query string to test order preservation
	req.Raw().URL.RawQuery = "b=2&a=1&b=3&a=2&b=1&a=3"

	srv.AppendResponse()
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify stable sort: parameters are sorted by name, but duplicates keep their original order
	expectedURL := "/kv?a=1&a=2&a=3&b=2&b=3&b=1"
	require.Contains(t, capturedURL, expectedURL, "URLs do not match - stable sort not maintained")
}

func TestQueryNormalizationPolicyEmptyItem(t *testing.T) {
	// Test that the policy handles empty query parameters correctly
	srv, close := mock.NewServer()
	defer close()

	var capturedURL string
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{NewPolicy()},
	}, &policy.ClientOptions{
		Transport: &urlCapturingTransporter{
			real: srv,
			onRequest: func(req *http.Request) {
				capturedURL = req.URL.String()
			},
		},
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL()+"/kv")
	require.NoError(t, err)

	// Manually set query string to test order preservation
	req.Raw().URL.RawQuery = "b=2&&a=1"

	srv.AppendResponse()
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify stable sort: parameters are sorted by name, but duplicates keep their original order
	expectedURL := "/kv?&a=1&b=2"
	require.Contains(t, capturedURL, expectedURL, "URLs do not match - stable sort not maintained")
}

type urlCapturingTransporter struct {
	onRequest func(*http.Request)
	real      policy.Transporter
}

func (t *urlCapturingTransporter) Do(req *http.Request) (*http.Response, error) {
	if t.onRequest != nil {
		t.onRequest(req)
	}
	return t.real.Do(req)
}
