// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
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

func ExampleReceiver_DeleteMessages() {
	count, err := receiver.DeleteMessages(context.TODO(), &azservicebus.DeleteMessagesOptions{
		Count:             4000,
		BeforeEnqueueTime: time.Now(),
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Fprintf(os.Stderr, "Number of messages deleted: %d\n", count)
}

func ExampleReceiver_DeleteMessages_loop() {
	// An example of how to delete messages in a loop.
	now := time.Now()

	for {
		count, err := receiver.DeleteMessages(context.TODO(), &azservicebus.DeleteMessagesOptions{
			BeforeEnqueueTime: now,
		})

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		if count == 0 {
			break
		}
	}
}
