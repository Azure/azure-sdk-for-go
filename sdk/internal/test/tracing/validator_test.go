// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

func TestNewSpanValidator(t *testing.T) {
	provider := NewSpanValidator(t, SpanMatcher{
		Name:   "TestSpan",
		Kind:   tracing.SpanKindClient,
		Status: tracing.SpanStatusUnset,
		Attributes: []tracing.Attribute{
			{Key: "testKey", Value: "testValue"},
		},
	})
	tracer := provider.NewTracer("module", "version")
	require.NotNil(t, tracer)
	require.True(t, tracer.Enabled())

	ctx, endSpan := runtime.StartSpan(context.Background(), "TestSpan", tracer, &runtime.StartSpanOptions{
		Kind: tracing.SpanKindClient,
		Attributes: []tracing.Attribute{
			{Key: "testKey", Value: "testValue"},
		},
	})
	defer func() { endSpan(nil) }()

	require.NotNil(t, tracer.SpanFromContext(ctx))
	require.NotNil(t, tracer.LinkFromContext(ctx))
}

func TestMatchingTracerStart(t *testing.T) {
	matcher := SpanMatcher{
		Name:   "TestSpan",
		Kind:   tracing.SpanKindProducer,
		Status: tracing.SpanStatusUnset,
		Attributes: []tracing.Attribute{
			{Key: "testKey1", Value: "testValue1"},
			{Key: "testKey2", Value: "testValue2"},
		},
	}
	tracer := matchingTracer{
		matcher: matcher,
	}
	ctx := context.Background()
	// no-op when SpanName doesn't match
	_, spn := tracer.Start(ctx, "BadSpanName", tracing.SpanKindProducer, nil, nil)
	require.EqualValues(t, spn, tracing.Span{})
	// tracks span when SpanName matches
	_, spn = tracer.Start(ctx, "TestSpan", tracing.SpanKindProducer, []tracing.Attribute{
		{Key: "testKey1", Value: "testValue1"},
		{Key: "testKey2", Value: "testValue2"},
	}, nil)
	require.NotNil(t, spn)
	spn.SetAttributes(tracing.Attribute{
		Key:   "TestAttributeKey",
		Value: "TestAttributeValue",
	})
	spn.AddLink(tracing.Link{})
	spn.SetStatus(tracing.SpanStatusOK, "ok")
}
