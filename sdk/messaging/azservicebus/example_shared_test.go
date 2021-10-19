// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/joho/godotenv"
)

// Helper functions to make the main examples compilable

func exitOnError(message string, err error) {
	// these errors/failures are expected since the example
	// code, at the moment, can't run.
	if err == nil {
		return
	}

	log.Panicf("(error in example): %s: %s", message, err.Error())
}

func yourLogicForProcessing(message *azservicebus.ReceivedMessage) {
	log.Printf("Message received")
}

const exampleQueue = "exampleQueue"
const exampleTopic = "exampleTopic"
const exampleSessionQueue = "exampleSessionQueue"
const exampleSubscription = "exampleSubscription"

var connectionString string

var client *azservicebus.Client
var sender *azservicebus.Sender
var receiver *azservicebus.Receiver

var messages []*azservicebus.ReceivedMessage
var err error

var initExamplesOnce sync.Once

func initExamples() {
	initExamplesOnce.Do(func() {
		_ = godotenv.Load()
		connectionString = os.Getenv("SERVICEBUS_CONNECTION_STRING")

		if connectionString == "" {
			log.Printf("SERVICEBUS_CONNECTION_STRING needs to be defined in the environment")
			return
		}

		createExampleEntities()

		ExampleNewClientWithConnectionString()

		sender, err = client.NewSender("exampleQueue") // or topicName
		exitOnError("Failed to create sender", err)

		// send some messages so the receiver tests will be fine running.
		for i := 0; i < 5; i++ {
			err = sender.SendMessage(context.TODO(), &azservicebus.Message{
				Body: []byte("hello world"),
			})
			exitOnError("Failed to send message", err)
		}
	})
}

func createExampleEntities() {
	queueManager, err := internal.NewQueueManagerWithConnectionString(connectionString)

	if err != nil {
		log.Fatalf("Failed to create queue manager : %s", err.Error())
	}

	if _, err := queueManager.Get(context.Background(), exampleQueue); err != nil {
		if _, err := queueManager.Put(context.Background(), "exampleQueue"); err != nil {
			log.Fatalf("Failed to create/update `exampleQueue`: %s", err.Error())
		}
	}

	if _, err := queueManager.Get(context.Background(), exampleSessionQueue); err != nil {
		if _, err := queueManager.Put(context.Background(), exampleSessionQueue, internal.QueueEntityWithRequiredSessions()); err != nil {
			log.Fatalf("Failed to create/update `exampleSessionQueue`: %s", err.Error())
		}
	}

	topicManager, err := internal.NewTopicManagerWithConnectionString(connectionString)

	if err != nil {
		log.Fatalf("Failed to create topic manager : %s", err.Error())
	}

	if _, err := topicManager.Get(context.Background(), exampleTopic); err != nil {
		if _, err := topicManager.Put(context.Background(), exampleTopic); err != nil {
			log.Fatalf("Failed to create/update `exampleTopic`: %s", err.Error())
		}
	}

	subscriptionManager, err := internal.NewSubscriptionManagerForConnectionString("exampleTopic", connectionString)

	if err != nil {
		log.Fatalf("Failed to create subscription manager : %s", err.Error())
	}

	if _, err := subscriptionManager.Get(context.Background(), exampleSubscription); err != nil {
		if _, err := subscriptionManager.Put(context.Background(), exampleSubscription); err != nil {
			log.Fatalf("Failed to create/update `exampleSubscription`: %s", err.Error())
		}
	}
}
