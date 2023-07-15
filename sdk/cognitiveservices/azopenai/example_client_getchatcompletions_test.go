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
		panic(err)
	}

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, modelDeploymentID, nil)

	if err != nil {
		panic(err)
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
		Messages: messages,
	}, nil)

	if err != nil {
		panic(err)
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
		panic(err)
	}

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, modelDeploymentID, nil)

	if err != nil {
		panic(err)
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
		Messages: messages,
		N:        to.Ptr[int32](1),
	}, nil)

	if err != nil {
		panic(err)
	}

	streamReader := resp.ChatCompletionsStream
	gotReply := false

	for {
		chatCompletions, err := streamReader.Read()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			panic(err)
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
