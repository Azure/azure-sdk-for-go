// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
)

func ExampleClient_GetCompletions() {
	azureOpenAIKey := os.Getenv("AOAI_API_KEY")
	modelDeploymentID := os.Getenv("AOAI_COMPLETIONS_MODEL_DEPLOYMENT")

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_ENDPOINT")

	if azureOpenAIKey == "" || modelDeploymentID == "" || azureOpenAIEndpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential, err := azopenai.NewKeyCredential(azureOpenAIKey)

	if err != nil {
		// TODO: handle error
	}

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		// TODO: handle error
	}

	resp, err := client.GetCompletions(context.TODO(), azopenai.CompletionsOptions{
		Prompt:       []string{"What is Azure OpenAI, in 20 words or less"},
		MaxTokens:    to.Ptr(int32(2048)),
		Temperature:  to.Ptr(float32(0.0)),
		DeploymentID: modelDeploymentID,
	}, nil)

	if err != nil {
		// TODO: handle error
	}

	for _, choice := range resp.Choices {
		fmt.Fprintf(os.Stderr, "Result: %s\n", *choice.Text)
	}

	// Output:
}

func ExampleClient_GetCompletionsStream() {
	azureOpenAIKey := os.Getenv("AOAI_API_KEY")
	modelDeploymentID := os.Getenv("AOAI_COMPLETIONS_MODEL_DEPLOYMENT")

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_ENDPOINT")

	if azureOpenAIKey == "" || modelDeploymentID == "" || azureOpenAIEndpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential, err := azopenai.NewKeyCredential(azureOpenAIKey)

	if err != nil {
		// TODO: handle error
	}

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		// TODO: handle error
	}

	resp, err := client.GetCompletionsStream(context.TODO(), azopenai.CompletionsOptions{
		Prompt:       []string{"What is Azure OpenAI, in 20 words or less?"},
		MaxTokens:    to.Ptr(int32(2048)),
		Temperature:  to.Ptr(float32(0.0)),
		DeploymentID: modelDeploymentID,
	}, nil)

	if err != nil {
		// TODO: handle error
	}

	for {
		entry, err := resp.CompletionsStream.Read()

		if errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "\n*** No more completions ***\n")
			break
		}

		if err != nil {
			// TODO: handle error
		}

		for _, choice := range entry.Choices {
			fmt.Fprintf(os.Stderr, "Result: %s\n", *choice.Text)
		}
	}

	// Output:
}
