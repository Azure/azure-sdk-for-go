// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/go-amqp"
)

const messagingSystemName = "servicebus"

type Provider = tracing.Provider
type Attribute = tracing.Attribute
type Link = tracing.Link
type Propagator = tracing.Propagator
type Carrier = tracing.Carrier

type Span = tracing.Span

type Tracer struct {
	tracer      tracing.Tracer
	propagator  tracing.Propagator
	destination string
}

type StartSpanOptions struct {
	Tracer        Tracer
	OperationName MessagingOperationName
	Attributes    []Attribute
}

func NewTracer(provider Provider, moduleName, version, hostName, queueOrTopic, subscription string) Tracer {
	t := Tracer{
		tracer:      provider.NewTracer(moduleName, version),
		propagator:  provider.NewPropagator(),
		destination: queueOrTopic,
	}
	t.tracer.SetAttributes(Attribute{Key: AttrMessagingSystem, Value: messagingSystemName},
		Attribute{Key: AttrDestinationName, Value: queueOrTopic})
	if hostName != "" {
		t.tracer.SetAttributes(Attribute{Key: AttrServerAddress, Value: hostName})
	}
	if subscription != "" {
		t.tracer.SetAttributes(Attribute{Key: AttrSubscriptionName, Value: subscription})
	}
	return t
}

func (t *Tracer) SpanFromContext(ctx context.Context) tracing.Span {
	return t.tracer.SpanFromContext(ctx)
}

func (t *Tracer) LinkFromContext(ctx context.Context, attrs ...Attribute) Link {
	return t.tracer.LinkFromContext(ctx, attrs...)
}

func (t *Tracer) Inject(ctx context.Context, message *amqp.Message) {
	t.propagator.Inject(ctx, messageCarrierAdapter(message))
}

func (t *Tracer) Extract(ctx context.Context, message *amqp.Message) context.Context {
	if message != nil {
		ctx = t.propagator.Extract(ctx, messageCarrierAdapter(message))
	}
	return ctx
}

func StartSpan(ctx context.Context, options *StartSpanOptions) (context.Context, func(error)) {
	if options == nil || options.OperationName == "" {
		return ctx, func(error) {}
	}
	attrs := append(options.Attributes, Attribute{Key: AttrOperationName, Value: string(options.OperationName)})

	operationType := getOperationType(options.OperationName)
	if operationType != "" {
		attrs = append(attrs, Attribute{Key: AttrOperationType, Value: string(operationType)})
	}
	if operationType == SettleOperationType {
		attrs = append(attrs, Attribute{Key: AttrDispositionStatus, Value: string(options.OperationName)})
	}

	spanKind := getSpanKind(operationType, options.OperationName, options.Attributes)

	tr := options.Tracer
	spanName := string(options.OperationName)
	if tr.destination != "" {
		spanName = fmt.Sprintf("%s %s", options.OperationName, tr.destination)
	}

	return runtime.StartSpan(ctx, spanName, tr.tracer,
		&runtime.StartSpanOptions{
			Kind:       spanKind,
			Attributes: attrs,
		})
}

func getOperationType(operationName MessagingOperationName) MessagingOperationType {
	switch operationName {
	case CreateOperationName:
		return CreateOperationType
	case SendOperationName, ScheduleOperationName, CancelScheduledOperationName:
		return SendOperationType
	case ReceiveOperationName, PeekOperationName, ReceiveDeferredOperationName, RenewMessageLockOperationName:
		return ReceiveOperationType
	case AbandonOperationName, CompleteOperationName, DeferOperationName, DeadLetterOperationName:
		return SettleOperationType
	default:
		return ""
	}
}

func getSpanKind(operationType MessagingOperationType, operationName MessagingOperationName, attrs []Attribute) SpanKind {
	isBatch := isBatchOperation(attrs)
	isSession := isSessionOperation(operationName)
	switch {
	case operationType == CreateOperationType, operationType == SendOperationType && !isBatch:
		return SpanKindProducer
	case operationType == SendOperationType, operationType == ReceiveOperationType, operationType == SettleOperationType, isSession:
		return SpanKindClient
	default:
		return SpanKindInternal
	}
}

func isSessionOperation(operationName MessagingOperationName) bool {
	switch operationName {
	case AcceptSessionOperationName, SetSessionStateOperationName, GetSessionStateOperationName, RenewSessionLockOperationName:
		return true
	default:
		return false
	}
}

func isBatchOperation(attrs []Attribute) bool {
	for _, attr := range attrs {
		if attr.Key == AttrBatchMessageCount {
			return true
		}
	}
	return false
}
