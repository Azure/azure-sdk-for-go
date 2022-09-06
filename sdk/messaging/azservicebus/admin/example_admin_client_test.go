// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
)

func ExampleNewClient() {
	// NOTE: If you'd like to authenticate using a Service Bus connection string
	// look at `NewClientWithConnectionString` instead.

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	exitOnError("Failed to create a DefaultAzureCredential", err)

	adminClient, err = admin.NewClient("<ex: myservicebus.servicebus.windows.net>", credential, nil)
	exitOnError("Failed to create ServiceBusClient in example", err)
}

func ExampleNewClientFromConnectionString() {
	// NOTE: If you'd like to authenticate via Azure Active Directory look at
	// the `NewClient` function instead.

	adminClient, err = admin.NewClientFromConnectionString("<connection string>", nil)
	exitOnError("Failed to create ServiceBusClient in example", err)
}

func ExampleClient_CreateQueue() {
	resp, err := adminClient.CreateQueue(context.TODO(), "queue-name", nil)
	exitOnError("Failed to add queue", err)

	// some example properties
	fmt.Printf("Max message delivery count = %d\n", *resp.MaxDeliveryCount)
	fmt.Printf("Lock duration: %s\n", *resp.LockDuration)
}

func ExampleClient_CreateQueue_usingproperties() {
	maxDeliveryCount := int32(10)

	resp, err := adminClient.CreateQueue(context.TODO(), "queue-name", &admin.CreateQueueOptions{
		Properties: &admin.QueueProperties{
			// some example properties
			LockDuration:     to.Ptr("PT1M"),
			MaxDeliveryCount: &maxDeliveryCount,
		},
	})
	exitOnError("Failed to create queue", err)

	// some example properties
	fmt.Printf("Max message delivery count = %d\n", *resp.MaxDeliveryCount)
	fmt.Printf("Lock duration: %s\n", *resp.LockDuration)
}

func ExampleClient_NewListQueuesPager() {
	queuePager := adminClient.NewListQueuesPager(nil)

	for queuePager.More() {
		page, err := queuePager.NextPage(context.TODO())

		if err != nil {
			panic(err)
		}

		for _, queue := range page.Queues {
			fmt.Printf("Queue name: %s, max size in MB: %d\n", queue.QueueName, *queue.MaxSizeInMegabytes)
		}
	}
}

func ExampleClient_NewListQueuesRuntimePropertiesPager() {
	queuePager := adminClient.NewListQueuesRuntimePropertiesPager(nil)

	for queuePager.More() {
		page, err := queuePager.NextPage(context.TODO())

		if err != nil {
			panic(err)
		}

		for _, queue := range page.QueueRuntimeProperties {
			fmt.Printf("Queue name: %s, active messages: %d\n", queue.QueueName, queue.ActiveMessageCount)
		}
	}
}

// NOTE: these are just here to keep the examples succinct.
var adminClient *admin.Client
var err error

func exitOnError(message string, err error) {}
