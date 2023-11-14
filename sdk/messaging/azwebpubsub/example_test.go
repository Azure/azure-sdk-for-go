// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azwebpubsub_test

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub"
)

// ExampleNewClient demonstrates how to create a new client with default Azure credentials.
func ExampleNewClient() {

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// handle error
	}

	endpoint := os.Getenv("WEBPUBSUB_ENDPOINT")

	if endpoint == "" {
		return
	}

	hub := os.Getenv("WEBPUBSUB_HUB")
	if hub == "" {
		return
	}
	client, err := azwebpubsub.NewClient(endpoint, cred, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // ignore

	// Output:
}

// ExampleNewClientFromConnectionString demonstrates how to create a new client with connection string.
func ExampleNewClientFromConnectionString() {
	connectionString := os.Getenv("WEBPUBSUB_CONNECTIONSTRING")
	if connectionString == "" {
		return
	}

	client, err := azwebpubsub.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // ignore

	// Output:
}
