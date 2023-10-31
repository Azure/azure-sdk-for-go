// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azwebpubsub_test

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub"
)

func ExampleNewClientWithDefaultAzureCredential() {

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// handle error
	}

	endpoint := os.Getenv("WEBPUBSUB_ENDPOINT")

	if endpoint == "" {
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

func ExampleNewClientWithConnectionString() *azwebpubsub.Client {
	connectionString := os.Getenv("WEBPUBSUB_CONNECTIONSTRING")

	if connectionString == "" {
		return nil
	}

	client, err := azwebpubsub.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	return client
}
