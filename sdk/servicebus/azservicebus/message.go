package azservicebus

import (
	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/Azure/go-amqp"
)

const (
	annotationPartitionKey    = "x-opt-partition-key"
	annotationViaPartitionKey = "x-opt-via-partition-key"
)

// Message represents a sendable message, including only user-settable properties.
// A message that has been received from Service Bus will be of type `ReceivedMessage`.
type Message struct {
	Body                    []byte
	ID                      string
	PartitionKey            *string
	TransactionPartitionKey *string

	// TODO: is an empty session ID valid?
	SessionID string
}

// ReceivedMessage represents a message received from Service Bus.
type ReceivedMessage struct {
	Message
	LockToken      *string
	SequenceNumber *int64

	PartitionID *int16

	legacyMessage *internal.Message
}

func (m *Message) toAMQPMessage() (*amqp.Message, error) {
	msg := &amqp.Message{
		Data: [][]byte{m.Body},
		Properties: &amqp.MessageProperties{
			MessageID: m.ID,
			GroupID:   m.SessionID,
		},
		Annotations: amqp.Annotations{},
	}

	if m.PartitionKey != nil {
		msg.Annotations[annotationPartitionKey] = m.PartitionKey
	}

	if m.TransactionPartitionKey != nil {
		msg.Annotations[annotationViaPartitionKey] = m.TransactionPartitionKey
	}

	return msg, nil
}

func convertToReceivedMessage(legacyMessage *internal.Message) *ReceivedMessage {
	var lockToken *string

	if legacyMessage.LockToken != nil {
		tmp := legacyMessage.LockToken.String()
		lockToken = &tmp
	}

	var sequenceNumber *int64
	var partitionKey *string
	var viaPartitionKey *string
	var partitionID *int16

	if legacyMessage.SystemProperties != nil {
		sequenceNumber = legacyMessage.SystemProperties.SequenceNumber
		partitionKey = legacyMessage.SystemProperties.PartitionKey
		viaPartitionKey = legacyMessage.SystemProperties.ViaPartitionKey
		partitionID = legacyMessage.SystemProperties.PartitionID
	}

	rm := &ReceivedMessage{
		// TODO: When we swap out the encoding from the legacy we should also make it so LockToken is simply a string, not expected to be a UUID.
		// Ie, it should be opaque to us.
		LockToken:      lockToken,
		SequenceNumber: sequenceNumber,
		Message: Message{
			Body:                    legacyMessage.Data,
			ID:                      legacyMessage.ID,
			PartitionKey:            partitionKey,
			TransactionPartitionKey: viaPartitionKey,
		},
		legacyMessage: legacyMessage,
		PartitionID:   partitionID,
	}

	return rm
}
