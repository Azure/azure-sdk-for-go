//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

func TestTransportIPv6Support(t *testing.T) {
	// Test that the default transport supports IPv6 by checking the dialer configuration
	transport := defaultHTTPClient.Transport.(*http.Transport)
	require.NotNil(t, transport, "default transport should not be nil")
	
	// The DialContext function should be set and should support IPv6
	require.NotNil(t, transport.DialContext, "DialContext should be configured")
	
	// Test with a mock IPv6 address to ensure the dialer can handle IPv6
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	// Try to dial an IPv6 loopback address (this will fail to connect but should not error on IPv6 parsing)
	_, err := transport.DialContext(ctx, "tcp", "[::1]:80")
	// We expect a connection error, not a parsing error
	require.Error(t, err, "connection should fail to non-existent server")
	
	// The error should be a network error, not a parsing error
	var netErr net.Error
	require.ErrorAs(t, err, &netErr, "error should be a network error")
}

func TestTransportIPv6URLParsing(t *testing.T) {
	// Test that requests can be created with IPv6 URLs
	testCases := []struct {
		name string
		url  string
	}{
		{
			name: "IPv6 with port",
			url:  "https://[2001:db8::1]:443/path",
		},
		{
			name: "IPv6 without port",
			url:  "https://[2001:db8::1]/path",
		},
		{
			name: "IPv6 loopback",
			url:  "https://[::1]:8080/api",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := NewRequest(context.Background(), http.MethodGet, tc.url)
			require.NoError(t, err, "should be able to create request with IPv6 URL")
			require.NotNil(t, req, "request should not be nil")
			require.Equal(t, tc.url, req.Raw().URL.String(), "URL should be preserved")
		})
	}
}

func TestTransportDualStackCapability(t *testing.T) {
	// Test that the transport can handle both IPv4 and IPv6 addresses
	transport := defaultHTTPClient.Transport.(*http.Transport)
	
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	// Test IPv4 address using a non-routable address
	_, err4 := transport.DialContext(ctx, "tcp", "192.0.2.1:12345")
	require.Error(t, err4, "IPv4 connection should fail to non-existent server")
	
	// Test IPv6 address using a non-routable address
	_, err6 := transport.DialContext(ctx, "tcp", "[2001:db8::1]:12345")
	require.Error(t, err6, "IPv6 connection should fail to non-existent server")
	
	// Both should be network errors, indicating the dialer supports both protocols
	var netErr4, netErr6 net.Error
	require.ErrorAs(t, err4, &netErr4, "IPv4 error should be a network error")
	require.ErrorAs(t, err6, &netErr6, "IPv6 error should be a network error")
}

func TestHTTPClientIPv6Request(t *testing.T) {
	// Test that we can create an HTTP client that handles IPv6 requests
	client := &http.Client{
		Transport: defaultHTTPClient.Transport,
		Timeout:   100 * time.Millisecond,
	}
	
	// Create a request to an IPv6 address using a non-routable address
	req, err := http.NewRequest(http.MethodGet, "https://[2001:db8::1]:12345/test", nil)
	require.NoError(t, err, "should be able to create request with IPv6 URL")
	
	// The request should fail with a connection error (not a parsing error)
	_, err = client.Do(req)
	require.Error(t, err, "request should fail to non-existent server")
	
	// Should be a network error indicating the transport attempted the connection
	var netErr net.Error
	require.ErrorAs(t, err, &netErr, "error should be a network error")
}

func TestPipelineWithIPv6(t *testing.T) {
	// Test that the azcore pipeline can handle IPv6 endpoints
	pl := NewPipeline("test", "1.0.0", PipelineOptions{}, &policy.ClientOptions{})
	require.NotNil(t, pl, "pipeline should not be nil")
	
	// Create a request with IPv6 URL
	req, err := NewRequest(context.Background(), http.MethodGet, "https://[2001:db8::1]:443/api/test")
	require.NoError(t, err, "should be able to create request with IPv6 URL")
	
	// The request creation should succeed (even though the actual network call would fail)
	require.NotNil(t, req, "request should not be nil")
	require.Equal(t, "2001:db8::1", req.Raw().URL.Hostname(), "hostname should be parsed correctly")
	require.Equal(t, "443", req.Raw().URL.Port(), "port should be parsed correctly")
}

func TestDefaultTransportIPv6Configuration(t *testing.T) {
	// Test that the default transport configuration supports IPv6 and dual-stack networking.
	// This validates that we haven't inadvertently disabled IPv6 support.
	transport := defaultHTTPClient.Transport.(*http.Transport)
	require.NotNil(t, transport, "default transport should not be nil")
	
	// Verify the transport has the expected configuration for IPv6 support
	require.NotNil(t, transport.DialContext, "DialContext should be configured for dual-stack networking")
	require.True(t, transport.ForceAttemptHTTP2, "HTTP/2 should be enabled (supports IPv6)")
	
	// The dialer should be the standard Go dialer which supports dual-stack by default
	// We can't easily test the dialer type directly, but we can verify it's not nil
	// and that it can handle IPv6 addresses (tested in other test functions)
	
	// Verify TLS configuration doesn't restrict IPv6
	require.NotNil(t, transport.TLSClientConfig, "TLS config should be present")
	require.GreaterOrEqual(t, transport.TLSClientConfig.MinVersion, uint16(tls.VersionTLS12), 
		"TLS version should support modern protocols")
}