// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/conn"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/go-amqp"
)

func newTracer(provider tracing.Provider, creds clientCreds, queueOrTopic, subscription string) tracing.Tracer {
	return tracing.NewTracer(provider, internal.ModuleName, internal.Version, getFullyQualifiedNamespace(creds), queueOrTopic, subscription)
}

// getFullyQualifiedNamespace returns fullyQualifiedNamespace from clientCreds if it is set.
// Otherwise, it parses the connection string and returns the FullyQualifiedNamespace from it.
// If both are empty, it returns an empty string.
func getFullyQualifiedNamespace(creds clientCreds) string {
	if creds.fullyQualifiedNamespace != "" {
		return creds.fullyQualifiedNamespace
	}
	csp, err := conn.ParseConnectionString(creds.connectionString)
	if err != nil {
		return ""
	}
	return csp.FullyQualifiedNamespace
}

func getMessageIDAttribute(message *amqp.Message) []tracing.Attribute {
	var attrs []tracing.Attribute
	if message != nil && message.Properties != nil && message.Properties.MessageID != nil && message.Properties.MessageID != "" {
		attrs = append(attrs, tracing.Attribute{Key: tracing.MessageID, Value: message.Properties.MessageID})
	}
	return attrs
}

func getMessageSpanAttributes(message *amqp.Message) []tracing.Attribute {
	attrs := getMessageIDAttribute(message)
	if message != nil && message.Properties != nil && message.Properties.CorrelationID != nil && message.Properties.CorrelationID != "" {
		attrs = append(attrs, tracing.Attribute{Key: tracing.ConversationID, Value: message.Properties.CorrelationID})
	}
	return attrs
}

func getReceivedMessageSpanAttributes(receivedMessage *ReceivedMessage) []tracing.Attribute {
	var attrs []tracing.Attribute
	if receivedMessage != nil {
		message := receivedMessage.Message().toAMQPMessage()
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
