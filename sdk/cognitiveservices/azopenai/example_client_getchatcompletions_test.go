// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
)

func ExampleClient_GetChatCompletions() {
	azureOpenAIKey := os.Getenv("AOAI_API_KEY")
	modelDeploymentID := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL_DEPLOYMENT")

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

	// This is a conversation in progress.
	// NOTE: all messages, regardless of role, count against token usage for this API.
	messages := []azopenai.ChatMessage{
		// You set the tone and rules of the conversation with a prompt as the system role.
		{Role: to.Ptr(azopenai.ChatRoleSystem), Content: to.Ptr("You are a helpful assistant. You will talk like a pirate.")},

		// The user asks a question
		{Role: to.Ptr(azopenai.ChatRoleUser), Content: to.Ptr("Can you help me?")},

		// The reply would come back from the ChatGPT. You'd add it to the conversation so we can maintain context.
		{Role: to.Ptr(azopenai.ChatRoleAssistant), Content: to.Ptr("Arrrr! Of course, me hearty! What can I do for ye?")},

		// The user answers the question based on the latest reply.
		{Role: to.Ptr(azopenai.ChatRoleUser), Content: to.Ptr("What's the best way to train a parrot?")},

		// from here you'd keep iterating, sending responses back from ChatGPT
	}

	gotReply := false

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		// This is a conversation in progress.
		// NOTE: all messages count against token usage for this API.
		Messages:     messages,
		DeploymentID: modelDeploymentID,
	}, nil)

	if err != nil {
		// TODO: handle error
	}

	for _, choice := range resp.Choices {
		gotReply = true
		fmt.Fprintf(os.Stderr, "Content[%d]: %s\n", *choice.Index, *choice.Message.Content)
	}

	if gotReply {
		fmt.Fprintf(os.Stderr, "Got chat completions reply\n")
	}

	// Output:
}

func ExampleClient_GetChatCompletions_functions() {
	azureOpenAIKey := os.Getenv("AOAI_API_KEY")
	modelDeploymentID := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL_DEPLOYMENT")

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

	resp, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		DeploymentID: modelDeploymentID,
		Messages: []azopenai.ChatMessage{
			{
				Role:    to.Ptr(azopenai.ChatRoleUser),
				Content: to.Ptr("What's the weather like in Boston, MA, in celsius?"),
			},
		},
		FunctionCall: &azopenai.ChatCompletionsOptionsFunctionCall{
			Value: to.Ptr("auto"),
		},
		Functions: []azopenai.FunctionDefinition{
			{
				Name:        to.Ptr("get_current_weather"),
				Description: to.Ptr("Get the current weather in a given location"),

				Parameters: map[string]any{
					"required": []string{"location"},
					"type":     "object",
					"properties": map[string]any{
						"location": map[string]any{
							"type":        "string",
							"description": "The city and state, e.g. San Francisco, CA",
						},
						"unit": map[string]any{
							"type": "string",
							"enum": []string{"celsius", "fahrenheit"},
						},
					},
				},
			},
		},
		Temperature: to.Ptr[float32](0.0),
	}, nil)

	if err != nil {
		// TODO: handle error
	}

	funcCall := resp.ChatCompletions.Choices[0].Message.FunctionCall

	// This is the function name we gave in the call to GetCompletions
	// Prints: Function name: "get_current_weather"
	fmt.Fprintf(os.Stderr, "Function name: %q\n", *funcCall.Name)

	// The arguments for your function come back as a JSON string
	var funcParams *struct {
		Location string `json:"location"`
		Unit     string `json:"unit"`
	}
	err = json.Unmarshal([]byte(*funcCall.Arguments), &funcParams)

	if err != nil {
		// TODO: handle error
	}

	// Prints:
	// Parameters: azopenai_test.location{Location:"Boston, MA", Unit:"celsius"}
	fmt.Fprintf(os.Stderr, "Parameters: %#v\n", *funcParams)

	// Output:
}

func ExampleClient_GetChatCompletionsStream() {
	azureOpenAIKey := os.Getenv("AOAI_API_KEY")
	modelDeploymentID := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL_DEPLOYMENT")

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

	// This is a conversation in progress.
	// NOTE: all messages, regardless of role, count against token usage for this API.
	messages := []azopenai.ChatMessage{
		// You set the tone and rules of the conversation with a prompt as the system role.
		{Role: to.Ptr(azopenai.ChatRoleSystem), Content: to.Ptr("You are a helpful assistant. You will talk like a pirate and limit your responses to 20 words or less.")},

		// The user asks a question
		{Role: to.Ptr(azopenai.ChatRoleUser), Content: to.Ptr("Can you help me?")},

		// The reply would come back from the ChatGPT. You'd add it to the conversation so we can maintain context.
		{Role: to.Ptr(azopenai.ChatRoleAssistant), Content: to.Ptr("Arrrr! Of course, me hearty! What can I do for ye?")},

		// The user answers the question based on the latest reply.
		{Role: to.Ptr(azopenai.ChatRoleUser), Content: to.Ptr("What's the best way to train a parrot?")},

		// from here you'd keep iterating, sending responses back from ChatGPT
	}

	resp, err := client.GetChatCompletionsStream(context.TODO(), azopenai.ChatCompletionsOptions{
		// This is a conversation in progress.
		// NOTE: all messages count against token usage for this API.
		Messages:     messages,
		N:            to.Ptr[int32](1),
		DeploymentID: modelDeploymentID,
	}, nil)

	if err != nil {
		// TODO: handle error
	}

	streamReader := resp.ChatCompletionsStream
	gotReply := false

	for {
		chatCompletions, err := streamReader.Read()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			// TODO: handle error
		}

		for _, choice := range chatCompletions.Choices {
			gotReply = true

			text := ""

			if choice.Delta.Content != nil {
				text = *choice.Delta.Content
			}

			role := ""

			if choice.Delta.Role != nil {
				role = string(*choice.Delta.Role)
			}

			fmt.Fprintf(os.Stderr, "Content[%d], role %q: %q\n", *choice.Index, role, text)
		}
	}

	if gotReply {
		fmt.Fprintf(os.Stderr, "Got chat completions streaming reply\n")
	}

	// Output:
}
