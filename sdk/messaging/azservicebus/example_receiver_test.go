// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ExampleClient_NewReceiverForSubscription() {
	receiver, err = client.NewReceiverForSubscription(
		"exampleTopic",
		"exampleSubscription",
		&azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.ReceiveModePeekLock,
		},
	)
	exitOnError("Failed to create Receiver", err)
}

func ExampleClient_NewReceiverForQueue() {
	receiver, err = client.NewReceiverForQueue(
		"exampleQueue",
		&azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.ReceiveModePeekLock,
		},
	)
	exitOnError("Failed to create Receiver", err)
}

func ExampleClient_NewReceiverForQueue_deadLetterQueue() {
	receiver, err = client.NewReceiverForQueue(
		"exampleQueue",
		&azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.ReceiveModePeekLock,
			SubQueue:    azservicebus.SubQueueDeadLetter,
		},
	)
	exitOnError("Failed to create Receiver for DeadLetterQueue", err)
}

func ExampleClient_NewReceiverForSubscription_deadLetterQueue() {
	receiver, err = client.NewReceiverForSubscription(
		"exampleTopic",
		"exampleSubscription",
		&azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.ReceiveModePeekLock,
			SubQueue:    azservicebus.SubQueueDeadLetter,
		},
	)
	exitOnError("Failed to create Receiver for DeadLetterQueue", err)
}

func ExampleReceiver_ReceiveMessages() {
	// ReceiveMessages respects the passed in context, and will gracefully stop
	// receiving when 'ctx' is cancelled.
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	messages, err = receiver.ReceiveMessages(ctx,
		// The number of messages to receive. Note this is merely an upper
		// bound. It is possible to get fewer message (or zero), depending
		// on the contents of the remote queue or subscription and network
		// conditions.
		1,
		nil,
	)

	exitOnError("Failed to receive messages", err)

	for _, message := range messages {
		err = receiver.CompleteMessage(context.TODO(), message)
		fmt.Printf("Received and completed message\n")
		exitOnError("Failed to complete message", err)
	}
}
