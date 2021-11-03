// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

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
	SpanRenewLock              = "sb.mgmt.RenewLock"
	SpanReceiveDeferred        = "sb.mgmt.ReceiveDeferred"
	SpanSendDisposition        = "sb.mgmt.SendDisposition"
	SpanScheduleMessage        = "sb.mgmt.Schedule"
	SpanCancelScheduledMessage = "sb.mgmt.CancelScheduled"
	SpanPeekFromSequenceNumber = "sb.mgmt.PeekSequenceNumber"
	SpanNameRecover            = "sb.mgmt.Recover"
	SpanTryRecover             = "sb.mgmt.TryRecover"
)

// admin client spans
const (
	SpanGetEntity = "sb.admin.Get"
)

// sender spans
const (
	SpanSendMessageFmt string = "sb.SendMessage.%s"
)

func ApplyRequestInfo(span tab.Spanner, req *http.Request) {
	span.AddAttributes(
		tab.StringAttribute("http.url", req.URL.String()),
		tab.StringAttribute("http.method", req.Method),
	)
}

func ApplyResponseInfo(span tab.Spanner, res *http.Response) {
	if res != nil {
		span.AddAttributes(tab.Int64Attribute("http.status_code", int64(res.StatusCode)))
	}
}

func StartConsumerSpanFromContext(ctx context.Context, operationName string, version string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	ApplyComponentInfo(span, version)
	span.AddAttributes(tab.StringAttribute("span.kind", "consumer"))
	return ctx, span
}

func ApplyComponentInfo(span tab.Spanner, version string) {
	span.AddAttributes(
		tab.StringAttribute("component", "github.com/Azure/azure-sdk-for-go"),
		tab.StringAttribute("version", version),
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
