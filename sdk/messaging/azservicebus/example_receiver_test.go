// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"errors"
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

func ExampleReceiver_ReceiveMessages_receiveAndDelete() {
	// ReceiveMessages respects the passed in context, and will gracefully stop
	// receiving when 'ctx' is cancelled.
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	messages, err = receiver.ReceiveMessages(ctx, 10, nil)

	if err != nil {
		panic(err)
	}

	for _, message := range messages {
		// The message body is a []byte. For this example we're just assuming that the body
		// was a string, converted to bytes but any []byte payload is valid.
		var body []byte = message.Body
		fmt.Printf("Message received with body: %s\n", string(body))
		fmt.Printf("Received and completed the message\n")
	}

	err := receiver.Close(ctx)

	if err != nil {
		panic(err)
	}

	// In ReceiveAndDelete mode, any messages stored in the internal cache are available after Close(). To avoid
	// message loss you'll want to loop after closing to ensure the cache is emptied.
	// NOTE: you don't need to do this when using PeekLock, which is the default.
	for {
		messages, err := receiver.ReceiveMessages(context.TODO(), 10, nil)

		if sbErr := (*azservicebus.Error)(nil); errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeClosed {
			// we've read all cached messages.
			break
		} else if err != nil {
			panic(err)
		} else {
			// process messages
			for _, message := range messages {
				var body []byte = message.Body
				fmt.Printf("Message received with body: %s\n", string(body))
				fmt.Printf("Received and completed the message\n")
			}
		}
	}
}

func ExampleReceiver_ReceiveMessages_amqpMessage() {
	// AMQP is the underlying protocol for all interaction with Service Bus.
	// You can, if needed, send and receive messages that have a 1:1 correspondence
	// with an AMQP message. This gives you full control over details that are not
	// exposed via the azservicebus.ReceivedMessage type.

	messages, err := receiver.ReceiveMessages(context.TODO(), 1, nil)

	if err != nil {
		panic(err)
	}

	// NOTE: For this example we'll assume we received at least one message.

	// Every received message carries a RawAMQPMessage.
	var rawAMQPMessage *azservicebus.AMQPAnnotatedMessage = messages[0].RawAMQPMessage

	// All the various body encodings available for AMQP messages are exposed via Body
	_ = rawAMQPMessage.Body.Data
	_ = rawAMQPMessage.Body.Value
	_ = rawAMQPMessage.Body.Sequence

	// delivery and message annotations
	_ = rawAMQPMessage.DeliveryAnnotations
	_ = rawAMQPMessage.MessageAnnotations

	// headers and footers
	_ = rawAMQPMessage.Header
	_ = rawAMQPMessage.Footer

	// Settlement (if in azservicebus.ReceiveModePeekLockMode) stil works on the ReceivedMessage.
	err = receiver.CompleteMessage(context.TODO(), messages[0], nil)

	if err != nil {
		panic(err)
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

func ExampleReceiver_ReceiveMessages_second() {
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
		fmt.Printf("DeadLetter Reason: %s\nDeadLetter Description: %s\n", *deadLetterMessage.DeadLetterReason, *deadLetterMessage.DeadLetterErrorDescription)
		err := deadLetterReceiver.CompleteMessage(context.TODO(), deadLetterMessage, nil)
		if err != nil {
			panic(err)
		}
	}
}

func Example_settleWithLockToken() {
	// This example shows you how to settle a message where you've only preserved the lock token. It's a bit more
	// work, on your part, than settling using the entire message but it does allow you to serialize the lock token
	// and then settle it in another process, or even on a completely separate machine.
	//
	// NOTE: this does not work if you're using Service Bus sessions.

	// ReceiveMessages respects the passed in context, and will gracefully stop
	// receiving when 'ctx' is cancelled.
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	// NOTE: we're receiving a single message, as an example - you can do this with multiple messages.
	messages, err = receiver.ReceiveMessages(ctx, 1, nil)

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

		// You can settle with just a lock token. This allows you to settle on a separate Receiver instance, or even
		// a Receiver instance in a completely separate process.
		//
		// To do this:
		// 1. Create a ReceivedMessage instance, like this:
		completelySeparateMsg := &azservicebus.ReceivedMessage{
			LockToken: message.LockToken,
		}

		// 2a. You can also renew the lock, with just the lock token.
		err = receiver.RenewMessageLock(context.TODO(), completelySeparateMsg, nil)

		if err != nil {
			panic(err)
		}

		// 2b. And settle it like you would any other ReceivedMessage.
		err = receiver.CompleteMessage(context.TODO(), completelySeparateMsg, nil)

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
