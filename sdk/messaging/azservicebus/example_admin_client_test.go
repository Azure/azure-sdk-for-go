// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ExampleNewAdminClient() {
	// NOTE: If you'd like to authenticate using a Service Bus connection string
	// look at `NewClientWithConnectionString` instead.

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	exitOnError("Failed to create a DefaultAzureCredential", err)

	adminClient, err = azservicebus.NewAdminClient("<ex: myservicebus.servicebus.windows.net>", credential, nil)
	exitOnError("Failed to create ServiceBusClient in example", err)
}

func ExampleNewAdminClientWithConnectionString() {
	// NOTE: If you'd like to authenticate via Azure Active Directory look at
	// the `NewClient` function instead.

	adminClient, err = azservicebus.NewAdminClientWithConnectionString(connectionString, nil)
	exitOnError("Failed to create ServiceBusClient in example", err)
}

func ExampleAdminClient_AddQueue() {
	resp, err := adminClient.AddQueue(context.TODO(), "queue-name")
	exitOnError("Failed to add queue", err)

	fmt.Printf("Queue name: %s", resp.Value.Name)
}

func ExampleAdminClient_AddQueueWithProperties() {
	lockDuration := time.Minute
	maxDeliveryCount := int32(10)

	resp, err := adminClient.AddQueueWithProperties(context.TODO(), &azservicebus.QueueProperties{
		Name: "queue-name",

		// some example properties
		LockDuration:     &lockDuration,
		MaxDeliveryCount: &maxDeliveryCount,
	})
	exitOnError("Failed to create queue", err)

	fmt.Printf("Queue name: %s", resp.Value.Name)
}

func ExampleAdminClient_ListQueues() {
	queuePager := adminClient.ListQueues(nil)

	for queuePager.NextPage(context.TODO()) {
		for _, queue := range queuePager.PageResponse().Value {
			fmt.Printf("Queue name: %s, max size in MB: %d", queue.Name, queue.MaxSizeInMegabytes)
		}
	}

	exitOnError("Failed when listing queues", queuePager.Err())
}

func ExampleAdminClient_ListQueuesRuntimeProperties() {
	queuePager := adminClient.ListQueuesRuntimeProperties(nil)

	for queuePager.NextPage(context.TODO()) {
		for _, queue := range queuePager.PageResponse().Value {
			fmt.Printf("Queue name: %s, active messages: %d", queue.Name, queue.ActiveMessageCount)
		}
	}

	exitOnError("Failed when listing queues runtime properties", queuePager.Err())
}
