// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"syscall"
)

var ErrPoolExhausted = errors.New("connection pool exhausted")

// IsWebExceptionRetriable checks if an error is retriable at the web/transport level.
// errors.As walks the full error chain, so no manual unwrapping is needed.
func IsWebExceptionRetriable(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrPoolExhausted) {
		return true
	}

	if isConnectError(err) {
		return true
	}

	if isUnknownHostError(err) {
		return true
	}

	if isNoRouteToHostError(err) {
		return true
	}

	if isSSLHandshakeError(err) {
		return true
	}

	return false
}

// IsNetworkFailure checks if an error represents a network-level failure.
// errors.As walks the full error chain, so no manual unwrapping is needed.
func IsNetworkFailure(err error) bool {
	if err == nil {
		return false
	}

	if isClosedChannelError(err) {
		return true
	}

	if isSocketError(err) {
		return true
	}

	if isSSLError(err) {
		return true
	}

	if isUnknownHostError(err) {
		return true
	}

	if isConnectError(err) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return true
		}
	}

	if isNoRouteToHostError(err) {
		return true
	}

	return false
}

func isConnectError(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if opErr.Op == "dial" || opErr.Op == "connect" {
			return true
		}
		if errors.Is(opErr.Err, syscall.ECONNREFUSED) {
			return true
		}
	}
	return false
}

func isUnknownHostError(err error) bool {
	var dnsErr *net.DNSError
	return errors.As(err, &dnsErr)
}

func isNoRouteToHostError(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if errors.Is(opErr.Err, syscall.EHOSTUNREACH) {
			return true
		}
		if errors.Is(opErr.Err, syscall.ENETUNREACH) {
			return true
		}
	}
	return false
}

func isSSLHandshakeError(err error) bool {
	var certErr x509.CertificateInvalidError
	if errors.As(err, &certErr) {
		return true
	}

	var hostErr x509.HostnameError
	if errors.As(err, &hostErr) {
		return true
	}

	var unknownAuthErr x509.UnknownAuthorityError
	if errors.As(err, &unknownAuthErr) {
		return true
	}

	return false
}

func isSSLError(err error) bool {
	if isSSLHandshakeError(err) {
		return true
	}

	var tlsErr tls.RecordHeaderError
	if errors.As(err, &tlsErr) {
		return true
	}

	var alertErr tls.AlertError
	if errors.As(err, &alertErr) {
		return true
	}

	return false
}

func isClosedChannelError(err error) bool {
	if errors.Is(err, io.ErrClosedPipe) {
		return true
	}
	if errors.Is(err, net.ErrClosed) {
		return true
	}
	return false
}

func isSocketError(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}
	return false
}
