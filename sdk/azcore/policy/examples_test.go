//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package policy_test

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func ExampleWithCaptureResponse() {
	// policy.WithCaptureResponse provides a mechanism for obtaining an API's underlying *http.Response
	var respFromCtx *http.Response
	ctx := policy.WithCaptureResponse(context.TODO(), &respFromCtx)
	// make some client method call using the updated context
	// resp, err := client.SomeMethod(ctx, ...)
	// *respFromCtx contains the raw *http.Response returned during the client method call.
	// if the HTTP transport didn't return a response due to an error then *respFromCtx will be nil.
	_ = ctx
}

func ExampleWithHTTPHeader() {
	// policy.WithHTTPHeader allows callers to augment API calls with custom headers
	customHeaders := http.Header{}
	customHeaders.Add("key", "value")
	ctx := policy.WithHTTPHeader(context.TODO(), customHeaders)
	// make some client method call using the updated context
	// resp, err := client.SomeMethod(ctx, ...)
	// the underlying HTTP request will include the values in customHeaders
	_ = ctx
}

func ExampleWithRetryOptions() {
	// policy.WithRetryOptions contains a [policy.RetryOptions] that can be used to customize the retry policy on a per-API call basis
	customRetryOptions := policy.RetryOptions{ /* populate with custom values */ }
	ctx := policy.WithRetryOptions(context.TODO(), customRetryOptions)
	// make some client method call using the updated context
	// resp, err := client.SomeMethod(ctx, ...)
	// the behavior of the retry policy will correspond to the values provided in customRetryPolicy
	_ = ctx
}
