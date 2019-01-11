package servicebus

//	MIT License
//
//	Copyright (c) Microsoft Corporation. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

import (
	"context"
	"net/http"
	"os"

	"github.com/Azure/azure-amqp-common-go/trace"
)

func (ns *Namespace) startSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Spanner) {
	ctx, span := trace.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	return ctx, span
}

func (m *Message) startSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Spanner) {
	ctx, span := trace.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	attrs := []trace.Attribute{trace.StringAttribute("amqp.message.id", m.ID)}
	if m.SessionID != nil {
		attrs = append(attrs, trace.StringAttribute("amqp.session.id", *m.SessionID))
	}
	if m.GroupSequence != nil {
		attrs = append(attrs, trace.Int64Attribute("amqp.sequence_number", int64(*m.GroupSequence)))
	}
	return ctx, span
}

func (em *entityManager) startSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Spanner) {
	ctx, span := trace.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	span.AddAttributes(trace.StringAttribute("span.kind", "client"))
	return ctx, span
}

func applyRequestInfo(span trace.Spanner, req *http.Request) {
	span.AddAttributes(
		trace.StringAttribute("http.url", req.URL.String()),
		trace.StringAttribute("http.method", req.Method),
	)
}

func applyResponseInfo(span trace.Spanner, res *http.Response) {
	if res != nil {
		span.AddAttributes(trace.Int64Attribute("http.status_code", int64(res.StatusCode)))
	}
}

func (e *entity) startSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Spanner) {
	ctx, span := trace.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	return ctx, span
}

func (si sessionIdentifiable) startSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Spanner) {
	ctx, span := trace.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	return ctx, span
}

func (s *Sender) startProducerSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Spanner) {
	ctx, span := trace.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	span.AddAttributes(
		trace.StringAttribute("span.kind", "producer"),
		trace.StringAttribute("message_bus.destination", s.getFullIdentifier()),
	)
	return ctx, span
}

func (r *Receiver) startConsumerSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Spanner) {
	ctx, span := startConsumerSpanFromContext(ctx, operationName)
	span.AddAttributes(trace.StringAttribute("message_bus.destination", r.entityPath))
	return ctx, span
}

func startConsumerSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Spanner) {
	ctx, span := trace.StartSpan(ctx, operationName)
	applyComponentInfo(span)
	span.AddAttributes(trace.StringAttribute("span.kind", "consumer"))
	return ctx, span
}

func applyComponentInfo(span trace.Spanner) {
	span.AddAttributes(
		trace.StringAttribute("component", "github.com/Azure/azure-service-bus-go"),
		trace.StringAttribute("version", Version),
	)
	applyNetworkInfo(span)
}

func applyNetworkInfo(span trace.Spanner) {
	hostname, err := os.Hostname()
	if err == nil {
		span.AddAttributes(
			trace.StringAttribute("peer.hostname", hostname),
		)
	}
}
