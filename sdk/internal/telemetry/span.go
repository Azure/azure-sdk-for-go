// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package telemetry

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

type SpanCallbackFn = func(context.Context) error

// WithSpan creates a span and executes the provided function with the span's context.
// The span is ended with the error returned from the function.
func WithSpan(ctx context.Context, spanName string, tracer tracing.Tracer, options *runtime.StartSpanOptions, fn SpanCallbackFn) error {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, spanName, tracer, options)
	defer func() { endSpan(err) }()
	err = fn(ctx)
	return err
}
