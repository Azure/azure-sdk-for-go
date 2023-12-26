//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package publisher_test

import (
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/publisher"
)

func ExampleNewClient() {
	// ex: https://<topic-name>.<region>.eventgrid.azure.net/api/events
	endpoint := os.Getenv("EVENTGRID_TOPIC_ENDPOINT")

	if endpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	// DefaultAzureCredential is a simplified credential type that tries to authenticate via several
	// different authentication mechanisms. For more control (or more credential types) see the documentation
	// for the azidentity module: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
	tokenCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := publisher.NewClient(endpoint, tokenCred, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client

	// Output:
}

func ExampleNewClientWithSAS() {
	// ex: https://<topic-name>.<region>.eventgrid.azure.net/api/events
	endpoint := os.Getenv("EVENTGRID_TOPIC_ENDPOINT")
	key := os.Getenv("EVENTGRID_TOPIC_KEY")

	if endpoint == "" || key == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	client, err := publisher.NewClientWithSAS(endpoint, azcore.NewSASCredential(key), &publisher.ClientOptions{
		ClientOptions: policy.ClientOptions{
			PerCallPolicies: []policy.Policy{
				dumpFullPolicy{"EventGridEvent"},
			},
		},
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client

	// Output:
}

func ExampleNewClientWithSharedKeyCredential() {
	// ex: https://<topic-name>.<region>.eventgrid.azure.net/api/events
	endpoint := os.Getenv("EVENTGRID_TOPIC_ENDPOINT")
	key := os.Getenv("EVENTGRID_TOPIC_KEY")

	if endpoint == "" || key == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	client, err := publisher.NewClientWithSharedKeyCredential(endpoint, azcore.NewKeyCredential(key), nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client

	// Output:
}
