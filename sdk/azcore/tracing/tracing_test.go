//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProviderZeroValues(t *testing.T) {
	pr := Provider{}
	tr := pr.NewTracer("name", "version")
	require.Zero(t, tr)
	ctx, sp := tr.Start(context.Background(), "spanName", nil)
	require.Equal(t, context.Background(), ctx)
	require.Zero(t, sp)
	sp.AddError(nil)
	sp.AddEvent("event")
	sp.End()
	sp.SetAttributes(Attribute{})
	sp.SetStatus(SpanStatusError, "boom")
}

func TestProvider(t *testing.T) {
	var addErrorCalled bool
	var addEventCalled bool
	var endCalled bool
	var setAttributesCalled bool
	var setStatusCalled bool

	pr := NewProvider(func(name, version string) Tracer {
		return NewTracer(func(context.Context, string, *SpanOptions) (context.Context, Span) {
			return nil, NewSpan(SpanImpl{
				AddError:      func(error) { addErrorCalled = true },
				AddEvent:      func(string, ...Attribute) { addEventCalled = true },
				End:           func() { endCalled = true },
				SetAttributes: func(...Attribute) { setAttributesCalled = true },
				SetStatus:     func(SpanStatus, string) { setStatusCalled = true },
			})
		}, nil)
	}, nil)
	tr := pr.NewTracer("name", "version")
	require.NotZero(t, tr)

	ctx, sp := tr.Start(context.Background(), "name", nil)
	require.NotEqual(t, context.Background(), ctx)
	require.NotZero(t, sp)

	sp.AddError(nil)
	sp.AddEvent("event")
	sp.End()
	sp.SetAttributes()
	sp.SetStatus(SpanStatusError, "desc")
	require.True(t, addErrorCalled)
	require.True(t, addEventCalled)
	require.True(t, endCalled)
	require.True(t, setAttributesCalled)
	require.True(t, setStatusCalled)
}
