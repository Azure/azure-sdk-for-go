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

type SpanConfig struct {
	name       SpanName
	attributes []Attribute
}

type SetAttributesFn func([]Attribute) []Attribute

func NewSpanConfig(name SpanName, options ...SetAttributesFn) *SpanConfig {
	sc := &SpanConfig{name: name}
	for _, fn := range options {
		sc.attributes = fn(sc.attributes)
	}
	return sc
}

// StartSpan creates a span with the specified name and attributes.
// If no span name is provided, no span is created.
func StartSpan(ctx context.Context, tracer Tracer, sc *SpanConfig) (context.Context, func(error)) {
	if sc == nil || sc.name == "" {
		return ctx, func(error) {}
	}
	return runtime.StartSpan(ctx, string(sc.name), tracer, &StartSpanOptions{Attributes: sc.attributes})
}

type StartSpanOptions = runtime.StartSpanOptions
