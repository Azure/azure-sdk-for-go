// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

func Example_usingReasoningEffort() {
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_ENDPOINT")
	model := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL") // This should be a model that supports reasoning like "o1"

	if endpoint == "" || model == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-12-01-preview"),
		azure.WithTokenCredential(tokenCredential),
	)

	// Create a complex planning request that can benefit from different reasoning levels
	prompt := "I have $50,000 to invest and I want to create a diversified portfolio. Suggest a mix of stocks, bonds, and other investment options, including the expected returns and risks."

	// First, use low reasoning effort
	fmt.Fprintf(os.Stderr, "Querying with LOW reasoning effort...\n")
	lowReasoningResponse, err := getCompletionWithReasoning(client, model, prompt, "low")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR with low reasoning: %s\n", err)
		return
	}

	// Then, use high reasoning effort
	fmt.Fprintf(os.Stderr, "Querying with HIGH reasoning effort...\n")
	highReasoningResponse, err := getCompletionWithReasoning(client, model, prompt, "high")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR with high reasoning: %s\n", err)
		return
	}

	// Compare the results
	fmt.Fprintf(os.Stderr, "\n=== Comparison of Token Usage ===\n")
	fmt.Fprintf(os.Stderr, "Low reasoning:\n")
	fmt.Fprintf(os.Stderr, "  Completion tokens: %v\n", lowReasoningResponse.Usage.CompletionTokens)
	fmt.Fprintf(os.Stderr, "  Reasoning tokens: %v\n", lowReasoningResponse.Usage.CompletionTokensDetails.ReasoningTokens)
	fmt.Fprintf(os.Stderr, "  Prompt tokens: %v\n", lowReasoningResponse.Usage.PromptTokens)
	fmt.Fprintf(os.Stderr, "  Total tokens: %v\n", lowReasoningResponse.Usage.TotalTokens)

	fmt.Fprintf(os.Stderr, "\nHigh reasoning:\n")
	fmt.Fprintf(os.Stderr, "  Completion tokens: %v\n", highReasoningResponse.Usage.CompletionTokens)
	fmt.Fprintf(os.Stderr, "  Reasoning tokens: %v\n", highReasoningResponse.Usage.CompletionTokensDetails.ReasoningTokens)
	fmt.Fprintf(os.Stderr, "  Prompt tokens: %v\n", highReasoningResponse.Usage.PromptTokens)
	fmt.Fprintf(os.Stderr, "  Total tokens: %v\n", highReasoningResponse.Usage.TotalTokens)

	// Show the first part of each response for comparison
	fmt.Fprintf(os.Stderr, "\n=== Sample Response Comparison ===\n")

	lowContent := lowReasoningResponse.Choices[0].Message.Content
	highContent := highReasoningResponse.Choices[0].Message.Content

	fmt.Fprintf(os.Stderr, "Low reasoning (first 1000 chars):\n%s...\n", truncateString(lowContent, 1000))
	fmt.Fprintf(os.Stderr, "\nHigh reasoning (first 1000 chars):\n%s...\n", truncateString(highContent, 1000))

	fmt.Fprintf(os.Stderr, "\nExample complete\n")

	// Output:
}

// Helper function to get completions with specified reasoning effort
func getCompletionWithReasoning(client *openai.Client, model string, prompt string, reasoningEffort string) (*openai.ChatCompletion, error) {
	var reasoningEffortParam openai.ChatCompletionReasoningEffort
	if reasoningEffort == "low" {
		reasoningEffortParam = openai.ChatCompletionReasoningEffortLow
	} else if reasoningEffort == "high" {
		reasoningEffortParam = openai.ChatCompletionReasoningEffortHigh
	} else {
		return nil, fmt.Errorf("invalid reasoning effort: %s", reasoningEffort)
	}
	chatParams := openai.ChatCompletionNewParams{
		Model:           openai.F(model),
		ReasoningEffort: openai.F(reasoningEffortParam),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionMessageParam{
				Role:    openai.F(openai.ChatCompletionMessageParamRoleSystem),
				Content: openai.F[any]("You are a financial advisor helping clients with investment strategies."),
			},
			openai.ChatCompletionMessageParam{
				Role:    openai.F(openai.ChatCompletionMessageParamRoleUser),
				Content: openai.F[any](prompt),
			},
		}),
	}

	resp, err := client.Chat.Completions.New(
		context.TODO(),
		chatParams)

	return resp, err
}

// Helper function to truncate a string to a specified length
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}
