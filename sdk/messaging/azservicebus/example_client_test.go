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

	// For more information about the DefaultAzureCredential:
	// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#NewDefaultAzureCredential
	credential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	client, err = azservicebus.NewClient("<ex: myservicebus.servicebus.windows.net>", credential, nil)

	if err != nil {
		panic(err)
	}
}

func ExampleNewClientFromConnectionString() {
	// NOTE: If you'd like to authenticate via Azure Active Directory look at
	// the `NewClient` function instead.

	client, err = azservicebus.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		panic(err)
	}
}
