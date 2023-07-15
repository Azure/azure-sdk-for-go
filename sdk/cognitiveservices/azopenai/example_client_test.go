// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
)

func ExampleNewClientForOpenAI() {
	// NOTE: this constructor creates a client that connects to the public OpenAI endpoint.
	// To connect to an Azure OpenAI endpoint, use azopenai.NewClient() or azopenai.NewClientWithyKeyCredential.
	keyCredential, err := azopenai.NewKeyCredential("<OpenAI-APIKey>")

	if err != nil {
		panic(err)
	}

	client, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCredential, nil)

	if err != nil {
		panic(err)
	}

	_ = client
}

func ExampleNewClient() {
	// NOTE: this constructor creates a client that connects to an Azure OpenAI endpoint.
	// To connect to the public OpenAI endpoint, use azopenai.NewClientForOpenAI
	dac, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	modelDeploymentID := "model deployment ID"
	client, err := azopenai.NewClient("https://<your-azure-openai-host>.openai.azure.com", dac, modelDeploymentID, nil)

	if err != nil {
		panic(err)
	}

	_ = client
}

func ExampleNewClientWithKeyCredential() {
	// NOTE: this constructor creates a client that connects to an Azure OpenAI endpoint.
	// To connect to the public OpenAI endpoint, use azopenai.NewClientForOpenAI
	keyCredential, err := azopenai.NewKeyCredential("<Azure-OpenAI-APIKey>")

	if err != nil {
		panic(err)
	}

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	modelDeploymentID := "model deployment ID"
	client, err := azopenai.NewClientWithKeyCredential("https://<your-azure-openai-host>.openai.azure.com", keyCredential, modelDeploymentID, nil)

	if err != nil {
		panic(err)
	}

	_ = client
}
