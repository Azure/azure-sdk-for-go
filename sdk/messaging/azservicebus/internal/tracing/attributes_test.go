// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

// write unit tests for go
import (
	"testing"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestGetMessageSpanAttributes(t *testing.T) {
	messageId := "message-id"
	correlationId := "correlation-id"

	testCases := []struct {
		name     string
		message  *amqp.Message
		expected []Attribute
	}{
		{
			name:     "empty message",
			message:  &amqp.Message{},
			expected: []Attribute{},
		},
		{
			name: "message with messageId",
			message: &amqp.Message{
				Properties: &amqp.MessageProperties{
					MessageID: messageId,
				},
			},
			expected: []Attribute{
				{Key: AttrMessageID, Value: messageId},
			},
		},
		{
			name: "message with correlationId",
			message: &amqp.Message{
				Properties: &amqp.MessageProperties{
					CorrelationID: correlationId,
				},
			},
			expected: []Attribute{
				{Key: AttrConversationID, Value: correlationId},
			},
		},
		{
			name: "message with all attributes",
			message: &amqp.Message{
				Properties: &amqp.MessageProperties{
					MessageID:     messageId,
					CorrelationID: correlationId,
				},
			},
			expected: []Attribute{
				{Key: AttrMessageID, Value: messageId},
				{Key: AttrConversationID, Value: correlationId},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetMessageSpanAttributes(tc.message)
			require.ElementsMatch(t, tc.expected, result)
		})
	}
}

func TestGetReceivedMessageSpanAttributes(t *testing.T) {
	messageId := "message-id"
	correlationId := "correlation-id"
	enqueuedTime := time.Now()

	testCases := []struct {
		name     string
		message  *amqp.Message
		expected []Attribute
	}{
		{
			name:     "empty message",
			message:  &amqp.Message{},
			expected: []Attribute{},
		},
		{
			name: "message with messageId and correlationId",
			message: &amqp.Message{
				Properties: &amqp.MessageProperties{
					MessageID:     messageId,
					CorrelationID: correlationId,
				},
			},
			expected: []Attribute{
				{Key: AttrMessageID, Value: messageId},
				{Key: AttrConversationID, Value: correlationId},
			},
		},
		{
			name: "message with all attributes",
			message: &amqp.Message{
				Properties: &amqp.MessageProperties{
					MessageID:     messageId,
					CorrelationID: correlationId,
				},
				Header: &amqp.MessageHeader{
					DeliveryCount: 1,
				},
				Annotations: map[any]any{
					enqueuedTimeAnnotation: enqueuedTime,
				},
			},
			expected: []Attribute{
				{Key: AttrMessageID, Value: messageId},
				{Key: AttrConversationID, Value: correlationId},
				{Key: AttrDeliveryCount, Value: int64(2)},
				{Key: AttrEnqueuedTime, Value: enqueuedTime.Unix()},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetReceivedMessageSpanAttributes(tc.message)
			require.ElementsMatch(t, tc.expected, result)
		})
	}
}

func TestGetMessageBatchSpanAttributes(t *testing.T) {
	expectedAttrs := []Attribute{
		{Key: AttrBatchMessageCount, Value: int64(1)},
	}
	result := GetMessageBatchSpanAttributes(1)
	require.ElementsMatch(t, expectedAttrs, result)
}
