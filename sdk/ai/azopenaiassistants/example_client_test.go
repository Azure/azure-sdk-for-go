//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleNewClientForOpenAI() {
	keyCredential := azcore.NewKeyCredential("<OpenAI-APIKey>")

	// NOTE: this constructor creates a client that connects to the public OpenAI endpoint.
	// To connect to an Azure OpenAI endpoint, use azopenaiassistants.NewClient() or azopenaiassistants.NewClientWithyKeyCredential.
	client, err := azopenaiassistants.NewClientForOpenAI("https://api.openai.com/v1", keyCredential, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client
}

func ExampleNewClient() {
	dac, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// NOTE: this constructor creates a client that connects to an Azure OpenAI endpoint.
	// To connect to the public OpenAI endpoint, use azopenaiassistants.NewClientForOpenAI
	client, err := azopenaiassistants.NewClient("https://<your-azure-openai-host>.openai.azure.com", dac, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client
}

func ExampleNewClientWithKeyCredential() {
	keyCredential := azcore.NewKeyCredential("<Azure-OpenAI-APIKey>")

	// NOTE: this constructor creates a client that connects to an Azure OpenAI endpoint.
	// To connect to the public OpenAI endpoint, use azopenaiassistants.NewClientForOpenAI
	client, err := azopenaiassistants.NewClientWithKeyCredential("https://<your-azure-openai-host>.openai.azure.com", keyCredential, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client
}
