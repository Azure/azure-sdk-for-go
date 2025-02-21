//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azotel

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

func TestNewTracingProvider(t *testing.T) {
	exporter := &testExporter{}

	otelTP := tracesdk.NewTracerProvider(tracesdk.WithBatcher(exporter))

	client, err := azcore.NewClient("azotel", internal.Version, azruntime.PipelineOptions{
		Tracing: azruntime.TracingOptions{
			Namespace: "TestNewTracingProvider",
		},
	}, &azcore.ClientOptions{
		TracingProvider: NewTracingProvider(otelTP, nil),
	})
	require.NoError(t, err)

	// returns a no-op span as there is no span yet
	emptySpan := client.Tracer().SpanFromContext(context.Background())
	emptySpan.AddEvent("noop_event")

	// returns an empty link as there is no span yet
	emptyLink := client.Tracer().LinkFromContext(context.Background())
	require.Empty(t, emptyLink)

	ctx, endSpan := azruntime.StartSpan(context.Background(), "test_span", client.Tracer(), nil)

	req, err := azruntime.NewRequest(ctx, http.MethodGet, "https://www.microsoft.com/")
	require.NoError(t, err)

	startedSpan := client.Tracer().SpanFromContext(req.Raw().Context())
	startedSpan.AddEvent("post_event")

	startedSpanLink := client.Tracer().LinkFromContext(req.Raw().Context())
	require.NotEmpty(t, startedSpanLink)

	_, err = client.Pipeline().Do(req)
	require.NoError(t, err)

	endSpan(nil)

	// shut down the tracing provider to flush all spans
	require.NoError(t, otelTP.Shutdown(context.Background()))

	require.Len(t, exporter.spans, 2)
}

func TestPropagator(t *testing.T) {
	provider := NewTracingProvider(tracesdk.NewTracerProvider(), nil)
	tracer := provider.NewTracer("test", "1.0")
	propagator := provider.NewPropagator()
	require.EqualValues(t, 2, len(propagator.Fields()))

	mapCarrier := propagation.MapCarrier{}
	carrier := tracing.NewCarrier(tracing.CarrierImpl{
		Get:  mapCarrier.Get,
		Set:  mapCarrier.Set,
		Keys: mapCarrier.Keys,
	})

	ctx, endSpan := azruntime.StartSpan(context.Background(), "test_span", tracer, nil)
	spanContext := tracer.SpanFromContext(ctx).SpanContext()
	require.NotNil(t, spanContext)
	require.False(t, spanContext.IsRemote())
	propagator.Inject(ctx, carrier)
	endSpan(nil)

	extractedCtx := propagator.Extract(context.Background(), carrier)
	extractedSpanContext := tracer.SpanFromContext(extractedCtx).SpanContext()
	require.NotNil(t, extractedSpanContext)
	require.True(t, extractedSpanContext.IsRemote())
	require.EqualValues(t, spanContext.TraceID(), extractedSpanContext.TraceID())
	require.EqualValues(t, spanContext.SpanID(), extractedSpanContext.SpanID())
}

func TestConvertSpan(t *testing.T) {
	ts := testSpan{t: t}
	span := convertSpan(&ts)

	const eventName = "event"
	eventAttr := tracing.Attribute{
		Key:   "key",
		Value: "value",
	}
	span.AddEvent(eventName, eventAttr)
	assert.EqualValues(t, eventName, ts.eventName)
	require.Len(t, ts.eventOptions, 1)

	span.End()
	assert.True(t, ts.endCalled)

	attr := tracing.Attribute{
		Key:   "key",
		Value: "value",
	}
	span.SetAttributes(attr)
	require.Len(t, ts.attributes, 1)

	span.AddLink(tracing.Link{})
	require.Len(t, ts.links, 1)

	const statusDesc = "everything is ok"
	span.SetStatus(tracing.SpanStatusOK, statusDesc)
	assert.EqualValues(t, tracing.SpanStatusOK, ts.statusCode)
	assert.EqualValues(t, statusDesc, ts.statusDesc)
}

func TestConvertAttributes(t *testing.T) {
	tests := []struct {
		label    string
		key      string
		value    any
		validate func(t *testing.T, kv attribute.KeyValue)
	}{
		{
			label: "int",
			key:   "int_key",
			value: 123,
			validate: func(t *testing.T, kv attribute.KeyValue) {
				require.EqualValues(t, "INT64", kv.Value.Type().String())
				assert.EqualValues(t, 123, kv.Value.AsInt64())
			},
		},
		{
			label: "int64",
			key:   "int64_key",
			value: int64(123),
			validate: func(t *testing.T, kv attribute.KeyValue) {
				require.EqualValues(t, "INT64", kv.Value.Type().String())
				assert.EqualValues(t, 123, kv.Value.AsInt64())
			},
		},
		{
			label: "float64",
			key:   "float64_key",
			value: 3.14159,
			validate: func(t *testing.T, kv attribute.KeyValue) {
				require.EqualValues(t, "FLOAT64", kv.Value.Type().String())
				assert.EqualValues(t, 3.14159, kv.Value.AsFloat64())
			},
		},
		{
			label: "bool",
			key:   "bool_key",
			value: true,
			validate: func(t *testing.T, kv attribute.KeyValue) {
				require.EqualValues(t, "BOOL", kv.Value.Type().String())
				assert.True(t, kv.Value.AsBool())
			},
		},
		{
			label: "string",
			key:   "string_key",
			value: "hello",
			validate: func(t *testing.T, kv attribute.KeyValue) {
				require.EqualValues(t, "STRING", kv.Value.Type().String())
				assert.EqualValues(t, "hello", kv.Value.AsString())
			},
		},
		{
			label: "float32",
			key:   "float32_key",
			value: float32(3.14159),
			validate: func(t *testing.T, kv attribute.KeyValue) {
				require.EqualValues(t, "STRING", kv.Value.Type().String())
				assert.EqualValues(t, "3.14159", kv.Value.AsString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			keyval := convertAttributes([]tracing.Attribute{
				{
					Key:   tt.key,
					Value: tt.value,
				},
			})

			require.Len(t, keyval, 1)
			assert.EqualValues(t, tt.key, keyval[0].Key)
			tt.validate(t, keyval[0])
		})
	}
}

func TestConvertLinks(t *testing.T) {
	attr := tracing.Attribute{
		Key:   "key",
		Value: "value",
	}
	spanContext := tracing.NewSpanContext(tracing.SpanContextConfig{
		TraceID:    tracing.TraceID{1, 2, 3, 4, 5, 6, 7, 8},
		SpanID:     tracing.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
		TraceFlags: tracing.TraceFlags(0x1),
		TraceState: "key1=val1,key2=val2",
		Remote:     true,
	})

	links := convertLinks([]tracing.Link{
		{
			SpanContext: spanContext,
		},
		{
			Attributes: []tracing.Attribute{attr},
		},
		{
			SpanContext: spanContext,
			Attributes:  []tracing.Attribute{attr},
		},
	})
	require.Len(t, links, 3)
	require.NotNil(t, links[0].SpanContext)
	require.True(t, links[0].SpanContext.IsRemote())
	require.Len(t, links[0].Attributes, 0)

	require.NotNil(t, links[1].SpanContext)
	require.False(t, links[1].SpanContext.IsRemote())
	require.Len(t, links[1].Attributes, 1)

	require.NotNil(t, links[2].SpanContext)
	require.True(t, links[2].SpanContext.IsRemote())
	require.Len(t, links[2].Attributes, 1)
}

func TestConvertSpanContext(t *testing.T) {
	traceID := tracing.TraceID{1, 2, 3, 4, 5, 6, 7, 8}
	spanID := tracing.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	traceFlags := tracing.TraceFlags(0x1)
	spanContext := tracing.NewSpanContext(tracing.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: traceFlags,
		TraceState: "key1=val1,key2=val2",
		Remote:     true,
	})

	otelSpanContext := convertSpanContext(spanContext)
	assert.EqualValues(t, traceID, otelSpanContext.TraceID())
	assert.EqualValues(t, spanID, otelSpanContext.SpanID())
	assert.EqualValues(t, traceFlags, otelSpanContext.TraceFlags())
	assert.EqualValues(t, "key1=val1,key2=val2", otelSpanContext.TraceState().String())
	assert.True(t, otelSpanContext.IsRemote())
}

func TestConvertSpanKind(t *testing.T) {
	assert.EqualValues(t, trace.SpanKindClient, convertSpanKind(tracing.SpanKindClient))
	assert.EqualValues(t, trace.SpanKindConsumer, convertSpanKind(tracing.SpanKindConsumer))
	assert.EqualValues(t, trace.SpanKindInternal, convertSpanKind(tracing.SpanKindInternal))
	assert.EqualValues(t, trace.SpanKindProducer, convertSpanKind(tracing.SpanKindProducer))
	assert.EqualValues(t, trace.SpanKindServer, convertSpanKind(tracing.SpanKindServer))
	assert.EqualValues(t, trace.SpanKindInternal, convertSpanKind(tracing.SpanKind(12345)))
}

func TestConvertStatus(t *testing.T) {
	assert.EqualValues(t, codes.Ok, convertStatus(tracing.SpanStatusOK))
	assert.EqualValues(t, codes.Error, convertStatus(tracing.SpanStatusError))
	assert.EqualValues(t, codes.Unset, convertStatus(tracing.SpanStatusUnset))
	assert.EqualValues(t, codes.Unset, convertStatus(tracing.SpanStatus(12345)))
}

func TestConvertPropagator(t *testing.T) {
	carrier := tracing.NewCarrier(tracing.CarrierImpl{
		Get:  func(key string) string { return "" },
		Set:  func(key, value string) {},
		Keys: func() []string { return nil },
	})
	propagator := &testPropagator{}
	otelPropagator := convertPropagator(propagator)
	require.NotNil(t, otelPropagator)
	otelPropagator.Inject(context.Background(), carrier)
	otelPropagator.Extract(context.Background(), carrier)
	require.True(t, propagator.injectCalled)
	require.True(t, propagator.extractCalled)
	require.Len(t, propagator.Fields(), 1)
}

type testExporter struct {
	spans []string
}

// ExportSpans implements the tracesdk.SpanExporter interface for the ConsoleExporter type.
func (c *testExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	for _, span := range spans {
		c.spans = append(c.spans, span.Name())
	}
	return nil
}

// Shutdown implements the tracesdk.SpanExporter interface for the ConsoleExporter type.
func (c *testExporter) Shutdown(ctx context.Context) error {
	return nil
}

type testSpan struct {
	embedded.Span

	t            *testing.T
	attributes   []attribute.KeyValue
	links        []trace.Link
	eventName    string
	eventOptions []trace.EventOption
	endCalled    bool
	statusCode   codes.Code
	statusDesc   string
}

func (ts *testSpan) End(options ...trace.SpanEndOption) {
	ts.endCalled = true
}

func (ts *testSpan) AddEvent(name string, options ...trace.EventOption) {
	ts.eventName = name
	ts.eventOptions = options
}

func (ts *testSpan) IsRecording() bool {
	ts.t.Fatal("IsRecording not required")
	return false
}

func (ts *testSpan) RecordError(err error, options ...trace.EventOption) {
	ts.t.Fatal("RecordError not required")
}

func (ts *testSpan) SpanContext() trace.SpanContext {
	ts.t.Fatal("SpanContext not required")
	return trace.SpanContext{}
}

func (ts *testSpan) SetStatus(code codes.Code, description string) {
	ts.statusCode = code
	ts.statusDesc = description
}

func (ts *testSpan) SetName(name string) {
	ts.t.Fatal("SetName not required")
}

func (ts *testSpan) SetAttributes(kv ...attribute.KeyValue) {
	ts.attributes = kv
}

func (ts *testSpan) AddLink(link trace.Link) {
	ts.links = append(ts.links, link)
}

func (ts *testSpan) TracerProvider() trace.TracerProvider {
	ts.t.Fatal("TracerProvider not required")
	return nil
}

type testPropagator struct {
	injectCalled  bool
	extractCalled bool
}

func (tp *testPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	tp.injectCalled = true
}

func (tp *testPropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	tp.extractCalled = true
	return ctx
}

func (tp *testPropagator) Fields() []string {
	return []string{"testfield"}
}
