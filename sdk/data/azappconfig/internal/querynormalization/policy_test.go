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

	// Manually set query string with empty parameter between b=2 and a=1
	req.Raw().URL.RawQuery = "b=2&&a=1"

	srv.AppendResponse()
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Empty entries should be filtered out, not preserved as consecutive '&' characters
	expectedURL := "/kv?a=1&b=2"
	require.Contains(t, capturedURL, expectedURL, "URLs do not match")
}

func TestQueryNormalizationPolicyURLEncodedChars(t *testing.T) {
	// Test that URL-encoded characters in parameter names and values are handled correctly
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

	// Set query with URL-encoded characters
	req.Raw().URL.RawQuery = "%24Select=key&api-version=2023-11-01"

	srv.AppendResponse()
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// $Select should be decoded, lowercased to $select, then re-encoded as %24select
	expectedURL := "/kv?%24select=key&api-version=2023-11-01"
	require.Contains(t, capturedURL, expectedURL, "URL-encoded characters not handled correctly")
}

func TestQueryNormalizationPolicyCaseCollision(t *testing.T) {
	// Test that different-cased keys that collide after lowercasing preserve
	// their original positional order. This is critical for HMAC signing
	// determinism: the same raw query must always produce the same output.
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

	// A=1 and a=2 collide after lowercasing; positional order must be preserved
	req.Raw().URL.RawQuery = "A=1&a=2"

	srv.AppendResponse()
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// After lowercasing both become "a"; original order A=1 (index 0) then a=2 (index 1)
	expectedURL := "/kv?a=1&a=2"
	require.Contains(t, capturedURL, expectedURL, "case-collision ordering not deterministic")

	// Run the same input multiple times to verify determinism
	for i := 0; i < 20; i++ {
		req, err = runtime.NewRequest(context.Background(), http.MethodGet, srv.URL()+"/kv")
		require.NoError(t, err)
		req.Raw().URL.RawQuery = "A=1&a=2"

		srv.AppendResponse()
		resp, err = pl.Do(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Contains(t, capturedURL, expectedURL, "case-collision ordering not deterministic on iteration %d", i)
	}
}

func TestQueryNormalizationPolicyCaseCollisionEncoded(t *testing.T) {
	// Test that URL-encoded keys that collide after decoding and lowercasing
	// preserve their original positional order.
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

	// %24Select and %24select collide after decoding ($Select/$select) and lowercasing
	req.Raw().URL.RawQuery = "%24Select=key&%24select=value&api-version=2023-11-01"

	srv.AppendResponse()
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Both become %24select after normalization; positional order preserved
	expectedURL := "/kv?%24select=key&%24select=value&api-version=2023-11-01"
	require.Contains(t, capturedURL, expectedURL, "encoded case-collision ordering not deterministic")
}

func TestQueryNormalizationPolicyCaseCollisionInterleaved(t *testing.T) {
	// Test interleaved different-cased keys that collide after lowercasing,
	// mixed with other keys.
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

	// B=x, A=1, b=y, a=2: after lowercasing, a and b each have two entries
	// positional order: B=x(0), A=1(1), b=y(2), a=2(3)
	// sorted by name: a group [A=1(1), a=2(3)], b group [B=x(0), b=y(2)]
	req.Raw().URL.RawQuery = "B=x&A=1&b=y&a=2"

	srv.AppendResponse()
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	expectedURL := "/kv?a=1&a=2&b=x&b=y"
	require.Contains(t, capturedURL, expectedURL, "interleaved case-collision ordering not deterministic")
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
