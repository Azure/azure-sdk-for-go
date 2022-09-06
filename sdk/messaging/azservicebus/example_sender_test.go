// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ExampleClient_NewSender() {
	sender, err = client.NewSender("<queue or topic>", nil)

	if err != nil {
		panic(err)
	}

	// close the sender when it's no longer needed
	defer sender.Close(context.TODO())
}

func ExampleSender_SendMessage_message() {
	message := &azservicebus.Message{
		Body: []byte("hello, this is a message"),
	}

	err = sender.SendMessage(context.TODO(), message, nil)
	exitOnError("Failed to send message", err)
}

func ExampleSender_SendMessage_messageBatch() {
	batch, err := sender.NewMessageBatch(context.TODO(), nil)
	exitOnError("Failed to create message batch", err)

	// By calling AddMessage multiple times you can add multiple messages into a
	// batch. This can help with message throughput, as you can send multiple
	// messages in a single send.
	err = batch.AddMessage(&azservicebus.Message{Body: []byte("hello world")}, nil)

	// We also support adding AMQPMessages directly to a batch as well
	// batch.AddAMQPMessage(&azservicebus.AMQPMessage{})

	if err != nil {
		switch err {
		case azservicebus.ErrMessageTooLarge:
			// At this point you can do a few things:
			// 1. Ignore this message
			// 2. Send this batch (it's full) and create a new batch.
			//
			// The batch can still be used after this error if you have
			// smaller messages you'd still like to add in.
			fmt.Printf("Failed to add message to batch\n")
		default:
			exitOnError("Error while trying to add message to batch", err)
		}
	}

	// After you add all the messages to the batch you send it using
	// Sender.SendMessageBatch()
	err = sender.SendMessageBatch(context.TODO(), batch, nil)
	exitOnError("Failed to send message batch", err)
}

func ExampleSender_ScheduleMessages() {
	// there are two ways of scheduling messages:
	// 1. Using the `Sender.ScheduleMessages()` function.
	// 2. Setting the `Message.ScheduledEnqueueTime` field on a message.

	// schedule the message to be delivered in an hour.
	sequenceNumbers, err := sender.ScheduleMessages(context.TODO(),
		[]*azservicebus.Message{
			{Body: []byte("hello world")},
		}, time.Now().Add(time.Hour), nil)
	exitOnError("Failed to schedule messages", err)

	err = sender.CancelScheduledMessages(context.TODO(), sequenceNumbers, nil)
	exitOnError("Failed to cancel scheduled messages", err)

	// or you can set the `ScheduledEnqueueTime` field on a message when you send it
	future := time.Now().Add(time.Hour)

	err = sender.SendMessage(context.TODO(),
		&azservicebus.Message{
			Body: []byte("hello world"),
			// schedule the message to be delivered in an hour.
			ScheduledEnqueueTime: &future,
		}, nil)
	exitOnError("Failed to schedule messages using SendMessage", err)
}

func ExampleSender_SendAMQPAnnotatedMessage() {
	// AMQP is the underlying protocol for all interaction with Service Bus.
	// You can, if needed, send and receive messages that have a 1:1 correspondence
	// with an AMQP message. This gives you full control over details that are not
	// exposed via the azservicebus.ReceivedMessage type.

	message := &azservicebus.AMQPAnnotatedMessage{
		Body: azservicebus.AMQPAnnotatedMessageBody{
			// there are three kinds of different body encodings
			// Data, Value and Sequence. See the azservicebus.AMQPMessageBody
			// documentation for more details.
			Value: "hello",
		},
		// full access to fields not normally exposed in azservicebus.Message, like
		// the delivery and message annotations.
		MessageAnnotations: map[any]any{
			"x-opt-message-test": "test-annotation-value",
		},
		DeliveryAnnotations: map[any]any{
			"x-opt-delivery-test": "test-annotation-value",
		},
	}

	err := sender.SendAMQPAnnotatedMessage(context.TODO(), message, nil)

	if err != nil {
		panic(err)
	}
}
