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
	require.False(t, tr.Enabled())
	tr.SetAttributes()
	ctx, sp := tr.Start(context.Background(), "spanName", nil)
	require.Equal(t, context.Background(), ctx)
	require.Zero(t, sp)
	sp.AddEvent("event")
	sp.End()
	sp.SetAttributes(Attribute{})
	sp.SetStatus(SpanStatusError, "boom")
	spCtx := tr.SpanFromContext(ctx)
	require.Zero(t, spCtx)
}

func TestProvider(t *testing.T) {
	var addEventCalled bool
	var endCalled bool
	var setAttributesCalled bool
	var setStatusCalled bool
	var spanFromContextCalled bool

	pr := NewProvider(func(name, version string) Tracer {
		return NewTracer(func(context.Context, string, *SpanOptions) (context.Context, Span) {
			return nil, NewSpan(SpanImpl{
				AddEvent:      func(string, ...Attribute) { addEventCalled = true },
				End:           func() { endCalled = true },
				SetAttributes: func(...Attribute) { setAttributesCalled = true },
				SetStatus:     func(SpanStatus, string) { setStatusCalled = true },
			})
		}, &TracerOptions{
			SpanFromContext: func(context.Context) Span {
				spanFromContextCalled = true
				return Span{}
			},
		})
	}, nil)
	tr := pr.NewTracer("name", "version")
	require.NotZero(t, tr)
	require.True(t, tr.Enabled())
	sp := tr.SpanFromContext(context.Background())
	require.Zero(t, sp)
	tr.SetAttributes(Attribute{Key: "some", Value: "attribute"})
	require.Len(t, tr.attrs, 1)
	require.EqualValues(t, tr.attrs[0].Key, "some")
	require.EqualValues(t, tr.attrs[0].Value, "attribute")

	ctx, sp := tr.Start(context.Background(), "name", nil)
	require.NotEqual(t, context.Background(), ctx)
	require.NotZero(t, sp)

	sp.AddEvent("event")
	sp.End()
	sp.SetAttributes()
	sp.SetStatus(SpanStatusError, "desc")
	require.True(t, addEventCalled)
	require.True(t, endCalled)
	require.True(t, setAttributesCalled)
	require.True(t, setStatusCalled)
	require.True(t, spanFromContextCalled)
}

func TestWrapTracer(t *testing.T) {
	var attrs []Attribute
	tr := NewTracer(func(ctx context.Context, _ string, options *SpanOptions) (context.Context, Span) {
		attrs = options.Attributes
		return ctx, Span{}
	}, nil)
	tr.SetAttributes(Attribute{Key: "key", Value: "original"})
	original := tr
	decorate := func(inner Tracer) Tracer {
		wrapped := NewTracer(inner.Start, &TracerOptions{SpanFromContext: inner.SpanFromContext})
		wrapped.SetAttributes(Attribute{Key: "key", Value: "wrapper"})
		return wrapped
	}
	tr = WrapTracer(tr, decorate)
	tr.SetAttributes(Attribute{Key: "key", Value: "override"})
	tr = WrapTracer(tr, decorate)
	tr.Start(context.Background(), "wrapped", &SpanOptions{
		Attributes: []Attribute{{Key: "key", Value: "call"}},
	})
	require.Equal(t, []Attribute{
		{Key: "key", Value: "call"},
		{Key: "key", Value: "original"},
		{Key: "key", Value: "wrapper"},
		{Key: "key", Value: "override"},
		{Key: "key", Value: "wrapper"},
	}, attrs)
	original.Start(context.Background(), "original", nil)
	require.Equal(t, []Attribute{{Key: "key", Value: "original"}}, attrs)
}
