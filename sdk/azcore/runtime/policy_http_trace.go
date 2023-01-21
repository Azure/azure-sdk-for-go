//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

// httpTracePolicy is a policy that creates a trace for the HTTP request and its response
func httpTracePolicy(req *policy.Request) (resp *http.Response, err error) {
	var tracer tracing.Tracer
	if req.OperationValue(&tracer) {
		attributes := []tracing.Attribute{
			{Key: "http.method", Value: req.Raw().Method},
			{Key: "http.url", Value: req.Raw().URL.String()},
		}

		if ua := req.Raw().Header.Get(shared.HeaderUserAgent); ua != "" {
			attributes = append(attributes, tracing.Attribute{Key: "http.user_agent", Value: ua})
		}
		if reqID := req.Raw().Header.Get(shared.HeaderXMSClientRequestID); reqID != "" {
			attributes = append(attributes, tracing.Attribute{Key: "requestId", Value: reqID})
		}

		ctx := req.Raw().Context()
		ctx, span := tracer.Start(ctx, "azure-sdk.http", &tracing.SpanOptions{
			Kind:       tracing.SpanKindClient,
			Attributes: attributes,
		})

		defer func() {
			if resp != nil {
				span.SetAttributes(tracing.Attribute{Key: "http.status_code", Value: resp.StatusCode})
				if reqID := resp.Header.Get(shared.HeaderXMSRequestID); reqID != "" {
					span.SetAttributes(tracing.Attribute{Key: "serviceRequestId", Value: reqID})
				}
			}
			span.End()
		}()

		req = req.WithContext(ctx)
	}
	resp, err = req.Next()
	return
}
