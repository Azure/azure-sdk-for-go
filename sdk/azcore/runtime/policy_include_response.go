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

// WithCaptureResponse applies the HTTP response retrieval annotation to the parent context.
// Call the returned function to retrieve the HTTP response after the request has completed.
// The created context is NOT safe to use across multiple goroutines.
func WithCaptureResponse(parent context.Context) (context.Context, func() *http.Response) {
	var resp *http.Response
	return context.WithValue(parent, shared.CtxIncludeResponseKey{}, &resp), func() *http.Response {
		cap := &resp
		return *cap
	}
}
