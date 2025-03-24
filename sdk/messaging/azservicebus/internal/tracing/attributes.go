// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"time"

	"github.com/Azure/go-amqp"
)

const enqueuedTimeAnnotation = "x-opt-enqueued-time"

func GetMessageIDAttribute(message *amqp.Message) Attribute {
	if message != nil && message.Properties != nil && message.Properties.MessageID != nil && message.Properties.MessageID != "" {
		return Attribute{Key: AttrMessageID, Value: message.Properties.MessageID}
	}
	return Attribute{}
}

func GetMessageSpanAttributes(message *amqp.Message) []Attribute {
	if message == nil {
		return nil
	}

	var attrs []Attribute

	messageIDAttr := GetMessageIDAttribute(message)
	if messageIDAttr.Key != "" {
		attrs = append(attrs, messageIDAttr)
	}

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
