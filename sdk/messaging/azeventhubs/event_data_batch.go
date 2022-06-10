// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"errors"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
)

// ErrMessageTooLarge is returned when a message cannot fit into a batch when using MessageBatch.Add()
var ErrMessageTooLarge = errors.New("the message could not be added because it is too large for the batch")

type (
	// EventDataBatch represents a batch of messages to send to Event Hubs in a single message
	EventDataBatch struct {
		mu sync.RWMutex

		marshaledMessages [][]byte
		batchEnvelope     *amqp.Message

		maxBytes    uint64
		currentSize uint64

		partitionID  *string
		partitionKey *string
	}
)

const (
	batchMessageFormat uint32 = 0x80013700
)

// AddEventDataOptions contains optional parameters for the AddEventData function.
type AddEventDataOptions struct {
	// For future expansion
}

// AddEventData adds an EventData to the batch if the message will not exceed the max size of the batch
// Returns:
// - ErrMessageTooLarge if the message cannot fit
// - a non-nil error for other failures
// - nil, otherwise
func (mb *EventDataBatch) AddEventData(ed *EventData, options *AddEventDataOptions) error {
	return mb.addAMQPMessage(ed.toAMQPMessage())
}

// NumBytes is the number of bytes in the message batch
func (mb *EventDataBatch) NumBytes() uint64 {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	return mb.currentSize
}

// NumMessages returns the # of messages in the batch.
func (mb *EventDataBatch) NumMessages() int32 {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	return int32(len(mb.marshaledMessages))
}

// toAMQPMessage converts this batch into a sendable *amqp.Message
// NOTE: not idempotent!
func (mb *EventDataBatch) toAMQPMessage() *amqp.Message {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.batchEnvelope.Data = make([][]byte, len(mb.marshaledMessages))
	mb.batchEnvelope.Format = batchMessageFormat

	if mb.partitionKey != nil {
		mb.batchEnvelope.Annotations[partitionKeyAnnotation] = *mb.partitionKey
	}

	if mb.partitionID != nil {
		mb.batchEnvelope.Annotations[partitionIDAnnotation] = *mb.partitionID
	}

	copy(mb.batchEnvelope.Data, mb.marshaledMessages)

	return mb.batchEnvelope
}

func (mb *EventDataBatch) addAMQPMessage(msg *amqp.Message) error {
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

	if len(mb.marshaledMessages) == 0 {
		// the first message is special - we use its properties and annotations as the
		// actual envelope for the batch message.
		batchEnv, batchEnvLen, err := createBatchEnvelope(msg)

		if err != nil {
			return err
		}

		// (we'll undo this if it turns out the message was too big)
		mb.currentSize = uint64(batchEnvLen)
		mb.batchEnvelope = batchEnv
	}

	actualPayloadSize := calcActualSizeForPayload(bin)

	if mb.currentSize+actualPayloadSize > mb.maxBytes {
		if len(mb.marshaledMessages) == 0 {
			// reset our our properties, this didn't end up being our first message.
			mb.currentSize = 0
			mb.batchEnvelope = nil
		}

		return ErrMessageTooLarge
	}

	mb.currentSize += actualPayloadSize
	mb.marshaledMessages = append(mb.marshaledMessages, bin)

	return nil
}

// createBatchEnvelope makes a copy of the properties of the message, minus any
// payload fields (like Data, Value or Sequence). The data field will be
// filled in with all the messages when the batch is completed.
func createBatchEnvelope(am *amqp.Message) (*amqp.Message, int, error) {
	batchEnvelope := *am

	batchEnvelope.Data = nil
	batchEnvelope.Value = nil
	batchEnvelope.Sequence = nil

	bytes, err := batchEnvelope.MarshalBinary()

	if err != nil {
		return nil, 0, err
	}

	return &batchEnvelope, len(bytes), nil
}

// calcActualSizeForPayload calculates the payload size based
// on overhead from AMQP encoding.
func calcActualSizeForPayload(payload []byte) uint64 {
	const vbin8Overhead = 5
	const vbin32Overhead = 8

	if len(payload) < 256 {
		return uint64(vbin8Overhead + len(payload))
	}

	return uint64(vbin32Overhead + len(payload))
}
