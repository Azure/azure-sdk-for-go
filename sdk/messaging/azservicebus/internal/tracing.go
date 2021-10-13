// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"net/http"
	"os"

	"github.com/devigned/tab"
)

// link/connection recovery spans
const (
	SpanRecover       = "sb.recover"
	SpanRecoverLink   = "sb.recover.link"
	SpanRecoverClient = "sb.recover.client"
)

// authentication
const (
	SpanNegotiateClaim = "sb.auth.negotiateClaim"
)

// settlement
const (
	SpanCompleteMessage = "sb.receiver.complete"
)

// processor
const (
	SpanProcessorLoop    = "sb.processor.main"
	SpanProcessorMessage = "sb.processor.message"
	SpanProcessorClose   = "sb.processor.close"
)

// mgmt client spans
const (
	spanNameRenewLock              = "sb.mgmt.RenewLock"
	spanNameReceiveDeferred        = "sb.mgmt.ReceiveDeferred"
	spanNameSendDisposition        = "sb.mgmt.SendDisposition"
	spanNameScheduleMessage        = "sb.mgmt.Schedule"
	spanNameCancelScheduledMessage = "sb.mgmt.CancelScheduled"
	spanPeekFromSequenceNumber     = "sb.mgmt.PeekSequenceNumber"
	spanNameRecover                = "sb.mgmt.Recover"
	spanNameTryRecover             = "sb.mgmt.TryRecover"
)

// sender spans
const (
	SpanSendMessageFmt string = "sb.SendMessage.%s"
)

func (ns *Namespace) startSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	ApplyComponentInfo(span)
	return ctx, span
}

func (em *EntityManager) startSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	ApplyComponentInfo(span)
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

func (mc *mgmtClient) startSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := startConsumerSpanFromContext(ctx, operationName)
	span.AddAttributes(tab.StringAttribute("message_bus.destination", mc.links.ManagementPath()))
	return ctx, span
}

func (mc *mgmtClient) startProducerSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	ApplyComponentInfo(span)
	span.AddAttributes(
		tab.StringAttribute("span.kind", "producer"),
		tab.StringAttribute("message_bus.destination", mc.links.ManagementPath()),
	)
	return ctx, span
}

func startConsumerSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	ApplyComponentInfo(span)
	span.AddAttributes(tab.StringAttribute("span.kind", "consumer"))
	return ctx, span
}

func ApplyComponentInfo(span tab.Spanner) {
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
