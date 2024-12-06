// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func ExampleNewClient() {
	// NOTE: If you'd like to authenticate using a Service Bus connection string
	// look at `NewClientFromConnectionString` instead.

	// `DefaultAzureCredential` tries several common credential types. For more credential types
	// see this link: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-credential-types.
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

func ExampleNewClient_configuringRetries() {
	// NOTE: If you'd like to authenticate using a Service Bus connection string
	// look at `NewClientFromConnectionString` instead.
	client, err = azservicebus.NewClient(endpoint, tokenCredential, &azservicebus.ClientOptions{
		// NOTE: you don't need to configure these explicitly if you like the defaults.
		// For more information see:
		//  https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus#RetryOptions
		RetryOptions: azservicebus.RetryOptions{
			// MaxRetries specifies the maximum number of attempts a failed operation will be retried
			// before producing an error.
			MaxRetries: 3,
			// RetryDelay specifies the initial amount of delay to use before retrying an operation.
			// The delay increases exponentially with each retry up to the maximum specified by MaxRetryDelay.
			RetryDelay: 4 * time.Second,
			// MaxRetryDelay specifies the maximum delay allowed before retrying an operation.
			// Typically the value is greater than or equal to the value specified in RetryDelay.
			MaxRetryDelay: 120 * time.Second,
		},
	})

	if err != nil {
		panic(err)
	}
}

func Example_enablingLogging() {
	// Required import:
	//   import azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"

	// print log output to stdout
	azlog.SetListener(func(event azlog.Event, s string) {
		fmt.Printf("[%s] %s\n", event, s)
	})

	// pick the set of events to log
	azlog.SetEvents(
		// EventConn is used whenever we create a connection or any links (ie: receivers, senders).
		azservicebus.EventConn,
		// EventAuth is used when we're doing authentication/claims negotiation.
		azservicebus.EventAuth,
		// EventReceiver represents operations that happen on Receivers.
		azservicebus.EventReceiver,
		// EventAdmin is used for operations in the azservicebus/admin.Client
		azservicebus.EventAdmin,
	)
}

func ExampleNewClient_usingCustomEndpoint() {
	// NOTE: If you'd like to authenticate using a Service Bus connection string
	// look at `NewClientFromConnectionString` instead.

	// `DefaultAzureCredential` tries several common credential types. For more credential types
	// see this link: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-credential-types.
	credential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	client, err = azservicebus.NewClient("<ex: myservicebus.servicebus.windows.net>", credential, &azservicebus.ClientOptions{
		// A custom endpoint can be used when you need to connect to a TCP proxy.
		CustomEndpoint: "<address/hostname of TCP proxy>",
	})

	if err != nil {
		panic(err)
	}
}

// fakes for examples
var endpoint string
var tokenCredential azcore.TokenCredential
