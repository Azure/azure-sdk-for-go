// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package telemetry

// unit test for matcher.go
import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

func TestNewSpanValidator(t *testing.T) {
	// validates a span with no attributes
	matcher := SpanMatcher{
		Name:   "TestSpan",
		Status: tracing.SpanStatusError,
	}
	provider := NewSpanValidator(t, matcher)
	tracer := provider.NewTracer("module", "version")
	require.NotNil(t, tracer)
	ctx, endSpan := runtime.StartSpan(context.Background(), "TestSpan", tracer, nil)
	defer endSpan(errors.New("test error"))
	require.NotNil(t, ctx)

	// does not track a span with a different name
	ctx, endSpan = runtime.StartSpan(context.Background(), "AnotherSpan", tracer, nil)
	defer endSpan(errors.New("test error"))
	require.NotNil(t, ctx)

	// validates when attributes are provided
	matcher = SpanMatcher{
		Name:   "TestSpan",
		Status: tracing.SpanStatusUnset,
		Attributes: []tracing.Attribute{
			{Key: "testKey", Value: "testValue"},
		},
	}
	provider = NewSpanValidator(t, matcher)
	tracer = provider.NewTracer("module", "version")
	require.NotNil(t, tracer)
	_, endSpan = runtime.StartSpan(context.Background(), "TestSpan", tracer, &runtime.StartSpanOptions{
		Attributes: []tracing.Attribute{
			{Key: "testKey", Value: "testValue"},
		},
	})
	defer endSpan(nil)
	require.NotNil(t, ctx)
}
