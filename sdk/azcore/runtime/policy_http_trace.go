//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

// httpTracePolicy is a policy that creates a trace for the HTTP request and its response
func httpTracePolicy(req *policy.Request) (resp *http.Response, err error) {
	rawTracer := req.Raw().Context().Value(shared.CtxWithTracingTracer{})
	if tracer, ok := rawTracer.(tracing.Tracer); ok {
		attributes := []tracing.Attribute{
			{Key: "http.method", Value: req.Raw().Method},
			{Key: "http.url", Value: req.Raw().URL.String()},
		}

		if ua := req.Raw().Header.Get(shared.HeaderUserAgent); ua != "" {
			attributes = append(attributes, tracing.Attribute{Key: "http.user_agent", Value: ua})
		}
		if reqID := req.Raw().Header.Get(shared.HeaderXMSClientRequestID); reqID != "" {
			attributes = append(attributes, tracing.Attribute{Key: "az.client_request_id", Value: reqID})
		}

		ctx := req.Raw().Context()
		ctx, span := tracer.Start(ctx, "HTTP "+req.Raw().Method, &tracing.SpanOptions{
			Kind:       tracing.SpanKindClient,
			Attributes: attributes,
		})

		defer func() {
			if resp != nil {
				span.SetAttributes(tracing.Attribute{Key: "http.status_code", Value: resp.StatusCode})
				if reqID := resp.Header.Get(shared.HeaderXMSRequestID); reqID != "" {
					span.SetAttributes(tracing.Attribute{Key: "az.service_request_id", Value: reqID})
				}
			} else if err != nil {
				span.SetStatus(tracing.SpanStatusError, err.Error())
			}
			span.End()
		}()

		req = req.WithContext(ctx)
		err = tracer.Inject(req.Raw().Context(), req.Raw())
		if err != nil {
			return
		}
	}
	resp, err = req.Next()
	return
}

// StartSpanOptions contains the optional values for StartSpan.
type StartSpanOptions struct {
	// for future expansion
}

// StartSpan starts a new tracing span.
// You must call the returned func to terminate the span. Pass the applicable error
// if the span will exit with an error condition.
//   - ctx is the parent context of the newly created context
//   - name is the name of the span. this is typically the fully qualified name of an API ("Client.Method")
//   - tracer is the client's Tracer for creating spans
//   - options contains optional values. pass nil to accept any default values
func StartSpan(ctx context.Context, name string, tracer tracing.Tracer, options *StartSpanOptions) (context.Context, func(error)) {
	if !tracer.Enabled() {
		return ctx, func(err error) {}
	}
	ctx, span := tracer.Start(ctx, name, &tracing.SpanOptions{
		Kind: tracing.SpanKindInternal,
	})
	ctx = context.WithValue(ctx, shared.CtxWithTracingTracer{}, tracer)
	return ctx, func(err error) {
		if err != nil {
			span.SetStatus(tracing.SpanStatusError, err.Error())
		}
		span.End()
	}
}
