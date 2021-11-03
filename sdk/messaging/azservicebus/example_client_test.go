// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ExampleNewClient() {
	// NOTE: If you'd like to authenticate using a Service Bus connection string
	// look at `NewClientFromConnectionString` instead.

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	exitOnError("Failed to create a DefaultAzureCredential", err)

	client, err = azservicebus.NewClient("<ex: myservicebus.servicebus.windows.net>", credential, nil)
	exitOnError("Failed to create ServiceBusClient in example", err)
}

func ExampleNewClientFromConnectionString() {
	// NOTE: If you'd like to authenticate via Azure Active Directory look at
	// the `NewClient` function instead.

	client, err = azservicebus.NewClientFromConnectionString(connectionString, nil)
	exitOnError("Failed to create ServiceBusClient in example", err)
}
