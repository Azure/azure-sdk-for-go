// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"errors"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/go-amqp"
)

// ErrMessageTooLarge is returned when a message cannot fit into a batch when using MessageBatch.Add()
var ErrMessageTooLarge = errors.New("the message could not be added because it is too large for the batch")

type (
	// MessageBatch represents a batch of messages to send to Service Bus in a single message
	MessageBatch struct {
		mu sync.RWMutex

		marshaledMessages [][]byte
		batchEnvelope     *amqp.Message
		maxBytes          uint64
		size              uint64
	}
)

const (
	batchMessageFormat uint32 = 0x80013700

	// TODO: should be calculated, not just a constant.
	batchMessageWrapperSize = uint64(100)
)

// NewMessageBatch builds a new message batch with a default standard max message size
func newMessageBatch(maxBytes uint64) *MessageBatch {
	mb := &MessageBatch{
		maxBytes: maxBytes,
	}

	return mb
}

// AddMessage adds a message to the batch if the message will not exceed the max size of the batch
// Returns:
// - ErrMessageTooLarge if the message cannot fit
// - a non-nil error for other failures
// - nil, otherwise
func (mb *MessageBatch) AddMessage(m *Message) error {
	return mb.addAMQPMessage(m.toAMQPMessage())
}

// NumBytes is the number of bytes in the message batch
func (mb *MessageBatch) NumBytes() uint64 {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	// calculated data size + batch message wrapper + data wrapper portions of the message
	return mb.numBytesNoLock()
}

func (mb *MessageBatch) numBytesNoLock() uint64 {
	return mb.size + batchMessageWrapperSize + (uint64(len(mb.marshaledMessages)) * 5)
}

// NumMessages returns the # of messages in the batch.
func (mb *MessageBatch) NumMessages() int32 {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	return int32(len(mb.marshaledMessages))
}

// toAMQPMessage converts this batch into a sendable *amqp.Message
// NOTE: not idempotent!
func (mb *MessageBatch) toAMQPMessage() *amqp.Message {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.batchEnvelope.Data = make([][]byte, len(mb.marshaledMessages))
	mb.batchEnvelope.Format = batchMessageFormat

	copy(mb.batchEnvelope.Data, mb.marshaledMessages)

	return mb.batchEnvelope
}

func (mb *MessageBatch) addAMQPMessage(msg *amqp.Message) error {
	if msg.Properties.MessageID == nil || msg.Properties.MessageID == "" {
		uid, err := uuid.New()
		if err != nil {
			return err
		}
		msg.Properties.MessageID = uid.String()
	}

	bin, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	mb.mu.Lock()
	defer mb.mu.Unlock()

	if int(mb.numBytesNoLock())+len(bin) > int(mb.maxBytes) {
		return ErrMessageTooLarge
	}

	mb.size += uint64(len(bin))

	if len(mb.marshaledMessages) == 0 {
		// first message, store it since we need to copy attributes from it
		// when we send the overall batch message.
		mb.batchEnvelope = createBatchEnvelope(msg)
	}

	mb.marshaledMessages = append(mb.marshaledMessages, bin)
	return nil
}

// createBatchEnvelope makes a copy of the properties of the message, minus any
// payload fields (like Data, Value or (eventually) Sequence). The data field will be
// filled in with all the messages when the batch is completed.
func createBatchEnvelope(am *amqp.Message) *amqp.Message {
	newAMQPMessage := *am

	newAMQPMessage.Data = nil
	newAMQPMessage.Value = nil

	return &newAMQPMessage
}
