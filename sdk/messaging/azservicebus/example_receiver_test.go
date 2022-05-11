// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"errors"
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

	// close the receiver when it's no longer needed
	defer receiver.Close(context.TODO())
}

func ExampleClient_NewReceiverForQueue() {
	receiver, err = client.NewReceiverForQueue(
		"exampleQueue",
		&azservicebus.ReceiverOptions{
			ReceiveMode: azservicebus.ReceiveModePeekLock,
		},
	)
	exitOnError("Failed to create Receiver", err)

	// close the receiver when it's no longer needed
	defer receiver.Close(context.TODO())
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

	// close the receiver when it's no longer needed
	defer receiver.Close(context.TODO())
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

	// close the receiver when it's no longer needed
	defer receiver.Close(context.TODO())
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
		// The message body is a []byte. For this example we're just assuming that the body
		// was a string, converted to bytes but any []byte payload is valid.
		var body []byte = message.Body
		fmt.Printf("Message received with body: %s\n", string(body))

		// For more information about settling messages:
		// https://docs.microsoft.com/azure/service-bus-messaging/message-transfers-locks-settlement#settling-receive-operations
		err = receiver.CompleteMessage(context.TODO(), message, nil)

		if err != nil {
			var sbErr *azservicebus.Error

			if errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeLockLost {
				// The message lock has expired. This isn't fatal for the client, but it does mean
				// that this message can be received by another Receiver (or potentially this one!).
				fmt.Printf("Message lock expired\n")

				// You can extend the message lock by calling receiver.RenewMessageLock(msg) before the
				// message lock has expired.
				continue
			}

			panic(err)
		}

		fmt.Printf("Received and completed the message\n")
	}
}
