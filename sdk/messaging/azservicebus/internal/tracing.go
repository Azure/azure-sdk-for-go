// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"net/http"
	"os"

	"github.com/devigned/tab"
)

func (ns *Namespace) startSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	return ctx, span
}

func (m *Message) startSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	attrs := []tab.Attribute{tab.StringAttribute("amqp.message.id", m.ID)}
	if m.SessionID != nil {
		attrs = append(attrs, tab.StringAttribute("amqp.session.id", *m.SessionID))
	}
	if m.GroupSequence != nil {
		attrs = append(attrs, tab.Int64Attribute("amqp.sequence_number", int64(*m.GroupSequence)))
	}
	span.AddAttributes(attrs...)
	return ctx, span
}

func (em *entityManager) startSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	span.AddAttributes(tab.StringAttribute("span.kind", "client"))
	return ctx, span
}

func applyRequestInfo(span tab.Spanner, req *http.Request) {
	span.AddAttributes(
		tab.StringAttribute("http.url", req.URL.String()),
		tab.StringAttribute("http.method", req.Method),
	)
}

func applyResponseInfo(span tab.Spanner, res *http.Response) {
	if res != nil {
		span.AddAttributes(tab.Int64Attribute("http.status_code", int64(res.StatusCode)))
	}
}

func (e *entity) startSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	span.AddAttributes(tab.StringAttribute("message_bus.destination", e.ManagementPath()))
	return ctx, span
}

func (s *Sender) startProducerSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	span.AddAttributes(
		tab.StringAttribute("span.kind", "producer"),
		tab.StringAttribute("message_bus.destination", s.getFullIdentifier()),
	)
	return ctx, span
}

func (r *Receiver) startConsumerSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := startConsumerSpanFromContext(ctx, operationName)
	span.AddAttributes(tab.StringAttribute("message_bus.destination", r.entityPath))
	return ctx, span
}

func (r *rpcClient) startSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := startConsumerSpanFromContext(ctx, operationName)
	span.AddAttributes(tab.StringAttribute("message_bus.destination", r.ec.ManagementPath()))
	return ctx, span
}

func (r *rpcClient) startProducerSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	span.AddAttributes(
		tab.StringAttribute("span.kind", "producer"),
		tab.StringAttribute("message_bus.destination", r.ec.ManagementPath()),
	)
	return ctx, span
}

func startConsumerSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	span.AddAttributes(tab.StringAttribute("span.kind", "consumer"))
	return ctx, span
}

func applyComponentInfo(span tab.Spanner) {
	span.AddAttributes(
		tab.StringAttribute("component", "github.com/Azure/azure-sdk-for-go"),
		tab.StringAttribute("version", Version),
	)
	applyNetworkInfo(span)
}

func applyNetworkInfo(span tab.Spanner) {
	hostname, err := os.Hostname()
	if err == nil {
		span.AddAttributes(
			tab.StringAttribute("peer.hostname", hostname),
		)
	}
}
