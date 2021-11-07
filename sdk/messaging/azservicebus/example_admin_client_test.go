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

func ExampleNewAdminClientFromConnectionString() {
	// NOTE: If you'd like to authenticate via Azure Active Directory look at
	// the `NewClient` function instead.

	adminClient, err = azservicebus.NewAdminClientFromConnectionString(connectionString, nil)
	exitOnError("Failed to create ServiceBusClient in example", err)
}

func ExampleAdminClient_CreateQueue() {
	resp, err := adminClient.CreateQueue(context.TODO(), "queue-name", nil, nil)
	exitOnError("Failed to add queue", err)

	// some example properties
	fmt.Printf("Max message delivery count = %d", resp.MaxDeliveryCount)
	fmt.Printf("Lock duration: %s", resp.LockDuration)
}

func ExampleAdminClient_CreateQueue_usingproperties() {
	lockDuration := time.Minute
	maxDeliveryCount := int32(10)

	resp, err := adminClient.CreateQueue(context.TODO(), "queue-name", &azservicebus.QueueProperties{
		// some example properties
		LockDuration:     &lockDuration,
		MaxDeliveryCount: &maxDeliveryCount,
	}, nil)
	exitOnError("Failed to create queue", err)

	// some example properties
	fmt.Printf("Max message delivery count = %d", resp.MaxDeliveryCount)
	fmt.Printf("Lock duration: %s", resp.LockDuration)
}

func ExampleAdminClient_ListQueues() {
	queuePager := adminClient.ListQueues(nil)

	for queuePager.NextPage(context.TODO()) {
		for _, queue := range queuePager.PageResponse().Items {
			fmt.Printf("Queue name: %s, max size in MB: %d", queue.QueueName, queue.MaxSizeInMegabytes)
		}
	}

	exitOnError("Failed when listing queues", queuePager.Err())
}

func ExampleAdminClient_ListQueuesRuntimeProperties() {
	queuePager := adminClient.ListQueuesRuntimeProperties(nil)

	for queuePager.NextPage(context.TODO()) {
		for _, queue := range queuePager.PageResponse().Items {
			fmt.Printf("Queue name: %s, active messages: %d", queue.QueueName, queue.ActiveMessageCount)
		}
	}

	exitOnError("Failed when listing queues runtime properties", queuePager.Err())
}
