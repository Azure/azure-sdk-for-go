// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/joho/godotenv"
)

// Helper functions to make the main examples compilable

func exitOnError(message string, err error) {
	// these errors/failures are expected since the example
	// code, at the moment, can't run.
	if err == nil {
		return
	}

	log.Printf("(error in example): %s: %s", message, err.Error())
}

func yourLogicForProcessing(message *azservicebus.ReceivedMessage) {}

func yourMessageHandler(message *azservicebus.ReceivedMessage) error {
	return nil
}

func yourErrorHandler(err error) {}

var connectionString string
var queueName string
var topicName string
var subscriptionName string
var client *azservicebus.Client
var sender *azservicebus.Sender
var receiver *azservicebus.Receiver
var messages []*azservicebus.ReceivedMessage
var err error

func init() {
	_ = godotenv.Load()
	connectionString, queueName = os.Getenv("SERVICEBUS_CONNECTION_STRING"), os.Getenv("SERVICEBUS_QUEUE")

	if connectionString == "" || queueName == "" {
		log.Printf("Need a connection string and queue for this example")
		return
	}

	ExampleNewClientWithConnectionString()

	sender, err = client.NewSender(queueName) // or topicName
	exitOnError("Failed to create sender", err)

	// send some messages so the receiver tests will be fine running.
	for i := 0; i < 5; i++ {
		err = sender.SendMessage(context.TODO(), &azservicebus.Message{
			Body: []byte("hello world"),
		})
		exitOnError("Failed to send message", err)
	}
}
