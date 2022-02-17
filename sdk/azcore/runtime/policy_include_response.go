//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// includeResponsePolicy creates a policy that retrieves the raw HTTP response upon request
func includeResponsePolicy(req *policy.Request) (*http.Response, error) {
	resp, err := req.Next()
	if resp == nil {
		return resp, err
	}
	if httpOutRaw := req.Raw().Context().Value(shared.CtxIncludeResponseKey{}); httpOutRaw != nil {
		httpOut := httpOutRaw.(**http.Response)
		*httpOut = resp
	}
	return resp, err
}

// IncludeResponse applies the HTTP response retrieval annotation to the parent context.
// Call runtime.ResponseFromContext() to retrieve the HTTP response from the context.
func IncludeResponse(parent context.Context) context.Context {
	var rawResp *http.Response
	return context.WithValue(parent, shared.CtxIncludeResponseKey{}, &rawResp)
}

// ResponseFromContext retrieves the raw HTTP response from the specified context.
// Disabled by default.  Use policy.IncludeResponse() to enable on the calling context.
func ResponseFromContext(ctx context.Context) *http.Response {
	if httpOutRaw := ctx.Value(shared.CtxIncludeResponseKey{}); httpOutRaw != nil {
		httpOut := httpOutRaw.(**http.Response)
		return *httpOut
	}
	return nil
}
