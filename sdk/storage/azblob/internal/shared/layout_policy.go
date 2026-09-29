// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type CtxLayoutEndpointKey struct{}

func WithLayoutEndpoint(ctx context.Context, endpoint string) context.Context {
	if endpoint == "" {
		return ctx
	}
	return context.WithValue(ctx, CtxLayoutEndpointKey{}, endpoint)
}

type LayoutPolicy struct {
}

func (l LayoutPolicy) Do(req *policy.Request) (*http.Response, error) {
	// Check if the layout endpoint is set in the context
	if layoutEndpoint := req.Raw().Context().Value(CtxLayoutEndpointKey{}); layoutEndpoint != nil && layoutEndpoint != "" {
		// Read the request endpoint (account) and set the Host header to the endpoint if not already set.
		req.Raw().Host = req.Raw().URL.Host

		// The layout endpoint may be a full URL (e.g. https://layout.blob.core.windows.net)
		// or just a host name (e.g. layout.blob.core.windows.net), so extract the host from it.
		// Set the request URL to the layout endpoint
		req.Raw().URL.Host = hostFromEndpoint(layoutEndpoint.(string))
	}
	return req.Next()
}

// hostFromEndpoint returns the host portion of endpoint, which may be a full URL or a bare host name.
func hostFromEndpoint(endpoint string) string {
	if i := strings.Index(endpoint, "://"); i >= 0 {
		endpoint = endpoint[i+len("://"):]
	}
	// drop any path, query or fragment
	if i := strings.IndexAny(endpoint, "/?#"); i >= 0 {
		endpoint = endpoint[:i]
	}
	return endpoint
}

func NewLayoutPolicy() policy.Policy {
	return LayoutPolicy{}
}
