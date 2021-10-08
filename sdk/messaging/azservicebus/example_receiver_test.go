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
		topicName,
		subscriptionName,
		&azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.PeekLock,
		},
	)
	exitOnError("Failed to create Receiver", err)
}

func ExampleClient_NewReceiverForQueue() {
	receiver, err = client.NewReceiverForQueue(
		queueName,
		&azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.PeekLock,
		},
	)
	exitOnError("Failed to create Receiver", err)
}

func ExampleClient_NewReceiverForQueue_deadLetterQueue() {
	receiver, err = client.NewReceiverForQueue(
		queueName,
		&azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.PeekLock,
			SubQueue:    azservicebus.SubQueueDeadLetter,
		},
	)
	exitOnError("Failed to create Receiver for DeadLetterQueue", err)
}

func ExampleClient_NewReceiverForSubscription_deadLetterQueue() {
	receiver, err = client.NewReceiverForSubscription(
		topicName,
		subscriptionName,
		&azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.PeekLock,
			SubQueue:    azservicebus.SubQueueDeadLetter,
		},
	)
	exitOnError("Failed to create Receiver for DeadLetterQueue", err)
}

func ExampleReceiver_ReceiveMessages() {
	// Receive a fixed set of messages. Note that the number of messages
	// to receive and the amount of time to wait are upper bounds.
	messages, err = receiver.ReceiveMessages(context.TODO(),
		// The number of messages to receive. Note this is merely an upper
		// bound. It is possible to get fewer message (or zero), depending
		// on the contents of the remote queue or subscription and network
		// conditions.
		1,
		&azservicebus.ReceiveOptions{
			// This configures the amount of time to wait for messages to arrive.
			// Note that this is merely an upper bound. It is possible to get messages
			// faster than the duration specified.
			MaxWaitTime: 60 * time.Second,
		},
	)

	exitOnError("Failed to receive messages", err)

	for _, message := range messages {
		err = receiver.CompleteMessage(context.TODO(), message)
		fmt.Printf("Received and completed message\n")
		exitOnError("Failed to complete message", err)
	}
}
