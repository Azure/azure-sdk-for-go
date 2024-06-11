// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleClient_GetCompletions() {
	azureOpenAIKey := os.Getenv("AOAI_COMPLETIONS_API_KEY")
	modelDeployment := os.Getenv("AOAI_COMPLETIONS_MODEL")

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_COMPLETIONS_ENDPOINT")

	if azureOpenAIKey == "" || modelDeployment == "" || azureOpenAIEndpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	resp, err := client.GetCompletions(context.TODO(), azopenai.CompletionsOptions{
		Prompt:         []string{"What is Azure OpenAI, in 20 words or less"},
		MaxTokens:      to.Ptr(int32(2048)),
		Temperature:    to.Ptr(float32(0.0)),
		DeploymentName: &modelDeployment,
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	for _, choice := range resp.Choices {
		fmt.Fprintf(os.Stderr, "Result: %s\n", *choice.Text)
	}

	// Output:
}

func ExampleClient_GetCompletionsStream() {
	azureOpenAIKey := os.Getenv("AOAI_COMPLETIONS_API_KEY")
	modelDeploymentID := os.Getenv("AOAI_COMPLETIONS_MODEL")

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_COMPLETIONS_ENDPOINT")

	if azureOpenAIKey == "" || modelDeploymentID == "" || azureOpenAIEndpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	resp, err := client.GetCompletionsStream(context.TODO(), azopenai.CompletionsOptions{
		Prompt:         []string{"What is Azure OpenAI, in 20 words or less?"},
		MaxTokens:      to.Ptr(int32(2048)),
		Temperature:    to.Ptr(float32(0.0)),
		DeploymentName: &modelDeploymentID,
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	defer resp.CompletionsStream.Close()

	for {
		entry, err := resp.CompletionsStream.Read()

		if errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "\n*** No more completions ***\n")
			break
		}

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		for _, choice := range entry.Choices {
			fmt.Fprintf(os.Stderr, "Result: %s\n", *choice.Text)
		}
	}

	// Output:
}
