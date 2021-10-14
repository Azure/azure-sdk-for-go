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

func ExampleAdminClient_CreateQueue() {
	queueProperties, err := adminClient.CreateQueue(context.TODO(), "queue-name")
	exitOnError("Failed to create queue", err)

	fmt.Printf("Queue name: %s", queueProperties.Name)
}

func ExampleAdminClient_CreateQueueWithProperties() {
	lockDuration := time.Minute
	maxDeliveryCount := int32(10)

	queueProperties, err := adminClient.CreateQueueWithProperties(context.TODO(), &azservicebus.QueueProperties{
		Name: "queue-name",

		// some example properties
		LockDuration:     &lockDuration,
		MaxDeliveryCount: &maxDeliveryCount,
	})
	exitOnError("Failed to create queue", err)

	fmt.Printf("Queue name: %s", queueProperties.Name)
}
