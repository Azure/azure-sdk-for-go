//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestChatCompletions_extensions_bringYourOwnData(t *testing.T) {
	t.Skip("TEMP DISABLED: search index authentication has been disabled temporarily: https://github.com/Azure/azure-sdk-for-go/issues/22966")

	client := newTestClient(t, azureOpenAI.ChatCompletionsOYD.Endpoint)

	resp, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("What does PR complete mean?")},
		},
		MaxTokens: to.Ptr[int32](512),
		AzureExtensionsOptions: []azopenai.AzureChatExtensionConfigurationClassification{
			&azureOpenAI.Cognitive,
		},
		DeploymentName: &azureOpenAI.ChatCompletionsOYD.Model,
	}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, resp)

	msgContext := resp.Choices[0].Message.Context
	require.NotEmpty(t, msgContext.Citations[0].Content)

	require.NotEmpty(t, *resp.Choices[0].Message.Content)
	require.Equal(t, azopenai.CompletionsFinishReasonStopped, *resp.Choices[0].FinishReason)
}

func TestChatExtensionsStreaming_extensions_bringYourOwnData(t *testing.T) {
	t.Skip("TEMP DISABLED: search index authentication has been disabled temporarily: https://github.com/Azure/azure-sdk-for-go/issues/22966")

	client := newTestClient(t, azureOpenAI.ChatCompletionsOYD.Endpoint)

	streamResp, err := client.GetChatCompletionsStream(context.Background(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("What does PR complete mean?")},
		},
		MaxTokens: to.Ptr[int32](512),
		AzureExtensionsOptions: []azopenai.AzureChatExtensionConfigurationClassification{
			&azureOpenAI.Cognitive,
		},
		DeploymentName: &azureOpenAI.ChatCompletionsOYD.Model,
	}, nil)

	require.NoError(t, err)
	defer streamResp.ChatCompletionsStream.Close()

	text := ""

	first := false

	for {
		event, err := streamResp.ChatCompletionsStream.Read()

		if errors.Is(err, io.EOF) {
			break
		}

		require.NoError(t, err)

		if first {
			// when you BYOD you get some extra content showing you metadata/info from the external
			// data source.
			first = false
			msgContext := event.Choices[0].Message.Context
			require.NotEmpty(t, msgContext.Citations[0].Content)
		}

		for _, choice := range event.Choices {
			if choice.Delta != nil && choice.Delta.Content != nil {
				text += *choice.Delta.Content
			}
		}
	}

	require.NotEmpty(t, text)
}
