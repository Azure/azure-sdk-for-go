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

func setMessageSpanAttributes(message amqpCompatibleMessage) tracing.SetAttributesFn {
	return func(attrs []tracing.Attribute) []tracing.Attribute {
		if message != nil {
			amqpMessage := message.toAMQPMessage()
			if amqpMessage != nil && amqpMessage.Properties != nil {
				if amqpMessage.Properties.MessageID != nil {
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

func setMessageBatchSpanAttributes(size int) tracing.SetAttributesFn {
	return func(attrs []tracing.Attribute) []tracing.Attribute {
		attrs = append(attrs, tracing.Attribute{Key: tracing.BatchMessageCount, Value: int64(size)})
		return attrs
	}
}
