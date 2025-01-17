// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

type Attribute = tracing.Attribute
type Tracer = tracing.Tracer
type Provider = tracing.Provider

type TracerOptions struct {
	Tracer     Tracer
	SpanName   SpanName
	Attributes []Attribute
}

// StartSpan creates a span with the specified name and attributes.
// If no span name is provided, no span is created.
func StartSpan(ctx context.Context, options *TracerOptions) (context.Context, func(error)) {
	if options == nil || options.SpanName == "" {
		return ctx, func(error) {}
	}
	spanKind := SpanKindInternal
	spanCaller := strings.Split(string(options.SpanName), ".")[0]
	if spanCaller == "Sender" {
		spanKind = SpanKindProducer
	} else if spanCaller == "Receiver" || spanCaller == "SessionReceiver" {
		spanKind = SpanKindConsumer
	}

	return runtime.StartSpan(ctx, string(options.SpanName), options.Tracer,
		&runtime.StartSpanOptions{
			Kind:       spanKind,
			Attributes: options.Attributes,
		})
}

func GetEntityPathAttributes(entityPath string) []tracing.Attribute {
	var attrs []tracing.Attribute
	queueOrTopic, subscription := splitEntityPath(entityPath)
	if queueOrTopic != "" {
		attrs = append(attrs, tracing.Attribute{Key: DestinationName, Value: queueOrTopic})
	}
	if subscription != "" {
		attrs = append(attrs, tracing.Attribute{Key: SubscriptionName, Value: subscription})
	}
	return attrs
}

func splitEntityPath(entityPath string) (string, string) {
	queueOrTopic := ""
	subscription := ""

	path := strings.Split(entityPath, "/")
	if len(path) >= 1 {
		queueOrTopic = path[0]
	}
	if len(path) >= 3 && path[1] == "Subscriptions" {
		subscription = path[2]
	}
	return queueOrTopic, subscription
}
