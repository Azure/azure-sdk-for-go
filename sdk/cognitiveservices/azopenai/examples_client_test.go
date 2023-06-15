// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
)

func ExampleNewClientForOpenAI() {
	keyCredential := azopenai.KeyCredential{
		APIKey: "open-ai-apikey",
	}

	client, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCredential, nil)

	if err != nil {
		panic(err)
	}

	_ = client
}

func ExampleNewClient() {
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
	keyCredential := azopenai.KeyCredential{
		APIKey: "Azure OpenAI apikey",
	}

	modelDeploymentID := "model deployment ID"
	client, err := azopenai.NewClientWithKeyCredential("https://<your-azure-openai-host>.openai.azure.com", keyCredential, modelDeploymentID, nil)

	if err != nil {
		panic(err)
	}

	_ = client
}
