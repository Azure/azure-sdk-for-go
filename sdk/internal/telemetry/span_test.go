// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package telemetry

// unit test for span.go
import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

func TestWithSpan(t *testing.T) {
	tracer := NewSpanValidator(t, SpanMatcher{
		Name:   "TestSpan",
		Status: tracing.SpanStatusError,
	}).NewTracer("module", "version")
	require.NotNil(t, tracer)
	err := WithSpan(context.Background(), "TestSpan", tracer, nil, func(ctx context.Context) error {
		return errors.New("test error")
	})
	require.Error(t, err)

	tracer = NewSpanValidator(t, SpanMatcher{
		Name:   "TestSpan",
		Status: tracing.SpanStatusUnset,
	}).NewTracer("module", "version")
	require.NotNil(t, tracer)
	err = WithSpan(context.Background(), "TestSpan", tracer, nil, func(ctx context.Context) error {
		return nil
	})
	require.Nil(t, err)
}
