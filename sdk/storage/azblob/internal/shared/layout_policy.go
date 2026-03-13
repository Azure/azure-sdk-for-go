// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"context"
	"net/http"

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

		// Set the request URL to the layout endpoint
		req.Raw().URL.Host = layoutEndpoint.(string)
	}
	return req.Next()
}

func NewLayoutPolicy() policy.Policy {
	return LayoutPolicy{}
}
