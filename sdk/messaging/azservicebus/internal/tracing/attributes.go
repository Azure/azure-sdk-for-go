// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"time"

	"github.com/Azure/go-amqp"
)

const enqueuedTimeAnnotation = "x-opt-enqueued-time"

func GetMessageIDAttribute(message *amqp.Message) []Attribute {
	var attrs []Attribute
	if message != nil && message.Properties != nil && message.Properties.MessageID != nil && message.Properties.MessageID != "" {
		attrs = append(attrs, Attribute{Key: AttrMessageID, Value: message.Properties.MessageID})
	}
	return attrs
}

func GetMessageSpanAttributes(message *amqp.Message) []Attribute {
	if message == nil {
		return nil
	}
	attrs := GetMessageIDAttribute(message)
	if message.Properties != nil && message.Properties.CorrelationID != nil && message.Properties.CorrelationID != "" {
		attrs = append(attrs, Attribute{Key: AttrConversationID, Value: message.Properties.CorrelationID})
	}
	return attrs
}

func GetReceivedMessageSpanAttributes(message *amqp.Message) []Attribute {
	if message == nil {
		return nil
	}
	attrs := GetMessageSpanAttributes(message)
	if message.Header != nil {
		attrs = append(attrs, Attribute{Key: AttrDeliveryCount, Value: int64(message.Header.DeliveryCount + 1)})
	}
	if message.Annotations != nil {
		if enqueuedTime, ok := message.Annotations[enqueuedTimeAnnotation]; ok {
			attrs = append(attrs, Attribute{Key: AttrEnqueuedTime, Value: enqueuedTime.(time.Time).Unix()})
		}
	}
	return attrs
}

func GetMessageBatchSpanAttributes(size int) []Attribute {
	return []Attribute{{Key: AttrBatchMessageCount, Value: int64(size)}}
}
