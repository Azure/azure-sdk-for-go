// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsWebExceptionRetriable(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"RuntimeError", errors.New("runtime error"), false},
		{"ConnectError", &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}, true},
		{"ConnectTimeoutError", &net.OpError{Op: "dial", Err: &timeoutError{}}, true},
		{"UnknownHostError", &net.DNSError{IsNotFound: true, Name: "unknown.host"}, true},
		{"ReadTimeoutError", &timeoutError{}, false},
		{"SSLHandshakeError", x509.HostnameError{Host: "example.com"}, true},
		{"NoRouteToHostError", &net.OpError{Op: "dial", Err: syscall.EHOSTUNREACH}, true},
		{"SSLPeerUnverifiedError", x509.CertificateInvalidError{}, true},
		{"SocketTimeoutError", &timeoutError{}, false},
		{"PoolExhaustedError", ErrPoolExhausted, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsWebExceptionRetriable(tc.err)
			require.Equal(t, tc.expected, result, "IsWebExceptionRetriable(%v) = %v, want %v", tc.err, result, tc.expected)
		})
	}
}

func TestIsNetworkFailure(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"RuntimeError", errors.New("runtime error"), false},
		{"ConnectError", &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}, true},
		{"ConnectTimeoutError", &net.OpError{Op: "dial", Err: &timeoutError{}}, true},
		{"UnknownHostError", &net.DNSError{IsNotFound: true, Name: "unknown.host"}, true},
		{"ReadTimeoutError", &timeoutError{}, true},
		{"SSLHandshakeError", x509.HostnameError{Host: "example.com"}, true},
		{"NoRouteToHostError", &net.OpError{Op: "dial", Err: syscall.EHOSTUNREACH}, true},
		{"SSLPeerUnverifiedError", x509.CertificateInvalidError{}, true},
		{"SocketTimeoutError", &timeoutError{}, true},
		{"ChannelException", &net.OpError{Op: "read", Err: errors.New("connection reset")}, true},
		{"PlainIOError", io.EOF, false},
		{"ClosedChannelError", net.ErrClosed, true},
		{"SocketError", &net.OpError{Op: "write", Err: errors.New("broken pipe")}, true},
		{"UnknownServiceError", &net.OpError{Op: "dial", Net: "unknown", Err: errors.New("unknown service")}, true},
		{"InterruptedByTimeoutError", &timeoutError{}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsNetworkFailure(tc.err)
			require.Equal(t, tc.expected, result, "IsNetworkFailure(%v) = %v, want %v", tc.err, result, tc.expected)
		})
	}
}

func TestIsWebExceptionRetriable_WrappedErrors(t *testing.T) {
	dnsErr := &net.DNSError{IsNotFound: true, Name: "unknown.host"}
	wrappedErr := fmt.Errorf("connection failed: %w", dnsErr)

	require.True(t, IsWebExceptionRetriable(wrappedErr), "should find wrapped DNS error")
}

func TestIsNetworkFailure_WrappedErrors(t *testing.T) {
	opErr := &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}
	wrappedErr := fmt.Errorf("connection failed: %w", opErr)

	require.True(t, IsNetworkFailure(wrappedErr), "should find wrapped OpError")
}

func TestIsWebExceptionRetriable_NilError(t *testing.T) {
	require.False(t, IsWebExceptionRetriable(nil))
}

func TestIsNetworkFailure_NilError(t *testing.T) {
	require.False(t, IsNetworkFailure(nil))
}

func TestIsNetworkFailure_ClosedPipe(t *testing.T) {
	require.True(t, IsNetworkFailure(io.ErrClosedPipe))
}

func TestIsWebExceptionRetriable_PoolExhausted(t *testing.T) {
	require.True(t, IsWebExceptionRetriable(ErrPoolExhausted))

	wrappedPoolErr := fmt.Errorf("failed to get connection: %w", ErrPoolExhausted)
	require.True(t, IsWebExceptionRetriable(wrappedPoolErr))
}

type timeoutError struct{}

func (e *timeoutError) Error() string   { return "timeout" }
func (e *timeoutError) Timeout() bool   { return true }
func (e *timeoutError) Temporary() bool { return true }
