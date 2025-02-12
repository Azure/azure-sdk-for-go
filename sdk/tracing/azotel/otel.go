//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// azotel adapts OpenTelemetry tracing for consumption by the azcore/tracing package.
package azotel

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// TracingProviderOptions contains the optional values for NewTracingProvider.
type TracingProviderOptions struct {
	// for future expansion
}

// NewTracingProvider creates a new tracing.Provider that wraps the specified OpenTelemetry TracerProvider.
//   - tracerProvider - the TracerProvider to wrap
//   - opts - optional configuration. pass nil to accept the default values
func NewTracingProvider(tracerProvider trace.TracerProvider, opts *TracingProviderOptions) tracing.Provider {
	return tracing.NewProvider(func(namespace, version string) tracing.Tracer {
		tracer := tracerProvider.Tracer(namespace, trace.WithInstrumentationVersion(version), trace.WithSchemaURL(semconv.SchemaURL))

		return tracing.NewTracer(func(ctx context.Context, spanName string, options *tracing.SpanOptions) (context.Context, tracing.Span) {
			kind := tracing.SpanKindInternal
			var attrs []attribute.KeyValue
			var links []trace.Link
			if options != nil {
				kind = options.Kind
				attrs = convertAttributes(options.Attributes)
				links = convertLinks(options.Links)
			}
			ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(convertSpanKind(kind)), trace.WithAttributes(attrs...), trace.WithLinks(links...))
			return ctx, convertSpan(span)
		}, &tracing.TracerOptions{
			SpanFromContext: func(ctx context.Context) tracing.Span {
				return convertSpan(trace.SpanFromContext(ctx))
			},
			LinkFromContext: func(ctx context.Context, attrs ...tracing.Attribute) tracing.Link {
				link := trace.LinkFromContext(ctx, convertAttributes(attrs)...)
				return tracing.Link{
					SpanContext: convertOTelSpanContext(link.SpanContext),
					Attributes:  attrs,
				}
			},
		})
	}, &tracing.ProviderOptions{
		NewPropagatorFn: func() tracing.Propagator {
			return convertPropagator(propagation.TraceContext{})
		},
	})
}

func convertSpan(traceSpan trace.Span) tracing.Span {
	impl := tracing.SpanImpl{
		End: func() {
			traceSpan.End()
		},
		SetAttributes: func(attrs ...tracing.Attribute) {
			traceSpan.SetAttributes(convertAttributes(attrs)...)
		},
		AddEvent: func(name string, attrs ...tracing.Attribute) {
			traceSpan.AddEvent(name, trace.WithAttributes(convertAttributes(attrs)...))
		},
		AddLink: func(link tracing.Link) {
			traceSpan.AddLink(convertLink(link))
		},
		SpanContext: func() tracing.SpanContext {
			return convertOTelSpanContext(traceSpan.SpanContext())
		},
		SetStatus: func(code tracing.SpanStatus, desc string) {
			traceSpan.SetStatus(convertStatus(code), desc)
		},
	}
	return tracing.NewSpan(impl)
}

func convertAttributes(attrs []tracing.Attribute) []attribute.KeyValue {
	keyvals := []attribute.KeyValue{}
	for _, kv := range attrs {
		switch vv := kv.Value.(type) {
		case int:
			keyvals = append(keyvals, attribute.Int(kv.Key, vv))
		case int64:
			keyvals = append(keyvals, attribute.Int64(kv.Key, vv))
		case float64:
			keyvals = append(keyvals, attribute.Float64(kv.Key, vv))
		case bool:
			keyvals = append(keyvals, attribute.Bool(kv.Key, vv))
		case string:
			keyvals = append(keyvals, attribute.String(kv.Key, vv))
		default:
			keyvals = append(keyvals, attribute.String(kv.Key, fmt.Sprintf("%v", vv)))
		}
	}
	return keyvals
}

func convertLinks(links []tracing.Link) []trace.Link {
	var otelLinks []trace.Link
	for _, link := range links {
		otelLinks = append(otelLinks, convertLink(link))
	}
	return otelLinks
}

func convertLink(link tracing.Link) trace.Link {
	return trace.Link{
		SpanContext: convertSpanContext(link.SpanContext),
		Attributes:  convertAttributes(link.Attributes),
	}
}

func convertSpanContext(spanContext tracing.SpanContext) trace.SpanContext {
	oTelTraceState, _ := trace.ParseTraceState(string(spanContext.TraceState()))
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID(spanContext.TraceID()),
		SpanID:     trace.SpanID(spanContext.SpanID()),
		TraceFlags: trace.TraceFlags(spanContext.TraceFlags()),
		TraceState: oTelTraceState,
		Remote:     spanContext.IsRemote(),
	})
}

func convertOTelSpanContext(spanContext trace.SpanContext) tracing.SpanContext {
	return tracing.NewSpanContext(tracing.SpanContextConfig{
		TraceID:    tracing.TraceID(spanContext.TraceID()),
		SpanID:     tracing.SpanID(spanContext.SpanID()),
		TraceFlags: tracing.TraceFlags(spanContext.TraceFlags()),
		TraceState: tracing.TraceState(spanContext.TraceState().String()),
		Remote:     spanContext.IsRemote(),
	})
}

func convertSpanKind(sk tracing.SpanKind) trace.SpanKind {
	switch sk {
	case tracing.SpanKindServer:
		return trace.SpanKindServer
	case tracing.SpanKindClient:
		return trace.SpanKindClient
	case tracing.SpanKindProducer:
		return trace.SpanKindProducer
	case tracing.SpanKindConsumer:
		return trace.SpanKindConsumer
	default:
		return trace.SpanKindInternal
	}
}

func convertStatus(ss tracing.SpanStatus) codes.Code {
	switch ss {
	case tracing.SpanStatusError:
		return codes.Error
	case tracing.SpanStatusOK:
		return codes.Ok
	default:
		return codes.Unset
	}
}

func convertPropagator(pr propagation.TextMapPropagator) tracing.Propagator {
	return tracing.NewPropagator(tracing.PropagatorImpl{
		Inject: func(ctx context.Context, carrier tracing.Carrier) {
			pr.Inject(ctx, carrier)
		},
		Extract: func(ctx context.Context, carrier tracing.Carrier) context.Context {
			return pr.Extract(ctx, carrier)
		},
		Fields: func() []string {
			return pr.Fields()
		},
	})
}
