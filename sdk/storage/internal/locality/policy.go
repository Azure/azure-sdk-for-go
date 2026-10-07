// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Package locality routes a request to a layout endpoint chosen by STG105 data locality.
package locality

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type ctxEndpointKey struct{}

// WithEndpoint returns a context that routes the requests made with it to endpoint, which may
// be a full URL (https://host:port) or a bare host. An empty endpoint returns ctx unchanged.
func WithEndpoint(ctx context.Context, endpoint string) context.Context {
	if endpoint == "" {
		return ctx
	}
	return context.WithValue(ctx, ctxEndpointKey{}, endpoint)
}

// Endpoint returns the layout endpoint attached to ctx by [WithEndpoint], if any.
func Endpoint(ctx context.Context) string {
	endpoint, _ := ctx.Value(ctxEndpointKey{}).(string)
	return endpoint
}

// NewPolicy returns a policy that sends a request to the layout endpoint attached to its
// context, if any, while keeping the account host in the Host header so the service
// authenticates and authorizes the request against the account. Requests without a layout
// endpoint pass through untouched.
//
// Register it per call, not per retry. The rewrite moves the account host into the Host header
// and must happen exactly once: every retry then reuses the rewritten request and stays on the
// layout endpoint, whereas repeating the rewrite would copy the layout host into the Host header.
func NewPolicy() policy.Policy {
	return endpointPolicy{}
}

type endpointPolicy struct{}

func (endpointPolicy) Do(req *policy.Request) (*http.Response, error) {
	if endpoint := Endpoint(req.Raw().Context()); endpoint != "" {
		raw := req.Raw()
		raw.Host = raw.URL.Host
		raw.URL.Host = hostFromEndpoint(endpoint)
	}
	return req.Next()
}

// hostFromEndpoint returns the host[:port] portion of endpoint, which may be a full URL or a
// bare host name.
func hostFromEndpoint(endpoint string) string {
	if i := strings.Index(endpoint, "://"); i >= 0 {
		endpoint = endpoint[i+len("://"):]
	}
	if i := strings.IndexAny(endpoint, "/?#"); i >= 0 {
		endpoint = endpoint[:i]
	}
	return endpoint
}
