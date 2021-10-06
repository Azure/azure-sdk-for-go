// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"github.com/Azure/azure-amqp-common-go/v3/uuid"
	"github.com/Azure/go-amqp"
)

type (
	// MessageBatch represents a batch of messages to send to Service Bus in a single message
	MessageBatch struct {
		marshaledMessages [][]byte
		batchEnvelope     *amqp.Message
		maxBytes          int
		size              int
	}
)

const (
	batchMessageFormat uint32 = 0x80013700

	// TODO: should be calculated, not just a constant.
	batchMessageWrapperSize = 100
)

// NewMessageBatch builds a new message batch with a default standard max message size
func newMessageBatch(maxBytes int) *MessageBatch {
	mb := &MessageBatch{
		maxBytes: maxBytes,
	}

	return mb
}

// Add adds a message to the batch if the message will not exceed the max size of the batch
// This function will return:
// (true, nil) if the message was added.
// (false, nil) if the message was too large to fit into the batch.
// (false, err) if an error occurs when adding the message.
func (mb *MessageBatch) Add(m *Message) (bool, error) {
	msg := m.toAMQPMessage()

	if msg.Properties.MessageID == nil || msg.Properties.MessageID == "" {
		uid, err := uuid.NewV4()
		if err != nil {
			return false, err
		}
		msg.Properties.MessageID = uid.String()
	}

	// if mb.SessionID != nil {
	// 	msg.Properties.GroupID = *mb.SessionID
	// }

	bin, err := msg.MarshalBinary()
	if err != nil {
		return false, err
	}

	if mb.Size()+len(bin) > int(mb.maxBytes) {
		return false, nil
	}

	mb.size += len(bin)

	if len(mb.marshaledMessages) == 0 {
		// first message, store it since we need to copy attributes from it
		// when we send the overall batch message.
		amqpMessage := m.toAMQPMessage()

		// we don't need to hold onto this since it'll get encoded
		// into our marshaledMessages. We just want the metadata.
		amqpMessage.Data = nil
		mb.batchEnvelope = amqpMessage
	}

	mb.marshaledMessages = append(mb.marshaledMessages, bin)
	return true, nil
}

// Size is the number of bytes in the message batch
func (mb *MessageBatch) Size() int {
	// calculated data size + batch message wrapper + data wrapper portions of the message
	return mb.size + batchMessageWrapperSize + (len(mb.marshaledMessages) * 5)
}

// Len returns the # of messages in the batch.
func (mb *MessageBatch) Len() int {
	return len(mb.marshaledMessages)
}

// toAMQPMessage converts this batch into a sendable *amqp.Message
// NOTE: not idempotent!
func (mb *MessageBatch) toAMQPMessage() *amqp.Message {
	mb.batchEnvelope.Data = make([][]byte, len(mb.marshaledMessages))
	mb.batchEnvelope.Format = batchMessageFormat

	copy(mb.batchEnvelope.Data, mb.marshaledMessages)

	return mb.batchEnvelope
}

// MessageType indicates the type of this message. Used for tracing.
func (mb *MessageBatch) messageType() string {
	return "Batch"
}
