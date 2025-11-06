// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
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
//
// Legacy completions are useful for:
// - Simple text generation tasks
// - Completing partial text
// - Single-turn interactions
// - Basic language generation scenarios
func Example_completions() {
	if !CheckRequiredEnvVars("AOAI_COMPLETIONS_MODEL", "AOAI_COMPLETIONS_ENDPOINT") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	model := os.Getenv("AOAI_COMPLETIONS_MODEL")
	endpoint := os.Getenv("AOAI_COMPLETIONS_ENDPOINT")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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
//
// Streaming completions are useful for:
// - Real-time text generation display
// - Reduced latency in responses
// - Interactive text generation
// - Long-form content creation
func Example_streamCompletions() {
	if !CheckRequiredEnvVars("AOAI_COMPLETIONS_MODEL", "AOAI_COMPLETIONS_ENDPOINT") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	model := os.Getenv("AOAI_COMPLETIONS_MODEL")
	endpoint := os.Getenv("AOAI_COMPLETIONS_ENDPOINT")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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
