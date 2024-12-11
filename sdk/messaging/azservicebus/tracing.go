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

func getSpanAttributesForMessage(message *Message) []tracing.Attribute {
	attrs := []tracing.Attribute{}
	if message != nil {
		if message.MessageID != nil {
			attrs = append(attrs, tracing.Attribute{Key: tracing.MessageID, Value: *message.MessageID})
		}
		if message.CorrelationID != nil {
			attrs = append(attrs, tracing.Attribute{Key: tracing.ConversationID, Value: *message.CorrelationID})
		}
	}
	return attrs
}
