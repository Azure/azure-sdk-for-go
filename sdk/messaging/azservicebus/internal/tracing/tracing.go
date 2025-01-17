// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

type Attribute = tracing.Attribute
type SpanKind = tracing.SpanKind
type Tracer = tracing.Tracer
type Provider = tracing.Provider

func NewNoOpTracer() *Tracer {
	return &Tracer{}
}

type TracerOptions struct {
	Tracer     Tracer
	SpanName   SpanName
	Attributes []Attribute
}

// StartSpan creates a span with the specified name and attributes.
// If no span name is provided, no span is created.
func StartSpan(ctx context.Context, options *TracerOptions) (context.Context, func(error)) {
	if options == nil || options.SpanName == "" {
		return ctx, func(error) {}
	}
	return runtime.StartSpan(ctx, string(options.SpanName), options.Tracer, &runtime.StartSpanOptions{Attributes: options.Attributes})
}
