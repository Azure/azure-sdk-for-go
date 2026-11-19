// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package sanitizer

import (
	"context"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/stretchr/testify/require"
)

func TestWrapTracerDisabled(t *testing.T) {
	tr := WrapTracer(tracing.Tracer{})
	require.Zero(t, tr)
	ctx := context.Background()
	got, span := tr.Start(ctx, "disabled", nil)
	require.Equal(t, ctx, got)
	require.Zero(t, span)
}

func TestWrapTracerWithoutSpanLookup(t *testing.T) {
	tr := WrapTracer(tracing.NewTracer(func(ctx context.Context, _ string, _ *tracing.SpanOptions) (context.Context, tracing.Span) {
		return ctx, tracing.Span{}
	}, nil))
	ctx, _ := tr.Start(context.Background(), "operation", nil)
	require.True(t, tr.Enabled())
	require.Zero(t, tr.SpanFromContext(ctx))
}

func TestWrapTracerForwarding(t *testing.T) {
	type spanKey struct{}
	var attrs []tracing.Attribute
	var eventName string
	var eventAttrs []tracing.Attribute
	var status tracing.SpanStatus
	var description string
	var ended bool
	span := tracing.NewSpan(tracing.SpanImpl{
		End: func() { ended = true },
		SetAttributes: func(a ...tracing.Attribute) {
			attrs = append(attrs, a...)
		},
		AddEvent: func(name string, a ...tracing.Attribute) {
			eventName, eventAttrs = name, a
		},
		SetStatus: func(code tracing.SpanStatus, desc string) {
			status, description = code, desc
		},
	})
	tr := tracing.NewTracer(func(ctx context.Context, name string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
		require.Equal(t, "operation", name)
		require.Equal(t, tracing.SpanKindClient, options.Kind)
		attrs = options.Attributes
		return context.WithValue(ctx, spanKey{}, span), span
	}, &tracing.TracerOptions{
		SpanFromContext: func(ctx context.Context) tracing.Span {
			if span, ok := ctx.Value(spanKey{}).(tracing.Span); ok {
				return span
			}
			return tracing.Span{}
		},
	})
	tr.SetAttributes(tracing.Attribute{Key: "original", Value: true})
	tr = WrapTracer(tr)
	tr.SetAttributes(tracing.Attribute{Key: "wrapper", Value: true})
	require.Zero(t, tr.SpanFromContext(context.Background()))
	ctx, wrapped := tr.Start(context.Background(), "operation", &tracing.SpanOptions{
		Kind:       tracing.SpanKindClient,
		Attributes: []tracing.Attribute{{Key: "call", Value: true}},
	})
	require.NotNil(t, ctx.Value(spanKey{}))
	wrapped.SetAttributes(tracing.Attribute{Key: "span", Value: true})
	require.Equal(t, []tracing.Attribute{
		{Key: "call", Value: true},
		{Key: "original", Value: true},
		{Key: "wrapper", Value: true},
		{Key: "span", Value: true},
	}, attrs)
	wrapped.AddEvent("event", tracing.Attribute{Key: "event", Value: true})
	require.Equal(t, "event", eventName)
	require.Equal(t, []tracing.Attribute{{Key: "event", Value: true}}, eventAttrs)

	u, err := url.Parse("https://example.com?sig=secret&visible=value")
	require.NoError(t, err)
	RecordURL(ctx, *u, New(Options{AllowedQueryParams: []string{"visible"}}))
	// Recording takes a snapshot rather than keeping a mutable request URL.
	u.RawQuery = "sig=changed"
	tr.SpanFromContext(ctx).SetStatus(tracing.SpanStatusError, "failed: https://example.com?sig=secret&visible=value")
	require.Equal(t, tracing.SpanStatusError, status)
	require.Equal(t, "failed: https://example.com?sig=REDACTED&visible=value", description)
	wrapped.SetStatus(tracing.SpanStatusError, "unknown: https://other.com?sig=secret")
	require.Equal(t, "unknown: https://other.com?sig=secret", description)
	wrapped.End()
	require.True(t, ended)
}

func TestWrapTracerAttributeOverrides(t *testing.T) {
	var attrs []tracing.Attribute
	tr := tracing.NewTracer(func(ctx context.Context, _ string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
		attrs = options.Attributes
		return ctx, tracing.Span{}
	}, nil)
	tr.SetAttributes(tracing.Attribute{Key: "az.namespace", Value: "original"})
	original := tr
	tr = WrapTracer(tr)
	tr.SetAttributes(tracing.Attribute{Key: "az.namespace", Value: "override"})
	options := &tracing.SpanOptions{
		Attributes: []tracing.Attribute{{Key: "az.namespace", Value: "call"}},
	}
	tr.Start(context.Background(), "operation", options)
	require.Equal(t, []tracing.Attribute{
		{Key: "az.namespace", Value: "call"},
		{Key: "az.namespace", Value: "original"},
		{Key: "az.namespace", Value: "override"},
	}, attrs)
	original.Start(context.Background(), "original", nil)
	require.Equal(t, []tracing.Attribute{{Key: "az.namespace", Value: "original"}}, attrs)
}

func TestWrapTracerCredentialRequestURL(t *testing.T) {
	descriptions := map[string]string{}
	tr := WrapTracer(tracing.NewTracer(func(ctx context.Context, name string, _ *tracing.SpanOptions) (context.Context, tracing.Span) {
		return ctx, tracing.NewSpan(tracing.SpanImpl{
			SetStatus: func(_ tracing.SpanStatus, description string) { descriptions[name] = description },
		})
	}, nil))
	ctx, operation := tr.Start(context.Background(), "operation", nil)
	serviceSanitizer := New(Options{AllowedQueryParams: []string{"service"}})
	tokenSanitizer := New(Options{AllowedQueryParams: []string{"token"}})
	serviceURL, err := url.Parse("https://service.example.com?service=visible&token=service-secret")
	require.NoError(t, err)
	RecordURL(ctx, *serviceURL, serviceSanitizer)

	isolated := &shared.ContextWithDeniedValues{Context: ctx}
	tokenURL, err := url.Parse("https://identity.example.com?service=token-secret&token=visible")
	require.NoError(t, err)
	tokenCtx, token := tr.Start(isolated, "token", nil)
	tokenCtx, httpSpan := tr.Start(tokenCtx, "HTTP GET", nil)
	RecordURL(tokenCtx, *tokenURL, tokenSanitizer)

	description := "failed: " + tokenURL.String()
	token.SetStatus(tracing.SpanStatusError, description)
	httpSpan.SetStatus(tracing.SpanStatusError, description)
	operation.SetStatus(tracing.SpanStatusError, description)
	require.Equal(t, "failed: https://identity.example.com?service=REDACTED&token=visible", descriptions["token"])
	require.Equal(t, descriptions["token"], descriptions["HTTP GET"])
	require.Equal(t, descriptions["token"], descriptions["operation"])

	RecordURL(ctx, *serviceURL, serviceSanitizer)
	operation.SetStatus(tracing.SpanStatusError, "failed: "+serviceURL.String())
	require.Equal(t, "failed: https://service.example.com?service=visible&token=REDACTED", descriptions["operation"])
}

func TestWrapTracerLatestRequestURL(t *testing.T) {
	descriptions := map[string]string{}
	tr := WrapTracer(tracing.NewTracer(func(ctx context.Context, name string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
		return ctx, tracing.NewSpan(tracing.SpanImpl{
			SetStatus: func(_ tracing.SpanStatus, desc string) { descriptions[name] = desc },
		})
	}, nil))
	s := New(Options{})
	ctx, parent := tr.Start(context.Background(), "parent", nil)
	ctx, child := tr.Start(ctx, "child", nil)
	_, unrelated := tr.Start(context.Background(), "unrelated", nil)
	short, err := url.Parse("https://example.com?sig=secret")
	require.NoError(t, err)
	long, err := url.Parse("https://example.com?sig=secret&token=another")
	require.NoError(t, err)
	description := "failed: " + long.String()
	parent.SetStatus(tracing.SpanStatusError, description)
	require.Equal(t, description, descriptions["parent"])

	RecordURL(ctx, *short, s)
	RecordURL(ctx, *long, s)
	child.SetStatus(tracing.SpanStatusError, description)
	parent.SetStatus(tracing.SpanStatusError, description)
	unrelated.SetStatus(tracing.SpanStatusError, description)
	require.Equal(t, "failed: https://example.com?sig=REDACTED&token=REDACTED", descriptions["child"])
	require.Equal(t, descriptions["child"], descriptions["parent"])
	require.Equal(t, description, descriptions["unrelated"])

	RecordURL(ctx, *short, s)
	description = "failed: " + short.String()
	child.SetStatus(tracing.SpanStatusError, description)
	parent.SetStatus(tracing.SpanStatusError, description)
	require.Equal(t, "failed: https://example.com?sig=REDACTED", descriptions["child"])
	require.Equal(t, descriptions["child"], descriptions["parent"])

	previous, err := url.Parse("https://previous.example.com?sig=previous-secret")
	require.NoError(t, err)
	RecordURL(ctx, *previous, s)
	RecordURL(ctx, *short, s)
	description = "failed: " + previous.String()
	parent.SetStatus(tracing.SpanStatusError, description)
	require.Equal(t, description, descriptions["parent"])

	RecordURL(ctx, url.URL{Scheme: "https", Host: "example.com"}, s)
	description = "failed: " + short.String()
	parent.SetStatus(tracing.SpanStatusError, description)
	require.Equal(t, description, descriptions["parent"])
	RecordURL(context.Background(), *short, s)
}
