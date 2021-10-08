// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ExampleClient_NewSender() {
	sender, err = client.NewSender(queueName) // or topicName
	exitOnError("Failed to create sender", err)

	// Output:
}

func ExampleSender_SendMessage_message() {
	message := &azservicebus.Message{
		Body: []byte("hello, this is a message"),
	}

	err = sender.SendMessage(context.TODO(), message)
	exitOnError("Failed to send message", err)
}

func ExampleSender_SendMessage_messageBatch() {
	batch, err := sender.NewMessageBatch(context.TODO(), nil)
	exitOnError("Failed to create message batch", err)

	messagesToSend := []*azservicebus.Message{
		{Body: []byte("hello world")},
		{Body: []byte("hello world as well")},
	}

	for i := 0; i < len(messagesToSend); i++ {
		added, err := batch.Add(messagesToSend[i])

		if added {
			continue
		}

		if err == nil {
			// At this point you can do a few things:
			// 1. Ignore this message
			// 2. Send this batch (it's full) and create a new batch.
			//
			// The batch can still be used after this error.
			log.Fatal("Failed to add message to batch (batch is full)")
		}

		exitOnError("Error while trying to add message to batch", err)
	}

	// now let's send the batch
	err = sender.SendMessageBatch(context.TODO(), batch)
	exitOnError("Failed to send message batch", err)
}

func ExampleSender_ScheduleMessages() {
	// there are two ways of scheduling messages:
	// 1. Using the `Sender.ScheduleMessages()` function.
	// 2. Setting the `Message.ScheduledEnqueueTime` field on a message.

	// schedule the message to be delivered in an hour.
	sequenceNumbers, err := sender.ScheduleMessages(context.TODO(),
		[]azservicebus.SendableMessage{
			&azservicebus.Message{Body: []byte("hello world")},
		}, time.Now().Add(time.Hour))
	exitOnError("Failed to schedule messages", err)

	err = sender.CancelScheduledMessages(context.TODO(), sequenceNumbers)
	exitOnError("Failed to cancel scheduled messages", err)

	// or you can set the `ScheduledEnqueueTime` field on a message when you send it
	future := time.Now().Add(time.Hour)

	err = sender.SendMessages(context.TODO(),
		[]*azservicebus.Message{
			{
				Body: []byte("hello world"),
				// schedule the message to be delivered in an hour.
				ScheduledEnqueueTime: &future,
			},
		})
	exitOnError("Failed to schedule messages using SendMessages", err)

	// Output:
}
