// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
)

const instrumentationScope = "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"

// Diagnostics holds the application-supplied telemetry instruments for one client.
type Diagnostics struct {
	tracer     tracing.Tracer
	logger     otellog.Logger
	operations metric.Int64Counter
	events     metric.Int64Counter
	duration   metric.Float64Histogram
	lifecycle  metric.Int64Counter
	host       string
	hub        string
}

func NewDiagnostics(host, hub string, tp tracing.Provider, mp metric.MeterProvider, lp otellog.LoggerProvider) *Diagnostics {
	tracer := tp.NewTracer(instrumentationScope, Version)
	if !tracer.Enabled() && mp == nil && lp == nil {
		return nil
	}
	d := &Diagnostics{host: host, hub: hub, tracer: tracer}
	d.tracer.SetAttributes(tracing.Attribute{Key: "server.address", Value: host}, tracing.Attribute{Key: "messaging.system", Value: "eventhubs"})
	if lp != nil {
		d.logger = lp.Logger(instrumentationScope, otellog.WithInstrumentationVersion(Version))
	}
	if mp != nil {
		meter := mp.Meter(instrumentationScope, metric.WithInstrumentationVersion(Version))
		// Instrument registration can return a usable instrument together with an error.
		// A missing instrument only disables that measurement.
		d.operations, _ = meter.Int64Counter("az.eventhubs.operations", metric.WithUnit("{operation}"))
		d.events, _ = meter.Int64Counter("az.eventhubs.events", metric.WithUnit("{event}"))
		d.duration, _ = meter.Float64Histogram("az.eventhubs.operation.duration", metric.WithUnit("s"),
			metric.WithExplicitBucketBoundaries(.001, .005, .01, .05, .1, .5, 1, 5, 10))
		d.lifecycle, _ = meter.Int64Counter("az.eventhubs.amqp.lifecycle", metric.WithUnit("{transition}"))
	}
	return d
}

func outcome(err error) string {
	switch {
	case err == nil:
		return "success"
	case IsCancelError(err):
		return "cancelled"
	default:
		return "error"
	}
}

func endSpan(span tracing.Span, err error) {
	if err != nil {
		// Broker error descriptions and custom transport errors can contain customer data.
		span.SetStatus(tracing.SpanStatusError, outcome(err))
	}
	span.End()
}

// StartOperation covers one public send or receive, including preparation and retries.
func (d *Diagnostics) StartOperation(ctx context.Context, operation string, kind tracing.SpanKind) (context.Context, func(int, error)) {
	if d == nil {
		return ctx, func(int, error) {}
	}
	start := time.Now()
	ctx, span := d.tracer.Start(ctx, "EventHubs."+operation, &tracing.SpanOptions{Kind: kind, Attributes: []tracing.Attribute{
		{Key: "messaging.destination.name", Value: d.hub}, {Key: "messaging.operation.name", Value: operation},
	}})
	return ctx, func(count int, err error) {
		span.SetAttributes(tracing.Attribute{Key: "messaging.batch.message_count", Value: count})
		defer endSpan(span, err)
		if d.operations == nil && d.events == nil && d.duration == nil {
			return
		}
		attrs := metric.WithAttributes(attribute.String("messaging.system", "eventhubs"), attribute.String("server.address", d.host),
			attribute.String("messaging.destination.name", d.hub), attribute.String("operation", operation), attribute.String("outcome", outcome(err)))
		if d.duration != nil {
			d.duration.Record(ctx, time.Since(start).Seconds(), attrs)
		}
		if d.operations != nil {
			d.operations.Add(ctx, 1, attrs)
		}
		if d.events != nil && err == nil && count > 0 {
			d.events.Add(ctx, int64(count), attrs)
		}
	}
}

// StartAttempt covers the transport work within one retry attempt.
func (d *Diagnostics) StartAttempt(ctx context.Context, name string) (context.Context, func(error)) {
	if d == nil || !d.tracer.Enabled() {
		return ctx, func(error) {}
	}
	ctx, span := d.tracer.Start(ctx, name, &tracing.SpanOptions{Kind: tracing.SpanKindClient})
	return ctx, func(err error) { endSpan(span, err) }
}

func (d *Diagnostics) StartLifecycle(ctx context.Context, component, action string) (context.Context, func(error)) {
	if d == nil {
		return ctx, func(error) {}
	}
	ctx, end := d.StartAttempt(ctx, "EventHubs."+component+"."+action)
	d.logLifecycle(ctx, component, "started")
	return ctx, func(err error) {
		d.recordLifecycle(ctx, component, outcome(err))
		end(err)
	}
}

func (d *Diagnostics) recordLifecycle(ctx context.Context, component, result string) {
	if d == nil {
		return
	}
	d.logLifecycle(ctx, component, result)
	if d.lifecycle != nil {
		d.lifecycle.Add(ctx, 1, metric.WithAttributes(attribute.String("server.address", d.host),
			attribute.String("messaging.destination.name", d.hub), attribute.String("component", component), attribute.String("outcome", result)))
	}
}

func (d *Diagnostics) logLifecycle(ctx context.Context, component, result string) {
	if d.logger == nil || !d.logger.Enabled(ctx, otellog.EnabledParameters{Severity: otellog.SeverityInfo}) {
		return
	}
	var record otellog.Record
	record.SetTimestamp(time.Now())
	record.SetSeverity(otellog.SeverityInfo)
	record.SetBody(attribute.StringValue("Event Hubs AMQP lifecycle"))
	record.AddAttributes(attribute.String("server.address", d.host), attribute.String("messaging.destination.name", d.hub),
		attribute.String("component", component), attribute.String("outcome", result))
	d.logger.Emit(ctx, record)
}
