// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/tracingvalidator"
	"github.com/stretchr/testify/require"
)

func TestStartSpan(t *testing.T) {
	// no-op when StartSpanOptions is nil
	ctx := context.Background()
	subCtx, _ := StartSpan(ctx, nil)
	require.Equal(t, ctx, subCtx)

	// no-op when StartSpanOptions is empty
	subCtx, _ = StartSpan(ctx, &StartSpanOptions{})
	require.Equal(t, ctx, subCtx)

	// no-op when SpanName is empty
	subCtx, _ = StartSpan(ctx, &StartSpanOptions{OperationName: ""})
	require.Equal(t, ctx, subCtx)

	// creates a span when both tracer and SpanName are set
	tr := Tracer{
		tracer: tracingvalidator.NewSpanValidator(t, tracingvalidator.SpanMatcher{
			Name: "test queue",
			Kind: SpanKindInternal,
			Attributes: []Attribute{
				{Key: AttrOperationName, Value: "test"},
			},
		}, nil).NewTracer("module", "version"),
		destination: "queue"}
	subCtx1, endSpan1 := StartSpan(ctx, &StartSpanOptions{Tracer: tr, OperationName: "test"})
	defer endSpan1(nil)
	require.NotEqual(t, ctx, subCtx1)

	// creates a producer span when operation name is SendOperationName
	tr.tracer = tracingvalidator.NewSpanValidator(t, tracingvalidator.SpanMatcher{
		Name: "send queue",
		Kind: SpanKindProducer,
		Attributes: []Attribute{
			{Key: AttrOperationName, Value: string(SendOperationName)},
			{Key: AttrOperationType, Value: string(SendOperationType)},
		},
	}, nil).NewTracer("module", "version")
	subCtx2, endSpan2 := StartSpan(ctx, &StartSpanOptions{Tracer: tr, OperationName: SendOperationName})
	defer endSpan2(nil)
	require.NotEqual(t, ctx, subCtx2)
}

func TestGetOperationType(t *testing.T) {
	// returns CreateOperationType when operation name is CreateOperationName
	require.Equal(t, CreateOperationType, getOperationType(CreateOperationName))

	// returns SendOperationType when operation name is SendOperationName
	require.Equal(t, SendOperationType, getOperationType(SendOperationName))

	// returns ReceiveOperationType when operation name is ReceiveOperationName
	require.Equal(t, ReceiveOperationType, getOperationType(ReceiveOperationName))

	// returns SettleOperationType when operation name is SettleOperationName
	require.Equal(t, SettleOperationType, getOperationType(CompleteOperationName))
}

func TestGetSpanKind(t *testing.T) {
	// returns SpanKindProducer when operation type is CreateOperationType
	require.Equal(t, SpanKindProducer, getSpanKind(CreateOperationType, CreateOperationName, nil))

	// returns SpanKindProducer when operation type is SendOperationType and not a batch operation
	require.Equal(t, SpanKindProducer, getSpanKind(SendOperationType, SendOperationName, nil))

	// returns SpanKindClient when operation type is SendOperationType and a batch operation
	require.Equal(t, SpanKindClient, getSpanKind(SendOperationType, SendOperationName, []Attribute{{Key: AttrBatchMessageCount, Value: "1"}}))

	// returns SpanKindClient when operation type is ReceiveOperationType
	require.Equal(t, SpanKindClient, getSpanKind(ReceiveOperationType, ReceiveOperationName, nil))

	// returns SpanKindClient when operation type is SettleOperationType
	require.Equal(t, SpanKindClient, getSpanKind(SettleOperationType, CompleteOperationName, nil))

	// returns SpanKindClient with operation name is a session operation
	require.Equal(t, SpanKindClient, getSpanKind("", AcceptSessionOperationName, nil))

	// returns SpanKindInternal when operation type is unknown
	require.Equal(t, SpanKindInternal, getSpanKind("", "unknown", nil))
}
