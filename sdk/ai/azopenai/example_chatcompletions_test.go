// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go"
)

// Example_getChatCompletions demonstrates how to use the Chat Completions API.
func Example_getChatCompletions() {
	if !CheckRequiredEnvVars("AOAI_CHAT_COMPLETIONS_MODEL", "AOAI_CHAT_COMPLETIONS_ENDPOINT") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	model := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_ENDPOINT")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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

// Example_chatCompletionsFunctions demonstrates how to use the Chat Completions API with function calling.
func Example_chatCompletionsFunctions() {
	if !CheckRequiredEnvVars("AOAI_CHAT_COMPLETIONS_MODEL", "AOAI_CHAT_COMPLETIONS_ENDPOINT") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	model := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_ENDPOINT")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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
		Tools: []openai.ChatCompletionToolParam{
			{
				Function: openai.FunctionDefinitionParam{
					Name:        "get_current_weather",
					Description: openai.String("Get the current weather in a given location"),
					Parameters:  functionSchema,
				},
				Type: "function",
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

// Example_chatCompletionsLegacyFunctions demonstrates using legacy-style function calling with the Chat Completions API.
func Example_chatCompletionsLegacyFunctions() {
	if !CheckRequiredEnvVars("AOAI_CHAT_COMPLETIONS_MODEL_LEGACY_FUNCTIONS_MODEL", "AOAI_CHAT_COMPLETIONS_MODEL_LEGACY_FUNCTIONS_ENDPOINT") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	model := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL_LEGACY_FUNCTIONS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL_LEGACY_FUNCTIONS_ENDPOINT")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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
		Tools: []openai.ChatCompletionToolParam{
			{
				Type: "function",
				Function: openai.FunctionDefinitionParam{
					Name:        "get_current_weather",
					Description: openai.String("Get the current weather in a given location"),
					Parameters:  parametersJSON,
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

// Example_chatCompletionStream demonstrates how to use the Chat Completions API with streaming responses.
func Example_chatCompletionStream() {
	if !CheckRequiredEnvVars("AOAI_CHAT_COMPLETIONS_MODEL", "AOAI_CHAT_COMPLETIONS_ENDPOINT") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	model := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_ENDPOINT")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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

// Example_chatCompletionsStructuredOutputs demonstrates how to use the Chat Completions API with structured outputs using function calling.
func Example_chatCompletionsStructuredOutputs() {
	if !CheckRequiredEnvVars("AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_MODEL", "AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_ENDPOINT") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	model := os.Getenv("AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_ENDPOINT")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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
		Tools: []openai.ChatCompletionToolParam{
			{
				Type: "function",
				Function: openai.FunctionDefinitionParam{
					Name:       "query",
					Parameters: structuredJSONSchema,
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

// Example_structuredOutputsResponseFormat demonstrates how to use the Chat Completions API
// with structured outputs using response format.
func Example_structuredOutputsResponseFormat() {
	if !CheckRequiredEnvVars("AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_MODEL", "AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_ENDPOINT") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	model := os.Getenv("AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_MODEL")
	endpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_STRUCTURED_OUTPUTS_ENDPOINT")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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
