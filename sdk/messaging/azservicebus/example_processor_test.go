// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

var processor *azservicebus.Processor

func ExampleClient_NewProcessorForQueue() {
	processor, err = client.NewProcessorForQueue(
		queueName,
		// NOTE: this is a parameter you'll want to tune. It controls the number of
		// active message `handleMessage` calls that the processor will allow at any time.
		azservicebus.ProcessorWithMaxConcurrentCalls(1),

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
		azservicebus.ProcessorWithReceiveMode(azservicebus.PeekLock),

		// 'AutoComplete' (enabled by default) will use your `handleMessage`
		// function's return value to determine how it should settle your message.
		//
		// A non-nil return error will cause your message to be Abandon()'d.
		// Nil return errors will cause your message to be Complete'd.
		azservicebus.ProcessorWithAutoComplete(true),
	)

	if err != nil {
		log.Printf("Failed to create Processor: %s", err.Error())
	}

	err = processor.Start(context.TODO(), yourMessageHandler, yourErrorHandler)

	if err != nil {
		log.Printf("Failed to start processor")
	}
}

func ExampleClient_NewProcessorForSubscription() {
	processor, err = client.NewProcessorForSubscription(
		topicName,
		subscriptionName,
		// NOTE: this is a parameter you'll want to tune. It controls the number of
		// active message `handleMessage` calls that the processor will allow at any time.
		azservicebus.ProcessorWithMaxConcurrentCalls(1),

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
		azservicebus.ProcessorWithReceiveMode(azservicebus.PeekLock),

		// 'AutoComplete' (enabled by default) will use your `handleMessage`
		// function's return value to determine how it should settle your message.
		//
		// A non-nil return error will cause your message to be Abandon()'d.
		// Nil return errors will cause your message to be Complete'd.
		azservicebus.ProcessorWithAutoComplete(true),
	)

	if err != nil {
		log.Printf("Failed to create Processor: %s", err.Error())
	}

	err = processor.Start(context.TODO(), yourMessageHandler, yourErrorHandler)

	if err != nil {
		log.Printf("Failed to start processor")
	}
}

func ExampleProcessor_Start() {
	handleMessage := func(message *azservicebus.ReceivedMessage) error {
		// This is where your logic for handling messages goes

		yourLogicForProcessing(message)

		// 'AutoComplete' (enabled by default, and controlled by `ProcessorWithAutoComplete`)
		// will use this return value to determine how it should settle your message.
		//
		// Non-nil errors will cause your message to be Abandon()'d.
		// Nil errors will cause your message to be Complete'd.
		return nil
	}

	handleError := func(err error) {
		// handleError will be called on errors that are noteworthy
		// but the Processor internally will continue to attempt to
		// recover.

		// NOTE: errors returned from `handleMessage` above will also be
		// sent here, but do not affect the running of the Processor
		// itself.

		// We'll just print these out, as they're informational and
		// can indicate if there are longer lived problems that we might
		// want to resolve manually (for instance, longer term network
		// outages, or issues affecting your `handleMessage` handler)
		log.Printf("Error: %s", err.Error())
	}

	err = processor.Start(context.TODO(), handleMessage, handleError)

	if err != nil {
		log.Printf("Processor loop has exited: %s", err.Error())
	}
}

func ExampleProcessor_Close() {
	err = processor.Close(context.TODO())

	if err != nil {
		log.Printf("Processor failed to close: %s", err.Error())
	}
}
