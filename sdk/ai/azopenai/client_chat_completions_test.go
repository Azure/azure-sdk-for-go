//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func newTestChatCompletionOptions(deployment string) azopenai.ChatCompletionsOptions {
	return azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestUserMessage{
				Content: azopenai.NewChatRequestUserMessageContent("Count to 10, with a comma between each number, no newlines and a period at the end. E.g., 1, 2, 3, ..."),
			},
		},
		MaxTokens:      to.Ptr(int32(1024)),
		Temperature:    to.Ptr(float32(0.0)),
		DeploymentName: &deployment,
	}
}

var expectedContent = "1, 2, 3, 4, 5, 6, 7, 8, 9, 10."
var expectedRole = azopenai.ChatRoleAssistant

func TestClient_GetChatCompletions(t *testing.T) {
	client := newTestClient(t, azureOpenAI.Endpoint)
	testGetChatCompletions(t, client, azureOpenAI.ChatCompletionsRAI.Model, true)
}

func TestClient_GetChatCompletionsStream(t *testing.T) {
	chatClient := newTestClient(t, azureOpenAI.ChatCompletionsRAI.Endpoint)
	testGetChatCompletionsStream(t, chatClient, azureOpenAI.ChatCompletionsRAI.Model)
}

func TestClient_OpenAI_GetChatCompletions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping OpenAI tests when attempting to do quick tests")
	}

	chatClient := newOpenAIClientForTest(t)
	testGetChatCompletions(t, chatClient, openAI.ChatCompletions, false)
}

func TestClient_OpenAI_GetChatCompletionsStream(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping OpenAI tests when attempting to do quick tests")
	}

	chatClient := newOpenAIClientForTest(t)
	testGetChatCompletionsStream(t, chatClient, openAI.ChatCompletions)
}

func testGetChatCompletions(t *testing.T, client *azopenai.Client, deployment string, checkRAI bool) {
	expected := azopenai.ChatCompletions{
		Choices: []azopenai.ChatChoice{
			{
				Message: &azopenai.ChatResponseMessage{
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

	resp, err := client.GetChatCompletions(context.Background(), newTestChatCompletionOptions(deployment), nil)
	skipNowIfThrottled(t, err)
	require.NoError(t, err)

	if checkRAI {
		// Azure also provides content-filtering. This particular prompt and responses
		// will be considered safe.
		expected.PromptFilterResults = []azopenai.ContentFilterResultsForPrompt{
			{PromptIndex: to.Ptr[int32](0), ContentFilterResults: safeContentFilterResultDetailsForPrompt},
		}
		expected.Choices[0].ContentFilterResults = safeContentFilter
	}

	require.NotEmpty(t, resp.ID)
	require.NotEmpty(t, resp.Created)

	expected.ID = resp.ID
	expected.Created = resp.Created

	require.Equal(t, expected, resp.ChatCompletions)
}

func testGetChatCompletionsStream(t *testing.T, client *azopenai.Client, deployment string) {
	streamResp, err := client.GetChatCompletionsStream(context.Background(), newTestChatCompletionOptions(deployment), nil)

	if respErr := (*azcore.ResponseError)(nil); errors.As(err, &respErr) && respErr.StatusCode == http.StatusTooManyRequests {
		t.Skipf("OpenAI resource overloaded, skipping this test")
	}

	require.NoError(t, err)

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

		if completion.PromptFilterResults != nil {
			require.Equal(t, []azopenai.ContentFilterResultsForPrompt{
				{PromptIndex: to.Ptr[int32](0), ContentFilterResults: safeContentFilterResultDetailsForPrompt},
			}, completion.PromptFilterResults)
		}

		if len(completion.Choices) == 0 {
			// you can get empty entries that contain just metadata (ie, prompt annotations)
			continue
		}

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

	chatClient, err := azopenai.NewClient(azureOpenAI.Endpoint.URL, dac, &azopenai.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: recordingTransporter},
	})
	require.NoError(t, err)

	testGetChatCompletions(t, chatClient, azureOpenAI.ChatCompletions, true)
}

func TestClient_GetChatCompletions_InvalidModel(t *testing.T) {
	client := newTestClient(t, azureOpenAI.Endpoint)

	_, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestUserMessage{
				Content: azopenai.NewChatRequestUserMessageContent("Count to 100, with a comma between each number and no newlines. E.g., 1, 2, 3, ..."),
			},
		},
		MaxTokens:      to.Ptr(int32(1024)),
		Temperature:    to.Ptr(float32(0.0)),
		DeploymentName: to.Ptr("invalid model name"),
	}, nil)

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, "DeploymentNotFound", respErr.ErrorCode)
}

func TestClient_GetChatCompletionsStream_Error(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip()
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		client := newBogusAzureOpenAIClient(t)
		streamResp, err := client.GetChatCompletionsStream(context.Background(), newTestChatCompletionOptions(azureOpenAI.ChatCompletions), nil)
		require.Empty(t, streamResp)
		assertResponseIsError(t, err)
	})

	t.Run("OpenAI", func(t *testing.T) {
		client := newBogusOpenAIClient(t)
		streamResp, err := client.GetChatCompletionsStream(context.Background(), newTestChatCompletionOptions(openAI.ChatCompletions), nil)
		require.Empty(t, streamResp)
		assertResponseIsError(t, err)
	})
}

func TestClient_GetChatCompletions_Vision(t *testing.T) {
	testFn := func(t *testing.T, chatClient *azopenai.Client, deploymentName string) {
		imageURL := "https://www.bing.com/th?id=OHR.BradgateFallow_EN-US3932725763_1920x1080.jpg"

		content := azopenai.NewChatRequestUserMessageContent([]azopenai.ChatCompletionRequestMessageContentPartClassification{
			&azopenai.ChatCompletionRequestMessageContentPartText{
				Text: to.Ptr("Describe this image"),
			},
			&azopenai.ChatCompletionRequestMessageContentPartImage{
				ImageURL: &azopenai.ChatCompletionRequestMessageContentPartImageURL{
					URL: &imageURL,
				},
			},
		})

		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()

		resp, err := chatClient.GetChatCompletions(ctx, azopenai.ChatCompletionsOptions{
			Messages: []azopenai.ChatRequestMessageClassification{
				&azopenai.ChatRequestUserMessage{
					Content: content,
				},
			},
			DeploymentName: to.Ptr(deploymentName),
			MaxTokens:      to.Ptr[int32](512),
		}, nil)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Choices[0].Message.Content)

		t.Logf(*resp.Choices[0].Message.Content)
	}

	t.Run("OpenAI", func(t *testing.T) {
		chatClient := newOpenAIClientForTest(t)
		testFn(t, chatClient, openAI.Vision.Model)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		chatClient := newTestClient(t, azureOpenAI.Vision.Endpoint)
		testFn(t, chatClient, azureOpenAI.Vision.Model)
	})
}

func TestGetChatCompletions_usingResponseFormatForJSON(t *testing.T) {
	testFn := func(t *testing.T, chatClient *azopenai.Client, deploymentName string) {
		body := azopenai.ChatCompletionsOptions{
			DeploymentName: &deploymentName,
			Messages: []azopenai.ChatRequestMessageClassification{
				&azopenai.ChatRequestSystemMessage{Content: to.Ptr("You are a helpful assistant designed to output JSON.")},
				&azopenai.ChatRequestUserMessage{
					Content: azopenai.NewChatRequestUserMessageContent("List capital cities and their states"),
				},
			},
			// Without this format directive you end up getting JSON, but with a non-JSON preamble, like this:
			// "I'm happy to help! Here are some examples of capital cities and their corresponding states:\n\n```json\n{\n" (etc)
			ResponseFormat: &azopenai.ChatCompletionsJSONResponseFormat{},
			Temperature:    to.Ptr[float32](0.0),
		}

		resp, err := chatClient.GetChatCompletions(context.Background(), body, nil)
		require.NoError(t, err)

		// validate that it came back as JSON data
		var v any
		err = json.Unmarshal([]byte(*resp.Choices[0].Message.Content), &v)
		require.NoError(t, err)
		require.NotEmpty(t, v)
	}

	t.Run("OpenAI", func(t *testing.T) {
		chatClient := newOpenAIClientForTest(t)
		testFn(t, chatClient, "gpt-3.5-turbo-1106")
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		chatClient := newTestClient(t, azureOpenAI.DallE.Endpoint)
		testFn(t, chatClient, "gpt-4-1106-preview")
	})
}
