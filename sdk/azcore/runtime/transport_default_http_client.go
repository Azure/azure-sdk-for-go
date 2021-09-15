//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

var defaultHTTPClient *http.Client

func init() {
	defaultTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	defaultHTTPClient = &http.Client{
		Transport: defaultTransport,
	}
}

// AuthenticationOptions contains various options used to create a credential policy.
type AuthenticationOptions struct {
	// TokenRequest is a TokenRequestOptions that includes a scopes field which contains
	// the list of OAuth2 authentication scopes used when requesting a token.
	// This field is ignored for other forms of authentication (e.g. shared key).
	TokenRequest policy.TokenRequestOptions
	// AuxiliaryTenants contains a list of additional tenant IDs to be used to authenticate
	// in cross-tenant applications.
	AuxiliaryTenants []string
}
