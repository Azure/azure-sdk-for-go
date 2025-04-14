// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
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
