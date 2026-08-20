// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/stretchr/testify/require"
)

// messageReceiver contains only the receiver operations used by the
// application. The SDK's Receiver satisfies this interface implicitly.
type messageReceiver interface {
	ReceiveMessages(context.Context, int, *azservicebus.ReceiveMessagesOptions) ([]*azservicebus.ReceivedMessage, error)
	CompleteMessage(context.Context, *azservicebus.ReceivedMessage, *azservicebus.CompleteMessageOptions) error
}

type newReceiver func(string, *azservicebus.ReceiverOptions) (messageReceiver, error)

func processMessage(ctx context.Context, createReceiver newReceiver) error {
	receiver, err := createReceiver("orders", nil)
	if err != nil {
		return err
	}

	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	if err != nil {
		return err
	}

	for _, message := range messages {
		if err := receiver.CompleteMessage(ctx, message, nil); err != nil {
			return err
		}
	}

	return nil
}

type fakeReceiver struct {
	messages  []*azservicebus.ReceivedMessage
	completed int
}

func (f *fakeReceiver) ReceiveMessages(context.Context, int, *azservicebus.ReceiveMessagesOptions) ([]*azservicebus.ReceivedMessage, error) {
	return f.messages, nil
}

func (f *fakeReceiver) CompleteMessage(context.Context, *azservicebus.ReceivedMessage, *azservicebus.CompleteMessageOptions) error {
	f.completed++
	return nil
}

var _ messageReceiver = (*azservicebus.Receiver)(nil)
var _ messageReceiver = (*fakeReceiver)(nil)

func Example_unitTesting() {
	fake := &fakeReceiver{
		messages: []*azservicebus.ReceivedMessage{
			{
				Body:      []byte("example"),
				MessageID: "message-id",
			},
		},
	}

	createReceiver := func(string, *azservicebus.ReceiverOptions) (messageReceiver, error) {
		return fake, nil
	}

	err := processMessage(context.TODO(), createReceiver)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

func TestUnitTestingReceiver(t *testing.T) {
	fake := &fakeReceiver{
		messages: []*azservicebus.ReceivedMessage{
			{MessageID: "message-id"},
		},
	}

	createReceiver := func(string, *azservicebus.ReceiverOptions) (messageReceiver, error) {
		return fake, nil
	}

	err := processMessage(context.TODO(), createReceiver)
	require.NoError(t, err)
	require.Equal(t, 1, fake.completed)
}
