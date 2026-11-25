// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"net"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type logCapture struct{ records []sdklog.Record }

func (l *logCapture) Export(_ context.Context, records []sdklog.Record) error {
	for _, record := range records {
		l.records = append(l.records, record.Clone())
	}
	return nil
}
func (*logCapture) Shutdown(context.Context) error   { return nil }
func (*logCapture) ForceFlush(context.Context) error { return nil }

type emptyPrefetch struct{ amqpwrap.AMQPReceiverCloser }

func (emptyPrefetch) Prefetched() *amqp.Message { return nil }

type batchSender struct {
	amqpwrap.AMQPSenderCloser
	err      error
	attempts int
	messages int
}

func (s *batchSender) Send(_ context.Context, m *amqp.Message, _ *amqp.SendOptions) error {
	s.attempts++
	if s.err != nil {
		return s.err
	}
	if s.attempts == 1 {
		return &amqp.Error{Condition: "com.microsoft:server-busy", Description: "private broker details"}
	}
	s.messages += len(m.Data)
	return nil
}
func (*batchSender) LinkName() string       { return "test-sender" }
func (*batchSender) MaxMessageSize() uint64 { return 4096 }

// Exercise the public operations with and without providers, substituting only the AMQP transport.
func TestUnitEventHubsTelemetry(t *testing.T) {
	for _, enabled := range []bool{true, false} {
		name := "disabled"
		if enabled {
			name = "enabled"
		}
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			recorder := tracetest.NewSpanRecorder()
			tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
			reader := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			logs := &logCapture{}
			lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(logs)))
			defer func() {
				require.NoError(t, lp.Shutdown(ctx))
				require.NoError(t, mp.Shutdown(ctx))
				require.NoError(t, tp.Shutdown(ctx))
			}()

			previousMeter, previousLogger := otel.GetMeterProvider(), otel.GetLoggerProvider()
			otel.SetMeterProvider(mp)
			otel.SetLoggerProvider(lp)
			defer otel.SetMeterProvider(previousMeter)
			defer otel.SetLoggerProvider(previousLogger)

			var options TelemetryOptions
			if enabled {
				options = TelemetryOptions{TracingProvider: azotel.NewTracingProvider(tp, nil), EnableMetrics: true, EnableLogs: true}
			}
			const conn = "Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=unused;SharedAccessKey=c2VjcmV0"
			producer, err := NewProducerClientFromConnectionString(conn, "events", &ProducerClientOptions{Telemetry: options, RetryOptions: RetryOptions{RetryDelay: time.Millisecond}})
			require.NoError(t, err)
			consumer, err := NewConsumerClient("test.servicebus.windows.net", "events", DefaultConsumerGroup, fakeTokenCredential{}, &ConsumerClientOptions{Telemetry: options})
			require.NoError(t, err)
			partition, err := consumer.NewPartitionClient("0", nil)
			require.NoError(t, err)

			sender := &batchSender{}
			ns := &internal.FakeNSForPartClient{
				Sender:   &internal.FakeAMQPSender{AMQPSenderCloser: sender},
				Receiver: &internal.FakeAMQPReceiver{AMQPReceiverCloser: emptyPrefetch{}, Messages: fakeMessages(3)},
			}
			producer.links = internal.NewLinks(ns, "", producer.getEntityPath, producer.newEventHubProducerLink).WithDiagnostics(producer.diagnostics, "sender")
			partition.links = internal.NewLinks(ns, "", partition.getEntityPath, partition.newEventHubConsumerLink).WithDiagnostics(partition.diagnostics, "receiver")

			prepareCtx, prepare := tp.Tracer("test").Start(ctx, "batch preparation")
			batch, err := producer.NewEventDataBatch(prepareCtx, nil)
			prepare.End()
			require.NoError(t, err)
			require.Error(t, producer.SendEventDataBatch(ctx, batch, nil))
			for _, body := range []string{"one", "two"} {
				require.NoError(t, batch.AddEventData(&EventData{Body: []byte(body)}, nil))
			}
			require.NoError(t, producer.SendEventDataBatch(ctx, batch, nil))
			require.Equal(t, 2, sender.attempts)
			require.Equal(t, 2, sender.messages)
			events, err := partition.ReceiveEvents(ctx, 2, nil)
			require.NoError(t, err)
			require.Len(t, events, 2)
			sender.err = errors.New("context canceled")
			require.Error(t, producer.SendEventDataBatch(ctx, batch, nil))
			partialCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
			defer cancel()
			events, err = partition.ReceiveEvents(partialCtx, 2, nil)
			require.NoError(t, err)
			require.Len(t, events, 1)
			var data metricdata.ResourceMetrics
			require.NoError(t, reader.Collect(ctx, &data))
			if !enabled {
				require.Empty(t, data.ScopeMetrics)
				require.Len(t, recorder.Ended(), 1) // Only the application span.
				require.Empty(t, logs.records)
				return
			}

			spans := map[string][]sdktrace.ReadOnlySpan{}
			for _, span := range recorder.Ended() {
				spans[span.Name()] = append(spans[span.Name()], span)
			}
			require.Len(t, spans["EventHubs.send"], 3)
			require.Equal(t, "cancelled", spans["EventHubs.send"][2].Status().Description)
			require.Contains(t, spans["EventHubs.send"][2].Attributes(), attribute.Int("messaging.batch.message_count", 2))
			require.Equal(t, codes.Error, spans["EventHubs.send"][0].Status().Code)
			require.Equal(t, "error", spans["EventHubs.send"][0].Status().Description)
			require.Len(t, spans["EventHubs.send.disposition"], 3)
			require.Equal(t, "error", spans["EventHubs.send.disposition"][0].Status().Description)
			for _, attempt := range spans["EventHubs.send.disposition"][:2] {
				require.Equal(t, spans["EventHubs.send"][1].SpanContext().SpanID(), attempt.Parent().SpanID())
			}
			require.Len(t, spans["EventHubs.receive"], 2)
			require.Len(t, spans["EventHubs.receive.wait"], 2)
			require.Equal(t, codes.Unset, spans["EventHubs.receive.wait"][1].Status().Code)
			require.Equal(t, spans["EventHubs.receive"][1].SpanContext().SpanID(), spans["EventHubs.receive.wait"][1].Parent().SpanID())
			require.Len(t, spans["EventHubs.sender.attach"], 1)
			require.Equal(t, prepare.SpanContext().SpanID(), spans["EventHubs.sender.attach"][0].Parent().SpanID())
			for _, name := range []string{"EventHubs.receiver.attach", "EventHubs.receive.wait"} {
				require.NotEmpty(t, spans[name])
				require.Equal(t, spans["EventHubs.receive"][0].SpanContext().SpanID(), spans[name][0].Parent().SpanID())
			}
			require.Len(t, logs.records, 4)
			for i, record := range logs.records {
				component := []string{"sender", "receiver"}[i/2]
				attachments := spans["EventHubs."+component+".attach"]
				require.Len(t, attachments, 1)
				require.Equal(t, attachments[0].SpanContext().SpanID(), record.SpanID())
				require.Equal(t, "Event Hubs AMQP lifecycle", record.Body().AsString())
				attrs := map[string]string{}
				record.WalkAttributes(func(kv attribute.KeyValue) bool { attrs[string(kv.Key)] = kv.Value.AsString(); return true })
				require.Equal(t, component, attrs["component"])
				require.Equal(t, []string{"started", "success"}[i%2], attrs["outcome"])
			}
			totals := map[string]int64{}
			require.Len(t, data.ScopeMetrics, 1)
			for _, instrument := range data.ScopeMetrics[0].Metrics {
				switch values := instrument.Data.(type) {
				case metricdata.Sum[int64]:
					for _, point := range values.DataPoints {
						if instrument.Name == "az.eventhubs.amqp.lifecycle" {
							for key, want := range map[string]string{"server.address": "test.servicebus.windows.net", "messaging.destination.name": "events"} {
								got, ok := point.Attributes.Value(attribute.Key(key))
								require.True(t, ok)
								require.Equal(t, want, got.AsString())
							}
						}
						totals[instrument.Name] += point.Value
					}
				case metricdata.Histogram[float64]:
					for _, point := range values.DataPoints {
						totals[instrument.Name] += int64(point.Count)
					}
				}
			}
			require.Equal(t, map[string]int64{"az.eventhubs.operations": 5, "az.eventhubs.events": 5, "az.eventhubs.operation.duration": 5, "az.eventhubs.amqp.lifecycle": 2}, totals)
		})
	}
}

func TestUnitTelemetryColdReceive(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	defer func() { require.NoError(t, tp.Shutdown(t.Context())) }()
	consumer, err := NewConsumerClientFromConnectionString(
		"Endpoint=sb://test.servicebus.windows.net/;SharedAccessKeyName=unused;SharedAccessKey=c2VjcmV0", "hub", DefaultConsumerGroup,
		&ConsumerClientOptions{
			Telemetry: TelemetryOptions{TracingProvider: azotel.NewTracingProvider(tp, nil)},
			NewWebSocketConn: func(context.Context, WebSocketConnParams) (net.Conn, error) {
				return nil, context.Canceled
			},
		})
	require.NoError(t, err)
	defer func() { require.NoError(t, consumer.Close(t.Context())) }()
	partition, err := consumer.NewPartitionClient("0", nil)
	require.NoError(t, err)
	defer func() { require.NoError(t, partition.Close(t.Context())) }()
	events, err := partition.ReceiveEvents(t.Context(), 1, nil)
	require.Empty(t, events)
	require.ErrorIs(t, err, context.Canceled)
	spans := recorder.Ended()
	require.Len(t, spans, 2)
	require.Equal(t, "EventHubs.connection.open", spans[0].Name())
	require.Equal(t, "EventHubs.receive", spans[1].Name())
	require.Equal(t, spans[1].SpanContext().SpanID(), spans[0].Parent().SpanID())
	for _, span := range spans {
		require.Equal(t, codes.Error, span.Status().Code)
		require.Equal(t, "cancelled", span.Status().Description)
	}
}
