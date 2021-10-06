// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ExampleClient_NewSender() {
	sender, err = client.NewSender(queueName) // or topicName
	exitOnError("Failed to create sender", err)
}

func ExampleSender_SendMessage_message() {
	message := &azservicebus.Message{
		Body: []byte("hello, this is a message"),
	}

	err = sender.SendMessage(context.TODO(), message)
	exitOnError("Failed to send message", err)
}

func ExampleSender_SendMessage_messageBatch() {
	client, err := azservicebus.NewClientWithConnectionString(connectionString, nil)
	exitOnError("Failed to create client", err)

	sender, err := client.NewSender(queueName)
	exitOnError("Failed to create sender", err)

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
	err = sender.SendMessage(context.TODO(), batch)
	exitOnError("Failed to send message batch", err)
}
