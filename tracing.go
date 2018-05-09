package servicebus

import (
	"context"
	"os"

	"github.com/opentracing/opentracing-go"
	tag "github.com/opentracing/opentracing-go/ext"
)

//func (sb *serviceBus) startSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
//	span, ctx := opentracing.StartSpanFromContext(ctx, operationName, opts...)
//	ApplyComponentInfo(span)
//	return span, ctx
//}

func (ns *Namespace) startSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName, opts...)
	ApplyComponentInfo(span)
	return span, ctx
}

func (s *sender) startProducerSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName, opts...)
	ApplyComponentInfo(span)
	tag.SpanKindProducer.Set(span)
	tag.MessageBusDestination.Set(span, s.getFullIdentifier())
	return span, ctx
}

func (r *receiver) startConsumerSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName, opts...)
	ApplyComponentInfo(span)
	tag.SpanKindConsumer.Set(span)
	tag.MessageBusDestination.Set(span, r.entityPath)
	return span, ctx
}

func (r *receiver) startConsumerSpanFromWire(ctx context.Context, operationName string, reference opentracing.SpanContext, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	opts = append(opts, opentracing.FollowsFrom(reference))
	span := opentracing.StartSpan(operationName, opts...)
	ctx = opentracing.ContextWithSpan(ctx, span)
	ApplyComponentInfo(span)
	tag.SpanKindConsumer.Set(span)
	tag.MessageBusDestination.Set(span, r.entityPath)
	return span, ctx
}

// ApplyComponentInfo applies eventhub library and network info to the span
func ApplyComponentInfo(span opentracing.Span) {
	tag.Component.Set(span, "github.com/Azure/azure-service-bus-go")
	span.SetTag("version", Version)
	applyNetworkInfo(span)
}

func applyNetworkInfo(span opentracing.Span) {
	hostname, err := os.Hostname()
	if err == nil {
		tag.PeerHostname.Set(span, hostname)
	}
}
