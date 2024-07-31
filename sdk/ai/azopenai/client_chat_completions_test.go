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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
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
	testFn := func(t *testing.T, client *azopenai.Client, deployment string, returnedModel string, checkRAI bool) {
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
			Model: &returnedModel,
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

		t.Logf("isAzure: %t, deployment: %s, returnedModel: %s", checkRAI, deployment, *resp.ChatCompletions.Model)
		require.Equal(t, expected, resp.ChatCompletions)
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		client := newTestClient(t, azureOpenAI.ChatCompletionsRAI.Endpoint)
		testFn(t, client, azureOpenAI.ChatCompletionsRAI.Model, "gpt-4", true)
	})

	t.Run("AzureOpenAI.TokenCredential", func(t *testing.T) {
		if recording.GetRecordMode() == recording.PlaybackMode {
			t.Skipf("Not running this test in playback (for now)")
		}

		if os.Getenv("USE_TOKEN_CREDS") != "true" {
			t.Skipf("USE_TOKEN_CREDS is not true, disabling token credential tests")
		}

		recordingTransporter := newRecordingTransporter(t)

		cred, err := credential.New(nil)
		require.NoError(t, err)

		chatClient, err := azopenai.NewClient(azureOpenAI.ChatCompletions.Endpoint.URL, cred, &azopenai.ClientOptions{
			ClientOptions: policy.ClientOptions{Transport: recordingTransporter},
		})
		require.NoError(t, err)

		testFn(t, chatClient, azureOpenAI.ChatCompletions.Model, "gpt-4", true)
	})

	t.Run("OpenAI", func(t *testing.T) {
		chatClient := newTestClient(t, openAI.ChatCompletions.Endpoint)
		testFn(t, chatClient, openAI.ChatCompletions.Model, "gpt-4-0613", false)
	})
}

func TestClient_GetChatCompletions_LogProbs(t *testing.T) {
	testFn := func(t *testing.T, epm endpointWithModel) {
		client := newTestClient(t, epm.Endpoint)

		opts := azopenai.ChatCompletionsOptions{
			Messages: []azopenai.ChatRequestMessageClassification{
				&azopenai.ChatRequestUserMessage{
					Content: azopenai.NewChatRequestUserMessageContent("Count to 10, with a comma between each number, no newlines and a period at the end. E.g., 1, 2, 3, ..."),
				},
			},
			MaxTokens:      to.Ptr(int32(1024)),
			Temperature:    to.Ptr(float32(0.0)),
			DeploymentName: &epm.Model,
			LogProbs:       to.Ptr(true),
			TopLogProbs:    to.Ptr(int32(5)),
		}

		resp, err := client.GetChatCompletions(context.Background(), opts, nil)
		require.NoError(t, err)

		for _, choice := range resp.Choices {
			require.NotEmpty(t, choice.LogProbs)
		}
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, azureOpenAI.ChatCompletions)
	})

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, openAI.ChatCompletions)
	})
}

func TestClient_GetChatCompletions_LogitBias(t *testing.T) {
	// you can use LogitBias to constrain the answer to NOT contain
	// certain tokens. More or less following the technique in this OpenAI article:
	// https://help.openai.com/en/articles/5247780-using-logit-bias-to-alter-token-probability-with-the-openai-api

	testFn := func(t *testing.T, epm endpointWithModel) {
		client := newTestClient(t, epm.Endpoint)

		opts := azopenai.ChatCompletionsOptions{
			Messages: []azopenai.ChatRequestMessageClassification{
				&azopenai.ChatRequestUserMessage{
					Content: azopenai.NewChatRequestUserMessageContent("Briefly, what are some common roles for people at a circus, names only, one per line?"),
				},
			},
			MaxTokens:      to.Ptr(int32(200)),
			Temperature:    to.Ptr(float32(0.0)),
			DeploymentName: &epm.Model,
			LogitBias: map[string]*int32{
				// you can calculate these tokens using OpenAI's online tool:
				// https://platform.openai.com/tokenizer?view=bpe
				// These token IDs are all variations of "Clown", which I want to exclude from the response.
				"25":    to.Ptr(int32(-100)),
				"220":   to.Ptr(int32(-100)),
				"1206":  to.Ptr(int32(-100)),
				"2493":  to.Ptr(int32(-100)),
				"5176":  to.Ptr(int32(-100)),
				"43456": to.Ptr(int32(-100)),
				"99423": to.Ptr(int32(-100)),
			},
		}

		resp, err := client.GetChatCompletions(context.Background(), opts, nil)
		require.NoError(t, err)

		for _, choice := range resp.Choices {
			if choice.Message == nil || choice.Message.Content == nil {
				continue
			}

			require.NotContains(t, *choice.Message.Content, "clown")
			require.NotContains(t, *choice.Message.Content, "Clown")
		}
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, azureOpenAI.ChatCompletions)
	})

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, openAI.ChatCompletions)
	})
}

func TestClient_GetChatCompletionsStream(t *testing.T) {
	testFn := func(t *testing.T, client *azopenai.Client, deployment string, returnedDeployment string) {
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

		modelWasReturned := false

		for {
			completion, err := streamResp.ChatCompletionsStream.Read()

			if errors.Is(err, io.EOF) {
				break
			}

			// NOTE: this is actually the name of the _model_, not the deployment. They usually match (just
			// by convention) but if this fails because they _don't_ match we can just adjust the test.
			if returnedDeployment == *completion.Model {
				modelWasReturned = true
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

		require.True(t, modelWasReturned)

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

	t.Run("AzureOpenAI", func(t *testing.T) {
		chatClient := newTestClient(t, azureOpenAI.ChatCompletionsRAI.Endpoint)
		testFn(t, chatClient, azureOpenAI.ChatCompletionsRAI.Model, "gpt-4")
	})

	t.Run("OpenAI", func(t *testing.T) {
		chatClient := newTestClient(t, openAI.ChatCompletions.Endpoint)
		testFn(t, chatClient, openAI.ChatCompletions.Model, openAI.ChatCompletions.Model)
	})
}

func TestClient_GetChatCompletions_InvalidModel(t *testing.T) {
	client := newTestClient(t, azureOpenAI.ChatCompletions.Endpoint)

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
		streamResp, err := client.GetChatCompletionsStream(context.Background(), newTestChatCompletionOptions(azureOpenAI.ChatCompletions.Model), nil)
		require.Empty(t, streamResp)
		assertResponseIsError(t, err)
	})

	t.Run("OpenAI", func(t *testing.T) {
		client := newBogusOpenAIClient(t)
		streamResp, err := client.GetChatCompletionsStream(context.Background(), newTestChatCompletionOptions(openAI.ChatCompletions.Model), nil)
		require.Empty(t, streamResp)
		assertResponseIsError(t, err)
	})
}

func TestClient_GetChatCompletions_Vision(t *testing.T) {
	testFn := func(t *testing.T, chatClient *azopenai.Client, deploymentName string, azure bool) {
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

		// vision is a bit of an oversubscribed Azure resource. Allow 429, but mark the test as skipped.
		customRequireNoError(t, err, azure)
		require.NotEmpty(t, resp.Choices[0].Message.Content)

		t.Logf(*resp.Choices[0].Message.Content)
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		chatClient := newTestClient(t, azureOpenAI.Vision.Endpoint)
		testFn(t, chatClient, azureOpenAI.Vision.Model, true)
	})

	t.Run("OpenAI", func(t *testing.T) {
		chatClient := newTestClient(t, openAI.Vision.Endpoint)
		testFn(t, chatClient, openAI.Vision.Model, false)
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

	t.Run("AzureOpenAI", func(t *testing.T) {
		chatClient := newTestClient(t, azureOpenAI.ChatCompletionsWithJSONResponseFormat.Endpoint)
		testFn(t, chatClient, azureOpenAI.ChatCompletionsWithJSONResponseFormat.Model)
	})

	t.Run("OpenAI", func(t *testing.T) {
		chatClient := newTestClient(t, openAI.ChatCompletionsWithJSONResponseFormat.Endpoint)
		testFn(t, chatClient, openAI.ChatCompletionsWithJSONResponseFormat.Model)
	})
}
