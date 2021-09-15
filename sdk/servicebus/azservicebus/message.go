package azservicebus

import (
	"github.com/Azure/azure-amqp-common-go/v3/uuid"
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
	LockToken      string
	SequenceNumber *int64

	legacyMessage *internal.Message
}

func (m *Message) toAMQPMessage() (*amqp.Message, error) {
	messageID := m.ID

	// TODO: I don't think this should be strictly required. Need to
	// look into why it won't send properly without one.
	if messageID == "" {
		uuid, err := uuid.NewV4()

		if err != nil {
			return nil, err
		}

		messageID = uuid.String()
	}

	msg := &amqp.Message{
		Data: [][]byte{m.Body},
		Properties: &amqp.MessageProperties{
			MessageID: messageID,
			GroupID:   m.SessionID,
		},
		Annotations: amqp.Annotations{},
	}

	if m.PartitionKey != nil {
		msg.Annotations[annotationPartitionKey] = *m.PartitionKey
	}

	if m.TransactionPartitionKey != nil {
		msg.Annotations[annotationViaPartitionKey] = *m.TransactionPartitionKey
	}

	return msg, nil
}

func convertToReceivedMessage(legacyMessage *internal.Message) *ReceivedMessage {
	var lockToken string

	if legacyMessage.LockToken != nil {
		lockToken = legacyMessage.LockToken.String()
	}

	var sequenceNumber *int64
	var partitionKey *string
	var viaPartitionKey *string

	if legacyMessage.SystemProperties != nil {
		sequenceNumber = legacyMessage.SystemProperties.SequenceNumber
		partitionKey = legacyMessage.SystemProperties.PartitionKey
		viaPartitionKey = legacyMessage.SystemProperties.ViaPartitionKey
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
	}

	return rm
}
