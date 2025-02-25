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
			Name: "test",
			Kind: SpanKindInternal,
			Attributes: []Attribute{
				{Key: AttrOperationName, Value: "test"},
			},
		}, nil).NewTracer("module", "version")}
	subCtx1, endSpan1 := StartSpan(ctx, &StartSpanOptions{Tracer: tr, OperationName: "test"})
	defer endSpan1(nil)
	require.NotEqual(t, ctx, subCtx1)

	// creates a producer span when operation name is SendOperationName
	tr.tracer = tracingvalidator.NewSpanValidator(t, tracingvalidator.SpanMatcher{
		Name: string(SendOperationName),
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
