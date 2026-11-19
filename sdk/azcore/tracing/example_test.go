// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing_test

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

func ExampleWrapTracer() {
	tracer := tracing.NewTracer(func(ctx context.Context, _ string, _ *tracing.SpanOptions) (context.Context, tracing.Span) {
		return ctx, tracing.Span{}
	}, nil)
	tracer.SetAttributes(tracing.Attribute{Key: "az.namespace", Value: "original"})
	tracer = tracing.WrapTracer(tracer, func(inner tracing.Tracer) tracing.Tracer {
		return tracing.NewTracer(func(ctx context.Context, name string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
			ctx, span := inner.Start(ctx, name, options)
			// Customize the span before returning it to the caller.
			span.SetAttributes(tracing.Attribute{Key: "custom", Value: true})
			return ctx, span
		}, &tracing.TracerOptions{SpanFromContext: inner.SpanFromContext})
	})
	tracer.SetAttributes(tracing.Attribute{Key: "az.namespace", Value: "override"})
	_, span := tracer.Start(context.TODO(), "operation", nil)
	defer span.End()
}
