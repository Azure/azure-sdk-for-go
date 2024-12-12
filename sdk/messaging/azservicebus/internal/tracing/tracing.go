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

func NewNoOpTracer() Tracer {
	return Tracer{}
}

type SpanOptions struct {
	name       MessagingSpanName
	attributes []Attribute
}

type SetAttributesFn func([]Attribute) []Attribute

func NewSpanOptions(name MessagingSpanName, options ...SetAttributesFn) *SpanOptions {
	so := &SpanOptions{name: name}
	for _, fn := range options {
		so.attributes = fn(so.attributes)
	}
	return so
}

// StartSpan creates a span with the specified name and attributes.
// If no span name is provided, no span is created.
func StartSpan(ctx context.Context, tracer Tracer, so *SpanOptions) (context.Context, func(error)) {
	if so == nil || so.name == "" {
		return ctx, func(error) {}
	}
	return runtime.StartSpan(ctx, string(so.name), tracer, &StartSpanOptions{Attributes: so.attributes})
}

type StartSpanOptions = runtime.StartSpanOptions
