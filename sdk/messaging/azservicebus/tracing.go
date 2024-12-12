// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
)

const messagingSystemName = "servicebus"

func newTracer(provider tracing.Provider, hostName string) tracing.Tracer {
	tracer := provider.NewTracer(internal.ModuleName, internal.Version)
	if !tracer.Enabled() {
		return tracer
	}

	tracer.SetAttributes(
		tracing.Attribute{Key: tracing.MessagingSystem, Value: messagingSystemName},
	)
	if hostName != "" {
		tracer.SetAttributes(
			tracing.Attribute{Key: tracing.ServerAddress, Value: hostName},
		)
	}

	return tracer
}

func setSenderSpanAttributes(queueOrTopic string, operationName tracing.MessagingOperationName) tracing.SetAttributesFn {
	return func(attrs []tracing.Attribute) []tracing.Attribute {
		attrs = append(attrs,
			tracing.Attribute{Key: tracing.DestinationName, Value: queueOrTopic},
			tracing.Attribute{Key: tracing.OperationType, Value: string(tracing.SendOperationType)},
			tracing.Attribute{Key: tracing.OperationName, Value: string(operationName)},
		)
		return attrs
	}
}

func setReceiverSpanAttributes(entityPath string, operationName tracing.MessagingOperationName) tracing.SetAttributesFn {
	return func(attrs []tracing.Attribute) []tracing.Attribute {
		attrs = setEntityPathAttributes(entityPath)(attrs)
		attrs = append(attrs, tracing.Attribute{Key: tracing.OperationName, Value: string(operationName)})

		if operationName == tracing.CompleteOperationName || operationName == tracing.AbandonOperationName ||
			operationName == tracing.DeadLetterOperationName || operationName == tracing.DeferOperationName {
			attrs = append(attrs, tracing.Attribute{Key: tracing.DispositionStatus, Value: string(operationName)})
			attrs = append(attrs, tracing.Attribute{Key: tracing.OperationType, Value: string(tracing.SettleOperationType)})
		} else {
			attrs = append(attrs, tracing.Attribute{Key: tracing.OperationType, Value: string(tracing.ReceiveOperationType)})
		}

		return attrs
	}
}

func setSessionSpanAttributes(entityPath string, operationName tracing.MessagingOperationName) tracing.SetAttributesFn {
	return func(attrs []tracing.Attribute) []tracing.Attribute {
		attrs = setEntityPathAttributes(entityPath)(attrs)
		attrs = append(attrs, tracing.Attribute{Key: tracing.OperationName, Value: string(operationName)})
		attrs = append(attrs, tracing.Attribute{Key: tracing.OperationType, Value: string(tracing.SessionOperationType)})
		return attrs
	}
}

func setMessageSpanAttributes(message amqpCompatibleMessage) tracing.SetAttributesFn {
	return func(attrs []tracing.Attribute) []tracing.Attribute {
		if message != nil {
			amqpMessage := message.toAMQPMessage()
			if amqpMessage != nil && amqpMessage.Properties != nil {
				if amqpMessage.Properties.MessageID != nil && amqpMessage.Properties.MessageID != "" {
					attrs = append(attrs, tracing.Attribute{Key: tracing.MessageID, Value: amqpMessage.Properties.MessageID})
				}
				if amqpMessage.Properties.CorrelationID != nil {
					attrs = append(attrs, tracing.Attribute{Key: tracing.ConversationID, Value: amqpMessage.Properties.CorrelationID})
				}
			}
		}
		return attrs
	}
}

func setReceivedMessageSpanAttributes(receivedMessage *ReceivedMessage) tracing.SetAttributesFn {
	return func(attrs []tracing.Attribute) []tracing.Attribute {
		if receivedMessage != nil {
			message := receivedMessage.Message()
			attrs = setMessageSpanAttributes(message)(attrs)
			attrs = append(attrs, tracing.Attribute{Key: tracing.DeliveryCount, Value: int64(receivedMessage.DeliveryCount)})
			if receivedMessage.EnqueuedTime != nil {
				attrs = append(attrs, tracing.Attribute{Key: tracing.EnqueuedTime, Value: receivedMessage.EnqueuedTime.Unix()})
			}
		}
		return attrs
	}
}

func setMessageBatchSpanAttributes(size int) tracing.SetAttributesFn {
	return func(attrs []tracing.Attribute) []tracing.Attribute {
		attrs = append(attrs, tracing.Attribute{Key: tracing.BatchMessageCount, Value: int64(size)})
		return attrs
	}
}

func setEntityPathAttributes(entityPath string) tracing.SetAttributesFn {
	return func(attrs []tracing.Attribute) []tracing.Attribute {
		queueOrTopic, subscription := splitEntityPath(entityPath)
		if queueOrTopic != "" {
			attrs = append(attrs, tracing.Attribute{Key: tracing.DestinationName, Value: queueOrTopic})
		}
		if subscription != "" {
			attrs = append(attrs, tracing.Attribute{Key: tracing.SubscriptionName, Value: subscription})
		}
		return attrs
	}
}

func splitEntityPath(entityPath string) (string, string) {
	queueOrTopic := ""
	subscription := ""

	path := strings.Split(entityPath, "/")
	if len(path) == 1 {
		queueOrTopic = path[0]
	} else if len(path) == 2 {
		queueOrTopic = path[0]
		subscription = path[1]
	}
	return queueOrTopic, subscription
}
