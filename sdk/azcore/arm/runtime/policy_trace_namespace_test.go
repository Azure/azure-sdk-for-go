//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestHTTPTraceNamespacePolicy(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()

	pl := exported.NewPipeline(srv, exported.PolicyFunc(httpTraceNamespacePolicy))

	// no tracer
	req, err := exported.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	srv.AppendResponse()
	_, err = pl.Do(req)
	require.NoError(t, err)

	// wrong tracer type
	req, err = exported.NewRequest(context.WithValue(context.Background(), shared.CtxWithTracingTracer{}, 0), http.MethodGet, srv.URL())
	require.NoError(t, err)
	srv.AppendResponse()
	_, err = pl.Do(req)
	require.NoError(t, err)

	// no SpanFromContext impl
	tr := tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
		return ctx, tracing.Span{}
	}, nil)
	req, err = exported.NewRequest(context.WithValue(context.Background(), shared.CtxWithTracingTracer{}, tr), http.MethodGet, srv.URL())
	require.NoError(t, err)
	srv.AppendResponse()
	_, err = pl.Do(req)
	require.NoError(t, err)

	// failed to parse resource ID, shouldn't call SetAttributes
	var attrString string
	tr = tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
		return ctx, tracing.Span{}
	}, &tracing.TracerOptions{
		SpanFromContext: func(ctx context.Context) tracing.Span {
			spanImpl := tracing.SpanImpl{
				SetAttributes: func(a ...tracing.Attribute) {
					require.Len(t, a, 1)
					v, ok := a[0].Value.(string)
					require.True(t, ok)
					attrString = a[0].Key + ":" + v
				},
			}
			return tracing.NewSpan(spanImpl)
		},
	})
	req, err = exported.NewRequest(context.WithValue(context.Background(), shared.CtxWithTracingTracer{}, tr), http.MethodGet, srv.URL())
	require.NoError(t, err)
	srv.AppendResponse()
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.Empty(t, attrString)

	// success
	tr = tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
		return ctx, tracing.Span{}
	}, &tracing.TracerOptions{
		SpanFromContext: func(ctx context.Context) tracing.Span {
			spanImpl := tracing.SpanImpl{
				SetAttributes: func(a ...tracing.Attribute) {
					require.Len(t, a, 1)
					v, ok := a[0].Value.(string)
					require.True(t, ok)
					attrString = a[0].Key + ":" + v
				},
			}
			return tracing.NewSpan(spanImpl)
		},
	})
	req, err = exported.NewRequest(context.WithValue(context.Background(), shared.CtxWithTracingTracer{}, tr), http.MethodGet, srv.URL()+requestEndpoint)
	require.NoError(t, err)
	srv.AppendResponse()
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, "az.namespace:Microsoft.Storage", attrString)
}
