package azservicebus

import "github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"

type Message struct {
	Body []byte
	ID   string
}

type ReceivedMessage struct {
	Message
	LockToken      string
	SequenceNumber int64

	legacyMessage *internal.Message
}
