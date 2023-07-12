//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var chatCompletionsRequest = azopenai.ChatCompletionsOptions{
	Messages: []azopenai.ChatMessage{
		{
			Role:    to.Ptr(azopenai.ChatRole("user")),
			Content: to.Ptr("Count to 10, with a comma between each number, no newlines and a period at the end. E.g., 1, 2, 3, ..."),
		},
	},
	MaxTokens:   to.Ptr(int32(1024)),
	Temperature: to.Ptr(float32(0.0)),
	Model:       &openAIChatCompletionsModel,
}

var expectedContent = "1, 2, 3, 4, 5, 6, 7, 8, 9, 10."
var expectedRole = azopenai.ChatRoleAssistant

func TestClient_GetChatCompletions(t *testing.T) {
	cred, err := azopenai.NewKeyCredential(apiKey)
	require.NoError(t, err)

	chatClient, err := azopenai.NewClientWithKeyCredential(endpoint, cred, chatCompletionsModelDeployment, newClientOptionsForTest(t))
	require.NoError(t, err)

	testGetChatCompletions(t, chatClient)
}

func TestClient_GetChatCompletionsStream(t *testing.T) {
	cred, err := azopenai.NewKeyCredential(apiKey)
	require.NoError(t, err)

	chatClient, err := azopenai.NewClientWithKeyCredential(endpoint, cred, chatCompletionsModelDeployment, newClientOptionsForTest(t))
	require.NoError(t, err)

	testGetChatCompletionsStream(t, chatClient, true)
}

func TestClient_OpenAI_GetChatCompletions(t *testing.T) {
	chatClient := newOpenAIClientForTest(t)
	testGetChatCompletions(t, chatClient)
}

func TestClient_OpenAI_GetChatCompletionsStream(t *testing.T) {
	chatClient := newOpenAIClientForTest(t)
	testGetChatCompletionsStream(t, chatClient, false)
}

func testGetChatCompletions(t *testing.T, client *azopenai.Client) {
	expected := azopenai.ChatCompletions{
		Choices: []azopenai.ChatChoice{
			{
				Message: &azopenai.ChatChoiceMessage{
					Role:    &expectedRole,
					Content: &expectedContent,
				},
				Index:        to.Ptr(int32(0)),
				FinishReason: to.Ptr(azopenai.CompletionsFinishReason("stop")),
			},
		},
		Usage: &azopenai.CompletionsUsage{
			// these change depending on which model you use. These #'s work for gpt-4, which is
			// what I'm using for these tests.
			CompletionTokens: to.Ptr(int32(29)),
			PromptTokens:     to.Ptr(int32(42)),
			TotalTokens:      to.Ptr(int32(71)),
		},
	}

	resp, err := client.GetChatCompletions(context.Background(), chatCompletionsRequest, nil)
	require.NoError(t, err)

	require.NotEmpty(t, resp.ID)
	require.NotEmpty(t, resp.Created)

	expected.ID = resp.ID
	expected.Created = resp.Created

	require.Equal(t, expected, resp.ChatCompletions)
}

func testGetChatCompletionsStream(t *testing.T, client *azopenai.Client, isAzure bool) {
	streamResp, err := client.GetChatCompletionsStream(context.Background(), chatCompletionsRequest, nil)
	require.NoError(t, err)

	if isAzure {
		// there's a bug right now where the first event comes back empty
		// Issue: https://github.com/Azure/azure-sdk-for-go/issues/21086
		_, err := streamResp.ChatCompletionsStream.Read()
		require.NoError(t, err)
	}

	// the data comes back differently for streaming
	// 1. the text comes back in the ChatCompletion.Delta field
	// 2. the role is only sent on the first streamed ChatCompletion
	// check that the role came back as well.
	var choices []azopenai.ChatChoice

	for {
		completion, err := streamResp.ChatCompletionsStream.Read()

		if errors.Is(err, io.EOF) {
			break
		}

		require.NoError(t, err)
		require.Equal(t, 1, len(completion.Choices))
		choices = append(choices, completion.Choices[0])
	}

	var message string

	for _, choice := range choices {
		if choice.Delta.Content == nil {
			continue
		}

		message += *choice.Delta.Content
	}

	require.Equal(t, expectedContent, message, "Ultimately, the same result as GetChatCompletions(), just sent across the .Delta field instead")

	require.Equal(t, azopenai.ChatRoleAssistant, expectedRole)
}

func TestClient_GetChatCompletions_DefaultAzureCredential(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skipf("Not running this test in playback (for now)")
	}

	if os.Getenv("USE_TOKEN_CREDS") != "true" {
		t.Skipf("USE_TOKEN_CREDS is not true, disabling token credential tests")
	}

	recordingTransporter := newRecordingTransporter(t)

	dac, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: policy.ClientOptions{
			Transport: recordingTransporter,
		},
	})
	require.NoError(t, err)

	chatClient, err := azopenai.NewClient(endpoint, dac, chatCompletionsModelDeployment, &azopenai.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: recordingTransporter},
	})
	require.NoError(t, err)

	testGetChatCompletions(t, chatClient)
}

func TestClient_GetChatCompletions_InvalidModel(t *testing.T) {
	cred, err := azopenai.NewKeyCredential(apiKey)
	require.NoError(t, err)

	chatClient, err := azopenai.NewClientWithKeyCredential(endpoint, cred, "thisdoesntexist", newClientOptionsForTest(t))
	require.NoError(t, err)

	_, err = chatClient.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatMessage{
			{
				Role:    to.Ptr(azopenai.ChatRole("user")),
				Content: to.Ptr("Count to 100, with a comma between each number and no newlines. E.g., 1, 2, 3, ..."),
			},
		},
		MaxTokens:   to.Ptr(int32(1024)),
		Temperature: to.Ptr(float32(0.0)),
	}, nil)

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, "DeploymentNotFound", respErr.ErrorCode)
}

func TestClient_GetChatCompletionsStream_Error(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip()
	}

	doTest := func(t *testing.T, client *azopenai.Client) {
		t.Helper()
		streamResp, err := client.GetChatCompletionsStream(context.Background(), chatCompletionsRequest, nil)
		require.Empty(t, streamResp)
		assertResponseIsError(t, err)
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		client := newBogusAzureOpenAIClient(t, chatCompletionsModelDeployment)
		doTest(t, client)
	})

	t.Run("OpenAI", func(t *testing.T) {
		client := newBogusOpenAIClient(t)
		doTest(t, client)
	})
}
