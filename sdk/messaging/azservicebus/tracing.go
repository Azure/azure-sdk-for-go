// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
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

func getSenderSpanAttributes(queueOrTopic string, operationName tracing.MessagingOperationName) []tracing.Attribute {
	return append([]tracing.Attribute{},
		tracing.Attribute{Key: tracing.DestinationName, Value: queueOrTopic},
		tracing.Attribute{Key: tracing.OperationType, Value: string(tracing.SendOperationType)},
		tracing.Attribute{Key: tracing.OperationName, Value: string(operationName)},
	)
}

func getReceiverSpanAttributes(entityPath string, operationName tracing.MessagingOperationName) []tracing.Attribute {
	attrs := tracing.GetEntityPathAttributes(entityPath)
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

func getSessionSpanAttributes(entityPath string, operationName tracing.MessagingOperationName) []tracing.Attribute {
	attrs := tracing.GetEntityPathAttributes(entityPath)
	attrs = append(attrs, tracing.Attribute{Key: tracing.OperationName, Value: string(operationName)})
	attrs = append(attrs, tracing.Attribute{Key: tracing.OperationType, Value: string(tracing.SessionOperationType)})
	return attrs
}

func getMessageSpanAttributes(message amqpCompatibleMessage) []tracing.Attribute {
	var attrs []tracing.Attribute
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

func getReceivedMessageSpanAttributes(receivedMessage *ReceivedMessage) []tracing.Attribute {
	var attrs []tracing.Attribute
	if receivedMessage != nil {
		message := receivedMessage.Message()
		attrs = getMessageSpanAttributes(message)
		attrs = append(attrs, tracing.Attribute{Key: tracing.DeliveryCount, Value: int64(receivedMessage.DeliveryCount)})
		if receivedMessage.EnqueuedTime != nil {
			attrs = append(attrs, tracing.Attribute{Key: tracing.EnqueuedTime, Value: receivedMessage.EnqueuedTime.Unix()})
		}
	}
	return attrs
}

func getMessageBatchSpanAttributes(size int) []tracing.Attribute {
	return []tracing.Attribute{{Key: tracing.BatchMessageCount, Value: int64(size)}}
}
