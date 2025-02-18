// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

// write unit tests for tracing.go
import (
	"context"
	"testing"
	"time"

	tracing2 "github.com/Azure/azure-sdk-for-go/sdk/internal/test/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/stretchr/testify/require"
)

func TestNewTracer(t *testing.T) {
	hostName := "fake.something"
	provider := tracing2.NewSpanValidator(t, tracing2.SpanMatcher{
		Name:   "test_span topic",
		Status: tracing.SpanStatusUnset,
		Attributes: []tracing.Attribute{
			{Key: tracing.MessagingSystem, Value: "servicebus"},
			{Key: tracing.ServerAddress, Value: hostName},
			{Key: tracing.DestinationName, Value: "topic"},
			{Key: tracing.SubscriptionName, Value: "subscription"},
			{Key: tracing.OperationName, Value: "test_span"},
		},
	})
	tracer := newTracer(provider, clientCreds{fullyQualifiedNamespace: hostName}, "topic", "subscription")
	require.NotNil(t, tracer)

	_, endSpan := tracing.StartSpan(context.Background(), &tracing.StartSpanOptions{
		Tracer:        tracer,
		OperationName: "test_span",
	})
	defer func() { endSpan(nil) }()
}

func TestGetMessageSpanAttributes(t *testing.T) {
	messageId := "message-id"
	correlationId := "correlation-id"

	// can parse empty message
	message := &Message{}
	expectedAttrs := []tracing.Attribute{}
	result := getMessageSpanAttributes(message.toAMQPMessage())
	require.ElementsMatch(t, expectedAttrs, result)

	// can parse message with messageId
	message = &Message{
		MessageID: &messageId,
	}
	expectedAttrs = []tracing.Attribute{
		{Key: tracing.MessageID, Value: messageId},
	}
	result = getMessageSpanAttributes(message.toAMQPMessage())
	require.ElementsMatch(t, expectedAttrs, result)

	// can parse message with correlationId
	message = &Message{
		CorrelationID: &correlationId,
	}
	expectedAttrs = []tracing.Attribute{
		{Key: tracing.ConversationID, Value: correlationId},
	}
	result = getMessageSpanAttributes(message.toAMQPMessage())
	require.ElementsMatch(t, expectedAttrs, result)

	// can parse message with both messageId and correlationId
	message = &Message{
		MessageID:     &messageId,
		CorrelationID: &correlationId,
	}
	expectedAttrs = []tracing.Attribute{
		{Key: tracing.MessageID, Value: messageId},
		{Key: tracing.ConversationID, Value: correlationId},
	}
	result = getMessageSpanAttributes(message.toAMQPMessage())
	require.ElementsMatch(t, expectedAttrs, result)
}

func TestGetReceivedMessageSpanAttributes(t *testing.T) {
	messageId := "message-id"
	correlationId := "correlation-id"
	deliveryCount := 1
	enqueuedTime := time.Now()

	// can parse empty message
	receivedMessage := &ReceivedMessage{}
	expectedAttrs := []tracing.Attribute{
		{Key: tracing.DeliveryCount, Value: int64(0)},
	}
	result := getReceivedMessageSpanAttributes(receivedMessage)
	require.ElementsMatch(t, expectedAttrs, result)

	// can parse message with enqueued time
	receivedMessage = &ReceivedMessage{
		MessageID:     messageId,
		CorrelationID: &correlationId,
		DeliveryCount: uint32(deliveryCount),
		EnqueuedTime:  &enqueuedTime,
	}
	expectedAttrs = []tracing.Attribute{
		{Key: tracing.MessageID, Value: messageId},
		{Key: tracing.ConversationID, Value: correlationId},
		{Key: tracing.DeliveryCount, Value: int64(1)},
		{Key: tracing.EnqueuedTime, Value: enqueuedTime.Unix()},
	}
	result = getReceivedMessageSpanAttributes(receivedMessage)
	require.ElementsMatch(t, expectedAttrs, result)
}

func TestGetMessageBatchSpanAttributes(t *testing.T) {
	expectedAttrs := []tracing.Attribute{
		{Key: tracing.BatchMessageCount, Value: int64(1)},
	}
	result := getMessageBatchSpanAttributes(1)
	require.ElementsMatch(t, expectedAttrs, result)
}
