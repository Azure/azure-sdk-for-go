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

	ctx, endSpan := azruntime.StartSpan(context.Background(), "test_span", client.Tracer(), nil)

	req, err := azruntime.NewRequest(ctx, http.MethodGet, "https://www.microsoft.com/")
	require.NoError(t, err)

	startedSpan := client.Tracer().SpanFromContext(req.Raw().Context())
	startedSpan.AddEvent("post_event")

	_, err = client.Pipeline().Do(req)
	require.NoError(t, err)

	endSpan(nil)

	// shut down the tracing provider to flush all spans
	require.NoError(t, otelTP.Shutdown(context.Background()))

	require.Len(t, exporter.spans, 2)
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
	eventName    string
	eventOptions []trace.EventOption
	endCalled    bool
	statusCode   codes.Code
	statusDesc   string
	link         trace.Link
}

func (ts *testSpan) End(options ...trace.SpanEndOption) {
	ts.endCalled = true
}

func (ts *testSpan) AddEvent(name string, options ...trace.EventOption) {
	ts.eventName = name
	ts.eventOptions = options
}

func (ts *testSpan) AddLink(link trace.Link) {
	ts.link = link
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

func (ts *testSpan) TracerProvider() trace.TracerProvider {
	ts.t.Fatal("TracerProvider not required")
	return nil
}
