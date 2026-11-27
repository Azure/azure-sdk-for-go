// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package transport

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// DefaultHTTPClient is the default configured http.Client used by all clients.
var DefaultHTTPClient *http.Client

func init() {
	defaultTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: defaultTransportDialContext(&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}),
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion:    tls.VersionTLS12,
			Renegotiation: tls.RenegotiateFreelyAsClient,
		},
		HTTP2: &http.HTTP2Config{
			// if the connection has been idle for 10 seconds, send a ping frame for a health check
			SendPingTimeout: 10 * time.Second,
			// if there's no response to the ping within the timeout, the connection will be closed
			PingTimeout: 5 * time.Second,
		},
	}
	DefaultHTTPClient = &http.Client{
		Transport: defaultTransport,
	}
}
