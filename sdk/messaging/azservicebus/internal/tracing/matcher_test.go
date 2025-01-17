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

// write unit test for matcher.go

func TestNewSpanValidator(t *testing.T) {
	hostName := "fake.servicebus.windows.net"
	provider := NewSpanValidator(t, SpanMatcher{
		Name:   "TestSpan",
		Kind:   SpanKindProducer,
		Status: SpanStatusUnset,
		Attributes: []Attribute{
			{Key: ServerAddress, Value: hostName},
		},
	})
	tracer := provider.NewTracer("module", "version")
	require.NotNil(t, tracer)
	require.True(t, tracer.Enabled())

	_, endSpan := runtime.StartSpan(context.Background(), "TestSpan", tracer, &runtime.StartSpanOptions{
		Kind: tracing.SpanKindProducer,
		Attributes: []Attribute{
			{Key: ServerAddress, Value: hostName},
		},
	})
	defer func() { endSpan(nil) }()
}

func TestMatchingTracerStart(t *testing.T) {
	hostName := "fake.servicebus.windows.net"
	matcher := SpanMatcher{
		Name:   "TestSpan",
		Kind:   SpanKindProducer,
		Status: SpanStatusUnset,
		Attributes: []Attribute{
			{Key: MessagingSystem, Value: "servicebus"},
			{Key: ServerAddress, Value: hostName},
		},
	}
	tracer := matchingTracer{
		matcher: matcher,
	}
	ctx := context.Background()
	// no-op when SpanName doesn't match
	_, spn := tracer.Start(ctx, "BadSpanName", SpanKindProducer, nil)
	require.EqualValues(t, spn, tracing.Span{})
	// tracks span when SpanName matches
	_, spn = tracer.Start(ctx, "TestSpan", SpanKindProducer, []Attribute{
		{Key: MessagingSystem, Value: "servicebus"},
		{Key: ServerAddress, Value: hostName},
	})
	require.NotNil(t, spn)
	spn.SetAttributes(tracing.Attribute{
		Key:   "TestAttributeKey",
		Value: "TestAttributeValue",
	})
	spn.SetStatus(SpanStatusOK, "ok")
}
