// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiextensions"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

// This example demonstrates how to use Azure OpenAI extensions to extract
// Azure-specific information from responses, such as content filter results
func Example_usingAzureContentFiltering() {
	endpoint := os.Getenv("AOAI_ENDPOINT")
	model := os.Getenv("AOAI_MODEL")

	if endpoint == "" || model == "" {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not running example.")
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-08-01-preview"),
		azure.WithTokenCredential(tokenCredential),
	)

	// Standard OpenAI chat completion request
	chatParams := openai.ChatCompletionNewParams{
		Model:     openai.F(model),
		MaxTokens: openai.Int(256),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionMessageParam{
				Role:    openai.F(openai.ChatCompletionMessageParamRoleUser),
				Content: openai.F[any]("Explain briefly how solar panels work"),
			},
		}),
	}

	resp, err := client.Chat.Completions.New(
		context.TODO(),
		chatParams,
	)

	// Check if there's a content filter error
	var contentErr *azopenaiextensions.ContentFilterError
	if azopenaiextensions.ExtractContentFilterError(err, &contentErr) {
		fmt.Fprintf(os.Stderr, "Content was filtered by Azure OpenAI:\n")

		if contentErr.Hate != nil && contentErr.Hate.Filtered != nil && *contentErr.Hate.Filtered {
			fmt.Fprintf(os.Stderr, "- Hate content was filtered\n")
		}

		if contentErr.Violence != nil && contentErr.Violence.Filtered != nil && *contentErr.Violence.Filtered {
			fmt.Fprintf(os.Stderr, "- Violent content was filtered\n")
		}

		if contentErr.Sexual != nil && contentErr.Sexual.Filtered != nil && *contentErr.Sexual.Filtered {
			fmt.Fprintf(os.Stderr, "- Sexual content was filtered\n")
		}

		if contentErr.SelfHarm != nil && contentErr.SelfHarm.Filtered != nil && *contentErr.SelfHarm.Filtered {
			fmt.Fprintf(os.Stderr, "- Self-harm content was filtered\n")
		}

		return
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	if len(resp.Choices) == 0 {
		fmt.Fprintf(os.Stderr, "No choices returned in the response, the model may have failed to generate content\n")
		return
	}

	// Access the Azure-specific content filter results from the response
	azureChatChoice := azopenaiextensions.ChatCompletionChoice(resp.Choices[0])
	contentFilterResults, err := azureChatChoice.ContentFilterResults()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	} else if contentFilterResults != nil {
		fmt.Fprintf(os.Stderr, "Content Filter Results:\n")

		if contentFilterResults.Hate != nil && contentFilterResults.Hate.Severity != nil {
			fmt.Fprintf(os.Stderr, "- Hate severity: %s\n", *contentFilterResults.Hate.Severity)
		}

		if contentFilterResults.Violence != nil && contentFilterResults.Violence.Severity != nil {
			fmt.Fprintf(os.Stderr, "- Violence severity: %s\n", *contentFilterResults.Violence.Severity)
		}

		if contentFilterResults.Sexual != nil && contentFilterResults.Sexual.Severity != nil {
			fmt.Fprintf(os.Stderr, "- Sexual severity: %s\n", *contentFilterResults.Sexual.Severity)
		}

		if contentFilterResults.SelfHarm != nil && contentFilterResults.SelfHarm.Severity != nil {
			fmt.Fprintf(os.Stderr, "- Self-harm severity: %s\n", *contentFilterResults.SelfHarm.Severity)
		}
	}

	// Access the response content
	fmt.Fprintf(os.Stderr, "\nResponse: %s\n", resp.Choices[0].Message.Content)
}

// This example demonstrates how to apply Azure OpenAI prompt filtering with streaming responses
func Example_usingAzurePromptFilteringWithStreaming() {
	endpoint := os.Getenv("AOAI_ENDPOINT")
	model := os.Getenv("AOAI_MODEL")

	if endpoint == "" || model == "" {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not running example.")
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-08-01-preview"),
		azure.WithTokenCredential(tokenCredential),
	)

	// Example of streaming with Azure extensions
	fmt.Fprintf(os.Stderr, "Streaming example:\n")
	streamingParams := openai.ChatCompletionNewParams{
		Model:     openai.F(model),
		MaxTokens: openai.Int(256),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionMessageParam{
				Role:    openai.F(openai.ChatCompletionMessageParamRoleUser),
				Content: openai.F[any]("List 3 benefits of renewable energy"),
			},
		}),
	}

	stream := client.Chat.Completions.NewStreaming(
		context.TODO(),
		streamingParams,
	)

	var fullContent string

	for stream.Next() {
		chunk := stream.Current()

		// Get Azure-specific prompt filter results, if available
		azureChunk := azopenaiextensions.ChatCompletionChunk(chunk)
		promptFilterResults, err := azureChunk.PromptFilterResults()

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		if promptFilterResults != nil {
			fmt.Fprintf(os.Stderr, "- Prompt filter results detected\n")
		}

		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			fullContent += content
			fmt.Fprint(os.Stderr, content)
		}
	}

	if err := stream.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "\n\nStreaming complete. Full content length: %d characters\n", len(fullContent))
}
