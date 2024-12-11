// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

// write unit tests for tracing.go
import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/stretchr/testify/require"
)

func TestNewTracer(t *testing.T) {
	hostName := "fake.something"
	provider := tracing.NewSpanValidator(t, tracing.SpanMatcher{
		Name:   "TestSpan",
		Status: tracing.SpanStatusUnset,
		Attributes: []tracing.Attribute{
			{Key: tracing.MessagingSystem, Value: "servicebus"},
			{Key: tracing.ServerAddress, Value: hostName},
		},
	})
	tracer := newTracer(provider, hostName)
	require.NotNil(t, tracer)
	require.True(t, tracer.Enabled())

	_, endSpan := tracing.StartSpan(context.Background(), "TestSpan", tracer, nil)
	defer func() { endSpan(nil) }()
}

func TestSpanAttributesForMessage(t *testing.T) {
	attrs := getSpanAttributesForMessage(nil)
	require.Empty(t, attrs)

	msgID := "messageID"
	message := &Message{
		MessageID: &msgID,
	}
	attrs = getSpanAttributesForMessage(message)
	require.Equal(t, 1, len(attrs))

	correlationID := "correlationID"
	message.CorrelationID = &correlationID
	attrs = getSpanAttributesForMessage(message)
	require.Equal(t, 2, len(attrs))
}
