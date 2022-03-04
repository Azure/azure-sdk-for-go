//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"net/http"
)

func ExampleWithCaptureResponse() {
	var respFromCtx *http.Response
	WithCaptureResponse(context.Background(), &respFromCtx)
}
