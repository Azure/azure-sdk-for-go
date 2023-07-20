// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
)

func ExampleNewClientForOpenAI() {
	keyCredential, err := azopenai.NewKeyCredential("<OpenAI-APIKey>")

	if err != nil {
		// TODO: handle error
	}

	// NOTE: this constructor creates a client that connects to the public OpenAI endpoint.
	// To connect to an Azure OpenAI endpoint, use azopenai.NewClient() or azopenai.NewClientWithyKeyCredential.
	client, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCredential, nil)

	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleNewClient() {
	dac, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		// TODO: handle error
	}

	// NOTE: this constructor creates a client that connects to an Azure OpenAI endpoint.
	// To connect to the public OpenAI endpoint, use azopenai.NewClientForOpenAI
	client, err := azopenai.NewClient("https://<your-azure-openai-host>.openai.azure.com", dac, nil)

	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleNewClientWithKeyCredential() {
	keyCredential, err := azopenai.NewKeyCredential("<Azure-OpenAI-APIKey>")

	if err != nil {
		// TODO: handle error
	}

	// NOTE: this constructor creates a client that connects to an Azure OpenAI endpoint.
	// To connect to the public OpenAI endpoint, use azopenai.NewClientForOpenAI
	client, err := azopenai.NewClientWithKeyCredential("https://<your-azure-openai-host>.openai.azure.com", keyCredential, nil)

	if err != nil {
		// TODO: handle error
	}

	_ = client
}
