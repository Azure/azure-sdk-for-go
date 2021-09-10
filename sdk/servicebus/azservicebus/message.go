package azservicebus

import "github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"

// Message represents a sendable message, including only user-settable properties.
// A message that has been received from Service Bus will be of type `ReceivedMessage`.
type Message struct {
	Body []byte
	ID   string
}

// ReceivedMessage represents a message received from Service Bus.
type ReceivedMessage struct {
	Message
	LockToken      string
	SequenceNumber int64

	legacyMessage *internal.Message
}
