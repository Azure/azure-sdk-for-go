//go:build go1.21
// +build go1.21

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/require"
)

func TestChatCompletions_extensions_bringYourOwnData(t *testing.T) {
	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.ChatCompletionsOYD.Endpoint)

	inputParams := openai.ChatCompletionNewParams{
		Model:     openai.ChatModel(azureOpenAI.ChatCompletionsOYD.Model),
		MaxTokens: openai.Int(512),
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("What does the OpenAI package do?"),
					},
				},
			},
		},
	}

	resp, err := client.Chat.Completions.New(context.Background(), inputParams,
		azopenai.WithDataSources(&azureOpenAI.Cognitive))
	customRequireNoError(t, err)
	require.NotEmpty(t, resp)

	msg := azopenai.ChatCompletionMessage(resp.Choices[0].Message)

	msgContext, err := msg.Context()
	require.NoError(t, err)
	require.NotEmpty(t, msgContext.Citations[0].Content)

	require.NotEmpty(t, msg.Content)
	require.Equal(t, "stop", resp.Choices[0].FinishReason)

	t.Logf("Content = %s", resp.Choices[0].Message.Content)
}

func TestChatExtensionsStreaming_extensions_bringYourOwnData(t *testing.T) {
	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.ChatCompletionsOYD.Endpoint)

	inputParams := openai.ChatCompletionNewParams{
		Model:     openai.ChatModel(azureOpenAI.ChatCompletionsOYD.Model),
		MaxTokens: openai.Int(512),
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("What does the OpenAI package do?"),
				},
			},
		}},
	}

	streamer := client.Chat.Completions.NewStreaming(context.Background(), inputParams,
		azopenai.WithDataSources(
			&azureOpenAI.Cognitive,
		))

	t.Cleanup(func() {
		err := streamer.Close()
		require.NoError(t, err)
	})

	text := ""

	first := true

	for streamer.Next() {
		chunk := streamer.Current()

		if first {
			// when you BYOD you get some extra content showing you metadata/info from the external
			// data source.
			first = false

			msgContext, err := azopenai.ChatCompletionChunkChoiceDelta(chunk.Choices[0].Delta).Context()
			require.NoError(t, err)
			require.NotEmpty(t, msgContext.Citations[0].Content)
		}

		for _, choice := range chunk.Choices {
			text += choice.Delta.Content
		}
	}

	customRequireNoError(t, streamer.Err())
	require.NotEmpty(t, text)

	t.Logf("Streaming content = %s", text)
}
