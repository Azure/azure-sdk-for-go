// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

var processor *azservicebus.Processor

func ExampleClient_NewProcessorForSubscription() {
	processor, err = client.NewProcessorForSubscription(
		topicName,
		subscriptionName,
		&azservicebus.ProcessorOptions{
			// NOTE: this is a parameter you'll want to tune. It controls the number of
			// active message `handleMessage` calls that the processor will allow at any time.
			MaxConcurrentCalls: 1,
			ReceiveMode:        azservicebus.PeekLock,
			ManualComplete:     false,
		},
	)
	exitOnError("Failed to create Processor", err)
}

func ExampleClient_NewProcessorForQueue() {
	processor, err = client.NewProcessorForQueue(
		queueName,
		&azservicebus.ProcessorOptions{
			// NOTE: this is a parameter you'll want to tune. It controls the number of
			// active message `handleMessage` calls that the processor will allow at any time.
			MaxConcurrentCalls: 1,
			ReceiveMode:        azservicebus.PeekLock,
			ManualComplete:     false,
		},
	)
	exitOnError("Failed to create Processor", err)
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
	exitOnError("Failed to start Processor", err)
}

func ExampleProcessor_Close() {
	err = processor.Close(context.TODO())
	exitOnError("Processor failed to close", err)
}
