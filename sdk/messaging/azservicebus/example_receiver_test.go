// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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

	if err != nil {
		panic(err)
	}

	for _, message := range messages {
		// For more information about settling messages:
		// https://docs.microsoft.com/azure/service-bus-messaging/message-transfers-locks-settlement#settling-receive-operations
		err = receiver.CompleteMessage(context.TODO(), message, nil)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Received and completed the message\n")
	}
}

func ExampleReceiver_DeadLetterMessage() {
	// Send a message to a queue
	sbMessage := &azservicebus.Message{
		Body: []byte("body of message"),
	}
	err = sender.SendMessage(context.TODO(), sbMessage, nil)
	if err != nil {
		panic(err)
	}
	// Create a receiver
	receiver, err := client.NewReceiverForQueue("myqueue", nil)
	if err != nil {
		panic(err)
	}
	defer receiver.Close(context.TODO())
	// Get the message from a queue
	messages, err := receiver.ReceiveMessages(context.TODO(), 1, nil)
	if err != nil {
		panic(err)
	}
	// Send a message to the dead letter queue
	for _, message := range messages {
		deadLetterOptions := &azservicebus.DeadLetterOptions{
			ErrorDescription: to.Ptr("exampleErrorDescription"),
			Reason:           to.Ptr("exampleReason"),
		}
		err := receiver.DeadLetterMessage(context.TODO(), message, deadLetterOptions)
		if err != nil {
			panic(err)
		}
	}
}

func ExampleReceiver_GetDeadLetterMessage() {
	// Create a dead letter receiver
	deadLetterReceiver, err := client.NewReceiverForQueue(
		"myqueue",
		&azservicebus.ReceiverOptions{
			SubQueue: azservicebus.SubQueueDeadLetter,
		},
	)
	if err != nil {
		panic(err)
	}
	defer deadLetterReceiver.Close(context.TODO())
	// Get messages from the dead letter queue
	deadLetterMessages, err := deadLetterReceiver.ReceiveMessages(context.TODO(), 1, nil)
	if err != nil {
		panic(err)
	}
	// Make messages in the dead letter queue as complete
	for _, deadLetterMessage := range deadLetterMessages {
		fmt.Printf("DeadLetter Reason: %s\nDeadLetter Description: %s\n", *deadLetterMessage.DeadLetterReason, *deadLetter.DeadLetterErrorDescription)
		err := deadLetterReceiver.CompleteMessage(context.TODO(), deadLetterMessage, nil)
		if err != nil {
			panic(err)
		}
	}
}
