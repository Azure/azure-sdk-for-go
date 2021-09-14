package azservicebus

import (
	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/Azure/go-amqp"
)

// Message represents a sendable message, including only user-settable properties.
// A message that has been received from Service Bus will be of type `ReceivedMessage`.
type Message struct {
	Body                    []byte
	ID                      string
	PartitionKey            *string
	TransactionPartitionKey *string

	// TODO: SessionID seems like it should be a pointer (ie: optional)
	SessionID string
}

// ReceivedMessage represents a message received from Service Bus.
type ReceivedMessage struct {
	Message
	LockToken      string
	SequenceNumber int64

	PartitionID *int16

	legacyMessage *internal.Message
}

func (m *Message) toAMQPMessage() (*amqp.Message, error) {
	return &amqp.Message{
		Data: [][]byte{m.Body},
		Properties: &amqp.MessageProperties{
			MessageID: m.ID,
			GroupID:   m.SessionID,
		},
		Annotations: amqp.Annotations{
			"x-opt-partition-key":     m.PartitionKey,
			"x-opt-via-partition-key": m.TransactionPartitionKey,
		},
	}, nil
}

func convertToReceivedMessage(legacyMessage *internal.Message) *ReceivedMessage {
	rm := &ReceivedMessage{
		// TODO: When we swap out the encoding from the legacy we should also make it so LockToken is simply a string, not expected to be a UUID.
		// Ie, it should be opaque to us.
		LockToken:      legacyMessage.LockToken.String(),
		SequenceNumber: *legacyMessage.SystemProperties.SequenceNumber,
		Message: Message{
			Body:                    legacyMessage.Data,
			ID:                      legacyMessage.ID,
			PartitionKey:            legacyMessage.SystemProperties.PartitionKey,
			TransactionPartitionKey: legacyMessage.SystemProperties.ViaPartitionKey,
		},
		legacyMessage: legacyMessage,
		PartitionID:   legacyMessage.SystemProperties.PartitionID,
	}

	return rm
}
