// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ExampleClient_NewReceiverForQueue() {
	receiver, err = client.NewReceiverForQueue(
		queueName,
		// The receive mode controls when a message is deleted from Service Bus.
		//
		// `azservicebus.PeekLock` is the default. The message is locked, preventing multiple
		// receivers from processing the message at once. You control the lock state of the message
		//  using one of the message settlement functions, processor.CompleteMessage(), which removes
		// it from Service Bus, or processor.AbandonMessage(), which makes it available again.
		//
		// `azservicebus.ReceiveAndDelete` causes Service Bus to remove the message as soon
		// as it's received.
		//
		// More information about receive modes:
		// https://docs.microsoft.com/azure/service-bus-messaging/message-transfers-locks-settlement#settling-receive-operations
		azservicebus.ReceiverWithReceiveMode(azservicebus.PeekLock),
	)

	exitOnError("Failed to create Receiver", err)
}

func ExampleClient_NewReceiverForSubscription() {
	receiver, err = client.NewReceiverForSubscription(
		topicName,
		subscriptionName,
		// The receive mode controls when a message is deleted from Service Bus.
		//
		// `azservicebus.PeekLock` is the default. The message is locked, preventing multiple
		// receivers from processing the message at once. You control the lock state of the message
		//  using one of the message settlement functions, processor.CompleteMessage(), which removes
		// it from Service Bus, or processor.AbandonMessage(), which makes it available again.
		//
		// `azservicebus.ReceiveAndDelete` causes Service Bus to remove the message as soon
		// as it's received.
		//
		// More information about receive modes:
		// https://docs.microsoft.com/azure/service-bus-messaging/message-transfers-locks-settlement#settling-receive-operations
		azservicebus.ReceiverWithReceiveMode(azservicebus.PeekLock),
	)
	exitOnError("Failed to create receiver", err)
}

func ExampleReceiver_ReceiveMessages() {
	// Receive a fixed set of messages. Note that the number of messages
	// to receive and the amount of time to wait are upper bounds.
	messages, err = receiver.ReceiveMessages(context.TODO(),
		// The number of messages to receive. Note this is merely an upper
		// bound. It is possible to get fewer message (or zero), depending
		// on the contents of the remote queue or subscription and network
		// conditions.
		10,
		// This configures the amount of time to wait for messages to arrive.
		// Note that this is merely an upper bound. It is possible to get messages
		// faster than the duration specified.
		azservicebus.ReceiveWithMaxWaitTime(60*time.Second),
	)

	exitOnError("Failed to receive messages", err)

	for _, message := range messages {
		err = receiver.CompleteMessage(context.TODO(), message)
		exitOnError("Failed to complete message", err)
	}
}
