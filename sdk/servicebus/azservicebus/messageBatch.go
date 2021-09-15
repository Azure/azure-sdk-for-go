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

type MessageTooLarge interface {
	MessageTooLarge()
}

type errMessageTooLarge struct {
}

func (e errMessageTooLarge) Error() string {
	return "message too large to fit in batch"
}

func (e errMessageTooLarge) NonRetriable()    {}
func (e errMessageTooLarge) MessageTooLarge() {}

// Add adds a message to the batch if the message will not exceed the max size of the batch
// If the message is too large, an error of type 'ErrMessageTooLarge' will be returned.
func (mb *MessageBatch) Add(m *Message) error {
	msg, err := m.toAMQPMessage()
	if err != nil {
		return err
	}

	if msg.Properties.MessageID == nil || msg.Properties.MessageID == "" {
		uid, err := uuid.NewV4()
		if err != nil {
			return err
		}
		msg.Properties.MessageID = uid.String()
	}

	// if mb.SessionID != nil {
	// 	msg.Properties.GroupID = *mb.SessionID
	// }

	bin, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	if mb.Size()+len(bin) > int(mb.maxBytes) {
		return &errMessageTooLarge{}
	}

	mb.size += len(bin)

	if len(mb.marshaledMessages) == 0 {
		// first message, store it since we need to copy attributes from it
		// when we send the overall batch message.
		amqpMessage, err := m.toAMQPMessage()

		if err != nil {
			return err
		}

		// we don't need to hold onto this since it'll get encoded
		// into our marshaledMessages. We just want the metadata.
		amqpMessage.Data = nil
		mb.batchEnvelope = amqpMessage
	}

	mb.marshaledMessages = append(mb.marshaledMessages, bin)
	return nil
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
func (mb *MessageBatch) toAMQPMessage() (*amqp.Message, error) {
	mb.batchEnvelope.Data = make([][]byte, len(mb.marshaledMessages))
	mb.batchEnvelope.Format = batchMessageFormat

	for idx, bytes := range mb.marshaledMessages {
		mb.batchEnvelope.Data[idx] = bytes
	}
	return mb.batchEnvelope, nil
}
