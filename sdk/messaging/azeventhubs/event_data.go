// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
)

// EventData is an event that can be sent, using the ProducerClient, to an Event Hub.
type EventData struct {
	// ApplicationProperties can be used to store custom metadata for a message.
	ApplicationProperties map[string]any

	// Body is the payload for a message.
	Body []byte

	// ContentType describes the payload of the message, with a descriptor following
	// the format of Content-Type, specified by RFC2045 (ex: "application/json").
	ContentType *string

	// MessageID is an application-defined value that uniquely identifies
	// the message and its payload. The identifier is a free-form string.
	//
	// If enabled, the duplicate detection feature identifies and removes further submissions
	// of messages with the same MessageId.
	MessageID *string

	// PartitionKey is used with a partitioned entity and enables assigning related messages
	// to the same internal partition. This ensures that the submission sequence order is correctly
	// recorded. The partition is chosen by a hash function in Event Hubs and cannot be chosen
	// directly.
	PartitionKey *string
}

// ReceivedEventData is an event that has been received using the ConsumerClient.
type ReceivedEventData struct {
	// ApplicationProperties can be used to store custom metadata for a message.
	ApplicationProperties map[string]any

	// Body is the payload for a message.
	Body []byte

	// ContentType describes the payload of the message, with a descriptor following
	// the format of Content-Type, specified by RFC2045 (ex: "application/json").
	ContentType *string

	// EnqueuedTime is the UTC time when the message was accepted and stored by Event Hubs.
	EnqueuedTime *time.Time

	// MessageID is an application-defined value that uniquely identifies
	// the message and its payload. The identifier is a free-form string.
	//
	// If enabled, the duplicate detection feature identifies and removes further submissions
	// of messages with the same MessageId.
	MessageID *string

	// PartitionKey is used with a partitioned entity and enables assigning related messages
	// to the same internal partition. This ensures that the submission sequence order is correctly
	// recorded. The partition is chosen by a hash function in Event Hubs and cannot be chosen
	// directly.
	PartitionKey *string

	// Offset is the offset of the event.
	Offset *int64

	// SequenceNumber is a unique number assigned to a message by Event Hubs.
	SequenceNumber int64
}

// Event Hubs custom properties
const (
	// Annotation properties
	partitionKeyAnnotation   = "x-opt-partition-key"
	partitionIDAnnotation    = "x-opt-partition-id"
	sequenceNumberAnnotation = "x-opt-sequence-number"
	offsetNumberAnnotation   = "x-opt-offset"
	enqueuedTimeAnnotation   = "x-opt-enqueued-time"
)

func (e *EventData) toAMQPMessage() *amqp.Message {
	amqpMsg := amqp.NewMessage(e.Body)

	var messageID any

	if e.MessageID != nil {
		messageID = *e.MessageID
	}

	amqpMsg.Properties = &amqp.MessageProperties{
		MessageID: messageID,
	}

	amqpMsg.Properties.ContentType = e.ContentType

	if len(e.ApplicationProperties) > 0 {
		amqpMsg.ApplicationProperties = make(map[string]any)
		for key, value := range e.ApplicationProperties {
			amqpMsg.ApplicationProperties[key] = value
		}
	}

	amqpMsg.Annotations = map[any]any{}

	if e.PartitionKey != nil {
		amqpMsg.Annotations[partitionKeyAnnotation] = *e.PartitionKey
	}

	return amqpMsg
}

// newReceivedEventData creates a received message from an AMQP message.
// NOTE: this converter assumes that the Body of this message will be the first
// serialized byte array in the Data section of the messsage.
func newReceivedEventData(amqpMsg *amqp.Message) *ReceivedEventData {
	re := &ReceivedEventData{}

	if len(amqpMsg.Data) == 1 {
		re.Body = amqpMsg.Data[0]
	}

	if amqpMsg.Properties != nil {
		if id, ok := amqpMsg.Properties.MessageID.(string); ok {
			re.MessageID = &id
		}

		re.ContentType = amqpMsg.Properties.ContentType
	}

	if amqpMsg.ApplicationProperties != nil {
		re.ApplicationProperties = make(map[string]any, len(amqpMsg.ApplicationProperties))
		for key, value := range amqpMsg.ApplicationProperties {
			re.ApplicationProperties[key] = value
		}
	}

	if amqpMsg.Annotations != nil {
		if sequenceNumber, ok := amqpMsg.Annotations[sequenceNumberAnnotation]; ok {
			re.SequenceNumber = sequenceNumber.(int64)
		}

		if partitionKey, ok := amqpMsg.Annotations[partitionKeyAnnotation]; ok {
			re.PartitionKey = to.Ptr(partitionKey.(string))
		}

		if enqueuedTime, ok := amqpMsg.Annotations[enqueuedTimeAnnotation]; ok {
			t := enqueuedTime.(time.Time)
			re.EnqueuedTime = &t
		}

		if offset, ok := amqpMsg.ApplicationProperties[offsetNumberAnnotation]; ok {
			if offsetInt64, ok := offset.(int64); ok {
				re.Offset = to.Ptr(int64(offsetInt64))
			}
		}
	}

	return re
}
