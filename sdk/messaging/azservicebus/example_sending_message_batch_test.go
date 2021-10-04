// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func panicOnError(message string, err error) {
	if err == nil {
		return
	}

	panic(fmt.Sprintf("%s: %s", message, err))
}

func getConnectionString() (string, string) {
	return os.Getenv("SERVICEBUS_CONNECTION_STRING"), os.Getenv("QUEUE_NAME")
}

func ExampleSender_NewMessageBatch() {
	connectionString, queueName := getConnectionString()

	if connectionString == "" || queueName == "" {
		log.Printf("Need a connection string and queue for this example")
		return
	}

	client, err := azservicebus.NewClientWithConnectionString(connectionString)
	panicOnError("Failed to create client", err)

	sender, err := client.NewSender(queueName)
	panicOnError("Failed to create sender", err)

	batch, err := sender.NewMessageBatch(context.TODO())
	panicOnError("Failed to create message batch", err)

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

		panicOnError("Error while trying to add message to batch", err)
	}

	// now let's send the batch
	err = sender.SendMessage(context.TODO(), batch)
	panicOnError("Failed to send message batch", err)
}
