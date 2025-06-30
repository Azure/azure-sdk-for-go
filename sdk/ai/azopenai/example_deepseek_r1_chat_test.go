// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go"
)

// Example_deepseekReasoningBasic demonstrates basic chat completions using DeepSeek-R1 reasoning model.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Send a simple prompt to the DeepSeek-R1 reasoning model
// - Configure parameters for optimal reasoning performance
// - Process the response with step-by-step reasoning
//
// The example uses environment variables for configuration:
// - AOAI_DEEPSEEK_ENDPOINT: Your Azure OpenAI endpoint URL with DeepSeek model access
// - AOAI_DEEPSEEK_MODEL: The DeepSeek model deployment name (e.g., "deepseek-r1")
//
// DeepSeek-R1 is a reasoning model that provides detailed step-by-step analysis
// for complex problems, making it ideal for mathematical reasoning, logical deduction,
// and analytical problem solving.
func Example_deepseekReasoningBasic() {
	if !CheckRequiredEnvVars("AOAI_DEEPSEEK_ENDPOINT", "AOAI_DEEPSEEK_MODEL") {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not running example.\n")
		return
	}

	endpoint := os.Getenv("AOAI_DEEPSEEK_ENDPOINT")
	model := os.Getenv("AOAI_DEEPSEEK_MODEL")

	// Create a client with token credentials
	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	// Send a reasoning problem to DeepSeek-R1
	resp, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Model:       openai.ChatModel(model),
			MaxTokens:   openai.Int(1500),
			Temperature: openai.Float(0.1), // Lower temperature for more consistent reasoning
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfSystem: &openai.ChatCompletionSystemMessageParam{
						Content: openai.ChatCompletionSystemMessageParamContentUnion{
							OfString: openai.String("You are a helpful assistant that excels at step-by-step reasoning. Always show your thought process clearly and break down complex problems into manageable steps."),
						},
					},
				},
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String("A company has 100 employees. If 60% work in engineering, 25% work in sales, and the rest work in administration, how many people work in each department? Please show your reasoning step by step."),
						},
					},
				},
			},
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Fprintf(os.Stderr, "DeepSeek-R1 Reasoning Response:\n")
		fmt.Fprintf(os.Stderr, "%s\n", resp.Choices[0].Message.Content)

		choice := resp.Choices[0]

		// Show the internal reasoning process (DeepSeek-R1's thinking)
		if choice.Message.JSON.ExtraFields != nil {
			if reasoningField, ok := choice.Message.JSON.ExtraFields["reasoning_content"]; ok {
				reasoningContent := reasoningField.Raw()
				if reasoningContent != "" {
					fmt.Fprintf(os.Stderr, "=== DeepSeek-R1 Internal Reasoning Process ===\n")
					fmt.Fprintf(os.Stderr, "%s\n", reasoningContent)
					fmt.Fprintf(os.Stderr, "\n")
				}
			}
		}
	}

	fmt.Fprintf(os.Stderr, "\n=== Basic Reasoning Example Complete ===\n")
}

// Example_deepseekReasoningMultiTurn demonstrates multi-turn conversations with DeepSeek-R1.
// This example shows how to:
// - Maintain conversation context across multiple turns
// - Build upon previous reasoning steps
// - Ask follow-up questions that reference earlier parts of the conversation
// - Handle complex problem-solving scenarios that require multiple interactions
// - Manage conversation history in a chat application
//
// When using the model for a chat application, you'll need to manage the history
// of that conversation and send the latest messages to the model.
//
// The example uses environment variables for configuration:
// - AOAI_DEEPSEEK_ENDPOINT: Your Azure OpenAI endpoint URL with DeepSeek model access
// - AOAI_DEEPSEEK_MODEL: The DeepSeek model deployment name (e.g., "deepseek-r1")
func Example_deepseekReasoningMultiTurn() {
	if !CheckRequiredEnvVars("AOAI_DEEPSEEK_ENDPOINT", "AOAI_DEEPSEEK_MODEL") {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not running example.\n")
		return
	}

	endpoint := os.Getenv("AOAI_DEEPSEEK_ENDPOINT")
	model := os.Getenv("AOAI_DEEPSEEK_MODEL")

	// Create a client with token credentials
	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "=== Multi-Turn Reasoning Conversation ===\n\n")

	// Build conversation history with multiple messages
	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: openai.String("You are a helpful assistant."),
				},
			},
		},
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("I am going to Paris, what should I see?"),
				},
			},
		},
		{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				Content: openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: openai.String("Paris, the capital of France, is known for its stunning architecture, art museums, historical landmarks, and romantic atmosphere. Here are some of the top attractions to see in Paris:\n \n 1. The Eiffel Tower: The iconic Eiffel Tower is one of the most recognizable landmarks in the world and offers breathtaking views of the city.\n 2. The Louvre Museum: The Louvre is one of the world's largest and most famous museums, housing an impressive collection of art and artifacts, including the Mona Lisa.\n 3. Notre-Dame Cathedral: This beautiful cathedral is one of the most famous landmarks in Paris and is known for its Gothic architecture and stunning stained glass windows.\n \n These are just a few of the many attractions that Paris has to offer. With so much to see and do, it's no wonder that Paris is one of the most popular tourist destinations in the world."),
				},
			},
		},
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("What is so great about #1?"),
				},
			},
		},
	}

	// Send the multi-turn conversation
	resp, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Model:     openai.ChatModel(model),
			MaxTokens: openai.Int(2048),
			Messages:  messages,
		},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Fprintf(os.Stderr, "DeepSeek-R1: %s\n", resp.Choices[0].Message.Content)
	}

	fmt.Fprintf(os.Stderr, "\n=== Multi-Turn Conversation Complete ===\n")
}

// Example_deepseekReasoningStreaming demonstrates streaming responses with DeepSeek-R1.
// This example shows how to:
// - Create a streaming chat completion request
// - Process streaming responses as they arrive
// - Handle the reasoning process in real-time
// - Provide a better user experience with immediate feedback
//
// The example uses environment variables for configuration:
// - AOAI_DEEPSEEK_ENDPOINT: Your Azure OpenAI endpoint URL with DeepSeek model access
// - AOAI_DEEPSEEK_MODEL: The DeepSeek model deployment name (e.g., "deepseek-r1")
//
// This example uses a simple math problem to demonstrate DeepSeek-R1's step-by-step
// reasoning capabilities in a streaming context.
func Example_deepseekReasoningStreaming() {
	if !CheckRequiredEnvVars("AOAI_DEEPSEEK_ENDPOINT", "AOAI_DEEPSEEK_MODEL") {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not running example.\n")
		return
	}

	endpoint := os.Getenv("AOAI_DEEPSEEK_ENDPOINT")
	model := os.Getenv("AOAI_DEEPSEEK_MODEL")

	// Create a client with token credentials
	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}
	// Create a streaming chat completion
	stream := client.Chat.Completions.NewStreaming(
		context.TODO(), openai.ChatCompletionNewParams{
			Model:       openai.ChatModel(model),
			MaxTokens:   openai.Int(1500),  // Reduced for simpler problem
			Temperature: openai.Float(0.1), // Lower temperature for consistent reasoning
			Messages: []openai.ChatCompletionMessageParamUnion{{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: openai.String("You are a helpful assistant that excels at step-by-step reasoning. Always show your thought process clearly."),
					},
				},
			},
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String("If I have 24 apples and I want to divide them equally among 6 friends, how many apples will each friend get? Also, if I buy 3 more bags of apples and each bag contains 8 apples, how many total apples will I have? Please show your reasoning step by step."),
						},
					},
				},
			},
		},
	)

	for stream.Next() {
		evt := stream.Current()
		if len(evt.Choices) > 0 {
			choice := evt.Choices[0]

			// Output content
			if choice.Delta.Content != "" {
				fmt.Fprintf(os.Stderr, "%s", choice.Delta.Content)
			}

			// Output reasoning content if present
			if choice.Delta.JSON.ExtraFields != nil {
				if reasoningField, ok := choice.Delta.JSON.ExtraFields["reasoning_content"]; ok {
					reasoningText := reasoningField.Raw()
					// Format reasoning content properly
					if reasoningText != "" && reasoningText != " " {
						// Clean up basic formatting issues
						cleanedContent := strings.ReplaceAll(reasoningText, `"`, "")
						cleanedContent = strings.ReplaceAll(cleanedContent, "null", "")
						fmt.Fprintf(os.Stderr, "%s", cleanedContent)
					}
				}
			}
		}
	}

	if stream.Err() != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: %s\n", stream.Err())
		return
	}

	fmt.Fprintf(os.Stderr, "\n\n=== Streaming Example Complete ===\n")
}
