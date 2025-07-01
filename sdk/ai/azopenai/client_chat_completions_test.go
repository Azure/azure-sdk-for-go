//go:build go1.21
// +build go1.21

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"

	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared/constant"
	"github.com/stretchr/testify/require"
)

func newStainlessTestChatCompletionOptions(deployment string) openai.ChatCompletionNewParams {
	message := "Count to 10, with a comma between each number, no newlines and a period at the end. E.g., 1, 2, 3, ..."

	return openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String(message),
				},
			},
		}},
		MaxTokens:   openai.Int(1024),
		Temperature: openai.Float(0.0),
		Model:       openai.ChatModel(deployment),
	}
}

var expectedContent = "1, 2, 3, 4, 5, 6, 7, 8, 9, 10."
var expectedRole = constant.ValueOf[constant.Assistant]()

func TestClient_GetChatCompletions(t *testing.T) {
	testFn := func(t *testing.T, client *openai.ChatCompletionService, deployment string, returnedModel string, checkRAI bool) {
		resp, err := client.New(context.Background(), newStainlessTestChatCompletionOptions(deployment))
		skipNowIfThrottled(t, err)
		require.NoError(t, err)

		require.NotEmpty(t, resp.ID)
		require.NotEmpty(t, resp.Created)

		t.Logf("isAzure: %t, deployment: %s, returnedModel: %s", checkRAI, deployment, resp.Model)

		// check Choices
		require.Equal(t, 1, len(resp.Choices))
		choice := resp.Choices[0]

		t.Logf("Content = %s", choice.Message.Content)

		require.Zero(t, choice.Index)
		require.EqualValues(t, "assistant", choice.Message.Role)
		require.NotEmpty(t, choice.Message.Content)
		require.Equal(t, "stop", choice.FinishReason)

		require.Equal(t, openai.CompletionUsage{
			// these change depending on which model you use. These #'s work for gpt-4, which is
			// what I'm using for these tests.
			CompletionTokens: 29,
			PromptTokens:     42,
			TotalTokens:      71,
		}, openai.CompletionUsage{
			CompletionTokens: resp.Usage.CompletionTokens,
			PromptTokens:     resp.Usage.PromptTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		})
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		client := newStainlessTestClient(t, azureOpenAI.ChatCompletionsRAI.Endpoint)

		testFn(t, &client.Chat.Completions, azureOpenAI.ChatCompletionsRAI.Model, "gpt-4", true)
	})

	t.Run("AzureOpenAI.DefaultAzureCredential", func(t *testing.T) {
		client := newStainlessTestClient(t, azureOpenAI.ChatCompletionsRAI.Endpoint)
		testFn(t, &client.Chat.Completions, azureOpenAI.ChatCompletions.Model, "gpt-4", true)
	})
}

func TestClient_GetChatCompletions_LogProbs(t *testing.T) {
	testFn := func(t *testing.T, client *openai.ChatCompletionService, model string) {
		opts := openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("Count to 10, with a comma between each number, no newlines and a period at the end. E.g., 1, 2, 3, ..."),
					},
				},
			}},
			MaxTokens:   openai.Int(1024),
			Temperature: openai.Float(0.0),
			Model:       openai.ChatModel(model),
			Logprobs:    openai.Bool(true),
			TopLogprobs: openai.Int(5),
		}

		resp, err := client.New(context.Background(), opts)
		require.NoError(t, err)

		for _, choice := range resp.Choices {
			require.NotEmpty(t, choice.Logprobs)
		}
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		client := newStainlessTestClient(t, azureOpenAI.ChatCompletions.Endpoint)
		testFn(t, &client.Chat.Completions, azureOpenAI.ChatCompletions.Model)
	})

	t.Run("AzureOpenAI.Service", func(t *testing.T) {
		client := newStainlessChatCompletionService(t, azureOpenAI.ChatCompletions.Endpoint)
		testFn(t, &client, azureOpenAI.ChatCompletions.Model)
	})
}

func TestClient_GetChatCompletions_LogitBias(t *testing.T) {
	// you can use LogitBias to constrain the answer to NOT contain
	// certain tokens. More or less following the technique in this OpenAI article:
	// https://help.openai.com/en/articles/5247780-using-logit-bias-to-alter-token-probability-with-the-openai-api

	testFn := func(t *testing.T, epm endpointWithModel) {
		client := newStainlessTestClient(t, epm.Endpoint)

		body := openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String("Briefly, what are some common roles for people at a circus, names only, one per line?"),
					},
				},
			}},
			MaxTokens:   openai.Int(200),
			Temperature: openai.Float(0.0),
			Model:       openai.ChatModel(epm.Model),
			LogitBias: map[string]int64{
				// you can calculate these tokens using OpenAI's online tool:
				// https://platform.openai.com/tokenizer?view=bpe
				// These token IDs are all variations of "Clown", which I want to exclude from the response.
				"25":    -100,
				"220":   -100,
				"1206":  -100,
				"2493":  -100,
				"5176":  -100,
				"43456": -100,
				"69568": -100,
				"99423": -100,
			},
		}

		resp, err := client.Chat.Completions.New(context.Background(), body)
		require.NoError(t, err)

		for _, choice := range resp.Choices {
			require.NotContains(t, choice.Message.Content, "clown")
			require.NotContains(t, choice.Message.Content, "Clown")
		}
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, azureOpenAI.ChatCompletions)
	})
}

func TestClient_GetChatCompletionsStream(t *testing.T) {
	chatClient := newStainlessTestClient(t, azureOpenAI.ChatCompletionsRAI.Endpoint)
	stream := chatClient.Chat.Completions.NewStreaming(context.Background(), newStainlessTestChatCompletionOptions(azureOpenAI.ChatCompletionsRAI.Model))

	// the data comes back differently for streaming
	// 1. the text comes back in the ChatCompletion.Delta field
	// 2. the role is only sent on the first streamed ChatCompletion
	// check that the role came back as well.
	var choices []openai.ChatCompletionChunkChoice

	modelWasReturned := false

	for stream.Next() {
		chunk := stream.Current()

		// NOTE: this is actually the name of the _model_, not the deployment. They usually match (just
		// by convention) but if this fails because they _don't_ match we can just adjust the test.
		if len(chunk.Model) > 0 {
			modelWasReturned = true
		}

		azureChunk := azopenai.ChatCompletionChunk(chunk)

		// NOTE: prompt filter results are non-deterministic as they're based on their own criteria, which
		// can change over time. We'll check that we can safely attempt to deserialize it.
		_, err := azureChunk.PromptFilterResults()
		require.NoError(t, err)

		if len(chunk.Choices) == 0 {
			// you can get empty entries that contain just metadata (ie, prompt annotations)
			continue
		}

		require.Equal(t, 1, len(chunk.Choices))
		choices = append(choices, chunk.Choices[0])
	}

	require.NoError(t, stream.Err())

	require.True(t, modelWasReturned)

	var message string
	var role constant.Assistant

	for _, choice := range choices {
		message += choice.Delta.Content
		if len(choice.Delta.Role) > 0 {
			role = constant.Assistant(choice.Delta.Role)
		}
	}

	require.Equal(t, expectedContent, message)
	require.Equal(t, expectedRole, role)
}

func TestClient_GetChatCompletions_Vision(t *testing.T) {
	// testFn := func(t *testing.T, chatClient *azopenai.Client, deploymentName string, azure bool) {
	chatClient := newStainlessTestClient(t, azureOpenAI.Vision.Endpoint)

	imageURL := "https://www.bing.com/th?id=OHR.BradgateFallow_EN-US3932725763_1920x1080.jpg"

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	resp, err := chatClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfArrayOfContentParts: []openai.ChatCompletionContentPartUnionParam{{
						OfText: &openai.ChatCompletionContentPartTextParam{
							Text: "Describe this image",
						},
					}, {
						OfImageURL: &openai.ChatCompletionContentPartImageParam{
							ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
								URL: imageURL,
							},
						},
					}},
				},
			},
		}},
		Model:     openai.ChatModel(azureOpenAI.Vision.Model),
		MaxTokens: openai.Int(512),
	})

	// vision is a bit of an oversubscribed Azure resource. Allow 429, but mark the test as skipped.
	customRequireNoError(t, err)
	require.NotEmpty(t, resp.Choices[0].Message.Content)

	t.Logf("Content: %s", resp.Choices[0].Message.Content)
}
