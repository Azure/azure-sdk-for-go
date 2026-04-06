// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/option"
)

// Example_getChatCompletions demonstrates how to use Azure OpenAI's Chat Completions API.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Structure a multi-turn conversation with different message roles
// - Send a chat completion request and handle the response
// - Process multiple response choices and finish reasons
//
// The example uses environment variables for configuration:
// - AOAI_CHAT_COMPLETIONS_MODEL: The deployment name of your chat model
// - AOAI_CHAT_COMPLETIONS_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
//
// Chat completions are useful for:
// - Building conversational AI interfaces
// - Creating chatbots with personality
// - Maintaining context across multiple interactions
// - Generating human-like text responses
func Example_getChatCompletions() {
	model := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_ENDPOINT")

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// This is a conversation in progress.
	// NOTE: all messages, regardless of role, count against token usage for this API.
	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			// You set the tone and rules of the conversation with a prompt as the system role.
			{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: openai.String("You are a helpful assistant. You will talk like a pirate."),
					},
				},
			},
			// The user asks a question
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("Can you help me?"),
					},
				},
			},
			// The reply would come back from the ChatGPT. You'd add it to the conversation so we can maintain context.
			{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String("Arrrr! Of course, me hearty! What can I do for ye?"),
					},
				},
			},
			// The user answers the question based on the latest reply.
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("What's the best way to train a parrot?"),
					},
				},
			},
		},
	})

	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	gotReply := false

	for _, choice := range resp.Choices {
		gotReply = true

		if choice.Message.Content != "" {
			fmt.Fprintf(os.Stderr, "Content[%d]: %s\n", choice.Index, choice.Message.Content)
		}

		if choice.FinishReason != "" {
			fmt.Fprintf(os.Stderr, "Finish reason[%d]: %s\n", choice.Index, choice.FinishReason)
		}
	}

	if gotReply {
		fmt.Fprintf(os.Stderr, "Got chat completions reply\n")
	}
}

// Example_chatCompletionsFunctions demonstrates how to use Azure OpenAI's function calling feature.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Define a function schema for weather information
// - Request function execution through the chat API
// - Parse and handle function call responses
//
// The example uses environment variables for configuration:
// - AOAI_CHAT_COMPLETIONS_MODEL: The deployment name of your chat model
// - AOAI_CHAT_COMPLETIONS_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
//
// Tool calling is useful for:
// - Integrating external APIs and services
// - Structured data extraction from natural language
// - Task automation and workflow integration
// - Building context-aware applications
func Example_getChatCompletions_usingTools() {
	model := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_ENDPOINT")

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// Define the function schema
	functionSchema := map[string]interface{}{
		"required": []string{"location"},
		"type":     "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "The city and state, e.g. San Francisco, CA",
			},
			"unit": map[string]interface{}{
				"type": "string",
				"enum": []string{"celsius", "fahrenheit"},
			},
		},
	}

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("What's the weather like in Boston, MA, in celsius?"),
					},
				},
			},
		},
		Tools: []openai.ChatCompletionToolUnionParam{
			{
				OfFunction: &openai.ChatCompletionFunctionToolParam{
					Function: openai.FunctionDefinitionParam{
						Name:        "get_current_weather",
						Description: openai.String("Get the current weather in a given location"),
						Parameters:  functionSchema,
					},
				},
			},
		},
		Temperature: openai.Float(0.0),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		toolCall := resp.Choices[0].Message.ToolCalls[0]

		// This is the function name we gave in the call
		fmt.Fprintf(os.Stderr, "Function name: %q\n", toolCall.Function.Name)

		// The arguments for your function come back as a JSON string
		var funcParams struct {
			Location string `json:"location"`
			Unit     string `json:"unit"`
		}

		err = json.Unmarshal([]byte(toolCall.Function.Arguments), &funcParams)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		fmt.Fprintf(os.Stderr, "Parameters: %#v\n", funcParams)
	}

}

// Example_chatCompletionsLegacyFunctions demonstrates using the legacy function calling format.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Define a function schema using the legacy format
// - Use tools API for backward compatibility
// - Handle function calling responses
//
// The example uses environment variables for configuration:
// - AOAI_CHAT_COMPLETIONS_MODEL_LEGACY_FUNCTIONS_MODEL: The deployment name of your chat model
// - AOAI_CHAT_COMPLETIONS_MODEL_LEGACY_FUNCTIONS_ENDPOINT: Your Azure OpenAI endpoint URL
// - AZURE_OPENAI_API_VERSION: Azure OpenAI service API version to use. See https://learn.microsoft.com/azure/ai-foundry/openai/api-version-lifecycle?tabs=go for information about API versions.
//
// Legacy function support ensures:
// - Compatibility with older implementations
// - Smooth transition to new tools API
// - Support for existing function-based workflows
func Example_chatCompletionsLegacyFunctions() {
	model := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL_LEGACY_FUNCTIONS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL_LEGACY_FUNCTIONS_ENDPOINT")
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

	// Define the function schema
	parametersJSON := map[string]interface{}{
		"required": []string{"location"},
		"type":     "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "The city and state, e.g. San Francisco, CA",
			},
			"unit": map[string]interface{}{
				"type": "string",
				"enum": []string{"celsius", "fahrenheit"},
			},
		},
	}

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("What's the weather like in Boston, MA, in celsius?"),
					},
				},
			},
		},
		// Note: Legacy functions are supported through the Tools API in the OpenAI Go SDK
		Tools: []openai.ChatCompletionToolUnionParam{
			{
				OfFunction: &openai.ChatCompletionFunctionToolParam{
					Function: openai.FunctionDefinitionParam{
						Name:        "get_current_weather",
						Description: openai.String("Get the current weather in a given location"),
						Parameters:  parametersJSON,
					},
				},
			},
		},
		ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: openai.String("auto"),
		},
		Temperature: openai.Float(0.0),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		toolCall := resp.Choices[0].Message.ToolCalls[0]

		// This is the function name we gave in the call
		fmt.Fprintf(os.Stderr, "Function name: %q\n", toolCall.Function.Name)

		// The arguments for your function come back as a JSON string
		var funcParams struct {
			Location string `json:"location"`
			Unit     string `json:"unit"`
		}

		err = json.Unmarshal([]byte(toolCall.Function.Arguments), &funcParams)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		fmt.Fprintf(os.Stderr, "Parameters: %#v\n", funcParams)
	}

}

// Example_chatCompletionStream demonstrates streaming responses from the Chat Completions API.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Set up a streaming chat completion request
// - Process incremental response chunks
// - Handle streaming errors and completion
//
// The example uses environment variables for configuration:
// - AOAI_CHAT_COMPLETIONS_MODEL: The deployment name of your chat model
// - AOAI_CHAT_COMPLETIONS_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
//
// Streaming is useful for:
// - Real-time response display
// - Improved perceived latency
// - Interactive chat interfaces
// - Long-form content generation
func Example_chatCompletionStream() {
	model := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_ENDPOINT")

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// This is a conversation in progress
	stream := client.Chat.Completions.NewStreaming(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			// System message sets the tone
			{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: openai.String("You are a helpful assistant. You will talk like a pirate and limit your responses to 20 words or less."),
					},
				},
			},
			// User question
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("Can you help me?"),
					},
				},
			},
			// Assistant reply
			{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String("Arrrr! Of course, me hearty! What can I do for ye?"),
					},
				},
			},
			// User follow-up
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("What's the best way to train a parrot?"),
					},
				},
			},
		},
	})

	gotReply := false

	for stream.Next() {
		gotReply = true
		evt := stream.Current()
		if len(evt.Choices) > 0 {
			print(evt.Choices[0].Delta.Content)
		}
	}

	if stream.Err() != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}

	if gotReply {
		fmt.Fprintf(os.Stderr, "\nGot chat completions streaming reply\n")
	}

}

// Example_chatCompletionsStructuredOutputs demonstrates using structured outputs with function calling.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Define complex JSON schemas for structured output
// - Request specific data structures through function calls
// - Parse and validate structured responses
//
// The example uses environment variables for configuration:
// - AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_MODEL: The deployment name of your chat model
// - AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
//
// Structured outputs are useful for:
// - Database query generation
// - Data extraction and transformation
// - API request formatting
// - Consistent response formatting
func Example_chatCompletionsStructuredOutputs() {
	model := os.Getenv("AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_ENDPOINT")

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// Define the structured output schema
	structuredJSONSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"table_name": map[string]interface{}{
				"type": "string",
				"enum": []string{"orders"},
			},
			"columns": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
					"enum": []string{
						"id", "status", "expected_delivery_date", "delivered_at",
						"shipped_at", "ordered_at", "canceled_at",
					},
				},
			},
			"conditions": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"column": map[string]interface{}{
							"type": "string",
						},
						"operator": map[string]interface{}{
							"type": "string",
							"enum": []string{"=", ">", "<", ">=", "<=", "!="},
						},
						"value": map[string]interface{}{
							"anyOf": []map[string]interface{}{
								{"type": "string"},
								{"type": "number"},
								{
									"type": "object",
									"properties": map[string]interface{}{
										"column_name": map[string]interface{}{"type": "string"},
									},
									"required":             []string{"column_name"},
									"additionalProperties": false,
								},
							},
						},
					},
					"required":             []string{"column", "operator", "value"},
					"additionalProperties": false,
				},
			},
			"order_by": map[string]interface{}{
				"type": "string",
				"enum": []string{"asc", "desc"},
			},
		},
		"required":             []string{"table_name", "columns", "conditions", "order_by"},
		"additionalProperties": false,
	}

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String("You are a helpful assistant. The current date is August 6, 2024. You help users query for the data they are looking for by calling the query function."),
					},
				},
			},
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("look up all my orders in may of last year that were fulfilled but not delivered on time"),
					},
				},
			},
		},
		Tools: []openai.ChatCompletionToolUnionParam{
			{
				OfFunction: &openai.ChatCompletionFunctionToolParam{
					Function: openai.FunctionDefinitionParam{
						Name:       "query",
						Parameters: structuredJSONSchema,
					},
				},
			},
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	if len(resp.Choices) > 0 && len(resp.Choices[0].Message.ToolCalls) > 0 {
		fn := resp.Choices[0].Message.ToolCalls[0].Function

		argumentsObj := map[string]interface{}{}
		err = json.Unmarshal([]byte(fn.Arguments), &argumentsObj)

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		fmt.Fprintf(os.Stderr, "%#v\n", argumentsObj)
	}

}

// Example_structuredOutputsResponseFormat demonstrates using JSON response formatting.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Define JSON schema for response formatting
// - Request structured mathematical solutions
// - Parse and process formatted JSON responses
//
// The example uses environment variables for configuration:
// - AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_MODEL: The deployment name of your model
// - AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
//
// Response formatting is useful for:
// - Mathematical problem solving
// - Step-by-step explanations
// - Structured data generation
// - Consistent output formatting
func Example_structuredOutputsWithTools() {
	model := os.Getenv("AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_ENDPOINT")

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// Define the structured output schema
	mathResponseSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"steps": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"explanation": map[string]interface{}{"type": "string"},
						"output":      map[string]interface{}{"type": "string"},
					},
					"required":             []string{"explanation", "output"},
					"additionalProperties": false,
				},
			},
			"final_answer": map[string]interface{}{"type": "string"},
		},
		"required":             []string{"steps", "final_answer"},
		"additionalProperties": false,
	}

	resp, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String("You are a helpful math tutor."),
					},
				},
			},
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("solve 8x + 31 = 2"),
					},
				},
			},
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   "math_response",
					Schema: mathResponseSchema,
				},
			},
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		responseObj := map[string]interface{}{}
		err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &responseObj)

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		fmt.Fprintf(os.Stderr, "%#v", responseObj)
	}

}
