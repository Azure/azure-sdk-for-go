// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/require"
)

func TestClient_GetCompletions(t *testing.T) {
	t.Skip("Disabled until we find a compatible model")

	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.Completions.Endpoint)

	resp, err := client.Completions.New(context.Background(), openai.CompletionNewParams{
		Prompt: openai.CompletionNewParamsPromptUnion{
			OfArrayOfStrings: []string{"What is Azure OpenAI?"},
		},
		MaxTokens:   openai.Int(2048 - 127),
		Temperature: openai.Float(0.0),
		Model:       openai.CompletionNewParamsModel(azureOpenAI.Completions.Model),
	})
	skipNowIfThrottled(t, err)
	require.NoError(t, err)

	// we'll do a general check here - as models change the answers can also change, token usages are different,
	// etc... So we'll just make sure data is coming back and is reasonable.
	require.NotZero(t, resp.Usage.PromptTokens)
	require.NotZero(t, resp.Usage.CompletionTokens)
	require.NotZero(t, resp.Usage.TotalTokens)
	require.Equal(t, int64(0), resp.Choices[0].Index)
	require.Equal(t, openai.CompletionChoiceFinishReasonStop, resp.Choices[0].FinishReason)

	require.NotEmpty(t, resp.Choices[0].Text)

	azureChoice := azopenai.CompletionChoice(resp.Choices[0])
	contentFilterResults, err := azureChoice.ContentFilterResults()
	require.NoError(t, err)

	require.Equal(t, safeContentFilter, contentFilterResults)

	azureCompletion := azopenai.Completion(*resp)
	promptFilterResults, err := azureCompletion.PromptFilterResults()
	require.NoError(t, err)

	require.Equal(t, []azopenai.ContentFilterResultsForPrompt{{
		PromptIndex:          to.Ptr[int32](0),
		ContentFilterResults: safeContentFilterResultDetailsForPrompt,
	}}, promptFilterResults)
}

func TestGetCompletionsStream(t *testing.T) {
	t.Skip("Disabled until we find a compatible model")

	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.Completions.Endpoint)

	stream := client.Completions.NewStreaming(context.TODO(), openai.CompletionNewParams{
		Model:       openai.CompletionNewParamsModel(azureOpenAI.Completions.Model),
		MaxTokens:   openai.Int(2048),
		Temperature: openai.Float(0.0),
		Prompt: openai.CompletionNewParamsPromptUnion{
			OfArrayOfStrings: []string{"What is Azure OpenAI?"},
		},
	})

	t.Cleanup(func() {
		err := stream.Close()
		require.NoError(t, err)
	})

	var sb strings.Builder
	var eventCount int

	for stream.Next() {
		completion := azopenai.Completion(stream.Current())

		promptFilterResults, err := completion.PromptFilterResults()
		require.NoError(t, err)

		if promptFilterResults != nil {
			require.Equal(t, []azopenai.ContentFilterResultsForPrompt{
				{PromptIndex: to.Ptr[int32](0), ContentFilterResults: safeContentFilterResultDetailsForPrompt},
			}, promptFilterResults)
		}

		eventCount++

		if len(completion.Choices) > 0 {
			sb.WriteString(completion.Choices[0].Text)
		}
	}

	require.NoError(t, stream.Err())

	got := sb.String()

	require.NotEmpty(t, got)

	// there's no strict requirement of how the response is streamed so just
	// choosing something that's reasonable but will be lower than typical usage
	// (which is usually somewhere around the 80s).
	require.GreaterOrEqual(t, eventCount, 50)
}
