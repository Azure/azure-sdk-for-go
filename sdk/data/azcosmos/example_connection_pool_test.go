// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos_test

import (
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func ExampleNewHTTPClient() {
	httpClient, err := azcosmos.NewHTTPClient(&azcosmos.HTTPClientOptions{
		MaxConnectionsPerEndpoint: 64,
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	client, err := azcosmos.NewClientFromConnectionString(os.Getenv("AZURE_COSMOS_CONNECTION_STRING"), &azcosmos.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: httpClient},
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	defer httpClient.CloseIdleConnections()
	defer client.Close()
}
