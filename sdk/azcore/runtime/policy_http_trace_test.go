//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestHTTPTracePolicy(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()

	pl := exported.NewPipeline(srv, newHTTPTracePolicy([]string{"visibleqp"}))

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

	var fullSpanName string
	var spanKind tracing.SpanKind
	var spanAttrs []tracing.Attribute
	var spanStatus tracing.SpanStatus
	var spanStatusStr string
	tr := tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
		fullSpanName = spanName
		require.NotNil(t, options)
		spanKind = options.Kind
		spanAttrs = options.Attributes
		spanImpl := tracing.SpanImpl{
			SetAttributes: func(a ...tracing.Attribute) { spanAttrs = append(spanAttrs, a...) },
			SetStatus: func(ss tracing.SpanStatus, s string) {
				spanStatus = ss
				spanStatusStr = s
			},
		}
		return ctx, tracing.NewSpan(spanImpl)
	}, nil)

	// HTTP ok
	req, err = exported.NewRequest(context.WithValue(context.Background(), shared.CtxWithTracingTracer{}, tr), http.MethodGet, srv.URL()+"?foo=redactme&visibleqp=bar")
	require.NoError(t, err)
	req.Raw().Header.Add(shared.HeaderUserAgent, "my-user-agent")
	req.Raw().Header.Add(shared.HeaderXMSClientRequestID, "my-client-request")
	srv.AppendResponse(mock.WithHeader(shared.HeaderXMSRequestID, "request-id"))
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, tracing.SpanStatusUnset, spanStatus)
	require.EqualValues(t, "HTTP GET", fullSpanName)
	require.EqualValues(t, tracing.SpanKindClient, spanKind)
	require.Len(t, spanAttrs, 7)
	require.Contains(t, spanAttrs, tracing.Attribute{Key: attrHTTPMethod, Value: http.MethodGet})
	require.Contains(t, spanAttrs, tracing.Attribute{Key: attrHTTPURL, Value: srv.URL() + "?foo=REDACTED&visibleqp=bar"})
	require.Contains(t, spanAttrs, tracing.Attribute{Key: attrNetPeerName, Value: srv.URL()[7:]}) // strip off the http://
	require.Contains(t, spanAttrs, tracing.Attribute{Key: attrHTTPUserAgent, Value: "my-user-agent"})
	require.Contains(t, spanAttrs, tracing.Attribute{Key: attrAZClientReqID, Value: "my-client-request"})
	require.Contains(t, spanAttrs, tracing.Attribute{Key: attrHTTPStatusCode, Value: http.StatusOK})
	require.Contains(t, spanAttrs, tracing.Attribute{Key: attrAZServiceReqID, Value: "request-id"})

	// HTTP bad request
	req, err = exported.NewRequest(context.WithValue(context.Background(), shared.CtxWithTracingTracer{}, tr), http.MethodGet, srv.URL())
	require.NoError(t, err)
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	_, err = pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, tracing.SpanStatusError, spanStatus)
	require.EqualValues(t, "400 Bad Request", spanStatusStr)
	require.Contains(t, spanAttrs, tracing.Attribute{Key: attrHTTPStatusCode, Value: http.StatusBadRequest})

	// HTTP error
	req, err = exported.NewRequest(context.WithValue(context.Background(), shared.CtxWithTracingTracer{}, tr), http.MethodGet, srv.URL())
	require.NoError(t, err)
	srv.AppendError(net.ErrClosed)
	_, err = pl.Do(req)
	require.Error(t, err)
	require.ErrorIs(t, err, net.ErrClosed)
	require.EqualValues(t, tracing.SpanStatusError, spanStatus)
	require.EqualValues(t, "use of closed network connection", spanStatusStr)

	const urlErrText = "the endpoint is invalid"
	req, err = exported.NewRequest(context.WithValue(context.Background(), shared.CtxWithTracingTracer{}, tr), http.MethodGet, srv.URL())
	require.NoError(t, err)
	srv.AppendError(&url.Error{
		Op:  http.MethodGet,
		URL: srv.URL(),
		Err: errors.New(urlErrText),
	})
	_, err = pl.Do(req)
	require.Error(t, err)
	var urlErr *url.Error
	require.False(t, errors.As(err, &urlErr))
	require.EqualValues(t, tracing.SpanStatusError, spanStatus)
	require.EqualValues(t, urlErrText, spanStatusStr)
}

func TestStartSpan(t *testing.T) {
	// tracing disabled
	ctxIn := context.Background()
	ctx, end := StartSpan(ctxIn, "TestStartSpan", tracing.Tracer{}, nil)
	end(nil)
	require.Equal(t, ctxIn, ctx)

	// span no error
	var startCalled bool
	var endCalled bool
	tr := tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
		startCalled = true
		require.EqualValues(t, "TestStartSpan", spanName)
		require.NotNil(t, options)
		require.EqualValues(t, tracing.SpanKindInternal, options.Kind)
		spanImpl := tracing.SpanImpl{
			End: func() { endCalled = true },
		}
		return ctx, tracing.NewSpan(spanImpl)
	}, nil)
	ctx, end = StartSpan(context.Background(), "TestStartSpan", tr, nil)
	end(nil)
	ctxTr := ctx.Value(shared.CtxWithTracingTracer{})
	require.NotNil(t, ctxTr)
	_, ok := ctxTr.(tracing.Tracer)
	require.True(t, ok)
	require.True(t, startCalled)
	require.True(t, endCalled)

	// with error
	var spanStatus tracing.SpanStatus
	var errStr string
	tr = tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
		spanImpl := tracing.SpanImpl{
			End: func() { endCalled = true },
			SetStatus: func(ss tracing.SpanStatus, s string) {
				spanStatus = ss
				errStr = s
			},
		}
		return ctx, tracing.NewSpan(spanImpl)
	}, nil)
	_, end = StartSpan(context.Background(), "TestStartSpan", tr, nil)
	u, err := url.Parse("https://contoso.com")
	require.NoError(t, err)
	resp := &http.Response{
		Status:     "the operation failed",
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader(`{ "error": { "code": "ErrorItFailed", "message": "it's not working" } }`)),
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    u,
		},
	}
	end(exported.NewResponseError(resp))
	require.EqualValues(t, tracing.SpanStatusError, spanStatus)
	require.Contains(t, errStr, "*azcore.ResponseError")
	require.Contains(t, errStr, "ERROR CODE: ErrorItFailed")
}

func TestStartSpansDontNest(t *testing.T) {
	srv, close := mock.NewServer()
	srv.SetResponse() // always return http.StatusOK
	defer close()

	pl := exported.NewPipeline(srv, newHTTPTracePolicy(nil))

	apiSpanCount := 0
	httpSpanCount := 0
	endCalled := 0
	tr := tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
		if spanName == "HTTP GET" {
			httpSpanCount++
		} else if spanName == "FooMethod" {
			apiSpanCount++
		} else {
			t.Fatalf("unexpected span name %s", spanName)
		}
		spanImpl := tracing.SpanImpl{
			End: func() { endCalled++ },
		}
		return ctx, tracing.NewSpan(spanImpl)
	}, nil)

	barMethod := func(ctx context.Context) {
		ourCtx, endSpan := StartSpan(ctx, "BarMethod", tr, nil)
		defer endSpan(nil)
		req, err := exported.NewRequest(ourCtx, http.MethodGet, srv.URL()+"/bar")
		require.NoError(t, err)
		_, err = pl.Do(req)
		require.NoError(t, err)
	}

	fooMethod := func(ctx context.Context) {
		ctx, endSpan := StartSpan(ctx, "FooMethod", tr, nil)
		defer endSpan(nil)
		barMethod(ctx)
		req, err := exported.NewRequest(ctx, http.MethodGet, srv.URL()+"/foo")
		require.NoError(t, err)
		_, err = pl.Do(req)
		require.NoError(t, err)
	}

	fooMethod(context.Background())

	// there should be a total of three spans. one for FooMethod, and two HTTP spans
	require.EqualValues(t, 1, apiSpanCount)
	require.EqualValues(t, 2, httpSpanCount)
	require.EqualValues(t, 3, endCalled)
}
