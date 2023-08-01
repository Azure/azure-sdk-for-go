//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime_test

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func ExampleWithCaptureResponse() {
	var respFromCtx *http.Response
	ctx := context.TODO()
	ctx = runtime.WithCaptureResponse(ctx, &respFromCtx)
	// make some client method call using the updated context
	// resp, err := client.SomeMethod(ctx, ...)
	// *respFromCtx contains the raw *http.Response returned during the client method call.
	// if the HTTP transport didn't return a response due to an error then *respFromCtx will be nil.
	_ = ctx
}
