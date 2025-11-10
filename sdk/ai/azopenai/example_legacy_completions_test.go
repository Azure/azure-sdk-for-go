// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
)

// Example_completions demonstrates how to use Azure OpenAI's legacy Completions API.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Send a simple text completion request
// - Handle the completion response
// - Process the generated text output
//
// The example uses environment variables for configuration:
// - AOAI_COMPLETIONS_MODEL: The deployment name of your completions model
// - AOAI_COMPLETIONS_ENDPOINT: Your Azure OpenAI endpoint URL
// - AZURE_OPENAI_API_VERSION: Azure OpenAI service API version to use. See https://learn.microsoft.com/azure/ai-foundry/openai/api-version-lifecycle?tabs=go for information about API versions.
//
// Legacy completions are useful for:
// - Simple text generation tasks
// - Completing partial text
// - Single-turn interactions
// - Basic language generation scenarios
func Example_completions() {
	model := os.Getenv("AOAI_COMPLETIONS_MODEL")
	endpoint := os.Getenv("AOAI_COMPLETIONS_ENDPOINT")
	apiVersion := os.Getenv("AZURE_OPENAI_API_VERSION")

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, apiVersion),
		azure.WithTokenCredential(tokenCredential),
	)

	resp, err := client.Completions.New(context.TODO(), openai.CompletionNewParams{
		Model: openai.CompletionNewParamsModel(model),
		Prompt: openai.CompletionNewParamsPromptUnion{
			OfString: openai.String("What is Azure OpenAI, in 20 words or less"),
		},
		Temperature: openai.Float(0.0),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Fprintf(os.Stderr, "Result: %s\n", resp.Choices[0].Text)
	}

}

// Example_streamCompletions demonstrates streaming responses from the legacy Completions API.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Set up a streaming completion request
// - Process incremental text chunks
// - Handle streaming errors and completion
//
// The example uses environment variables for configuration:
// - AOAI_COMPLETIONS_MODEL: The deployment name of your completions model
// - AOAI_COMPLETIONS_ENDPOINT: Your Azure OpenAI endpoint URL
// - AZURE_OPENAI_API_VERSION: Azure OpenAI service API version to use. See https://learn.microsoft.com/azure/ai-foundry/openai/api-version-lifecycle?tabs=go for information about API versions.
//
// Streaming completions are useful for:
// - Real-time text generation display
// - Reduced latency in responses
// - Interactive text generation
// - Long-form content creation
func Example_streamCompletions() {
	model := os.Getenv("AOAI_COMPLETIONS_MODEL")
	endpoint := os.Getenv("AOAI_COMPLETIONS_ENDPOINT")

	apiVersion := os.Getenv("AZURE_OPENAI_API_VERSION")

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, apiVersion),
		azure.WithTokenCredential(tokenCredential),
	)

	stream := client.Completions.NewStreaming(context.TODO(), openai.CompletionNewParams{
		Model: openai.CompletionNewParamsModel(model),
		Prompt: openai.CompletionNewParamsPromptUnion{
			OfString: openai.String("What is Azure OpenAI, in 20 words or less"),
		},
		MaxTokens:   openai.Int(2048),
		Temperature: openai.Float(0.0),
	})

	for stream.Next() {
		evt := stream.Current()
		if len(evt.Choices) > 0 {
			print(evt.Choices[0].Text)
		}
	}

	if stream.Err() != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}

}
