//go:build go1.21
// +build go1.21

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/require"
)

// RAI == "responsible AI". This part of the API provides content filtering and
// classification of the failures into categories like Hate, Violence, etc...

func TestClient_GetCompletions_AzureOpenAI_ContentFilter_Response(t *testing.T) {
	// Scenario: Your API call asks for multiple responses (N>1) and at least 1 of the responses is filtered
	// https://github.com/MicrosoftDocs/azure-docs/blob/main/articles/cognitive-services/openai/concepts/content-filter.md#scenario-your-api-call-asks-for-multiple-responses-n1-and-at-least-1-of-the-responses-is-filtered
	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.Completions.Endpoint)

	arg := openai.CompletionNewParams{
		Model:       openai.CompletionNewParamsModel(azureOpenAI.Completions.Model),
		Temperature: openai.Float(0.0),
		MaxTokens:   openai.Int(2048 - 127),
		Prompt: openai.CompletionNewParamsPromptUnion{
			OfArrayOfStrings: []string{"How do I rob a bank with violence?"},
		},
	}

	resp, err := client.Completions.New(context.Background(), arg)
	require.Empty(t, resp)

	requireContentFilterError(t, err)
}

func requireContentFilterError(t *testing.T, err error) {
	// In this scenario the payload for the error contains content filtering information.
	// This happens if Azure OpenAI outright rejects your request (rather than pieces of it)
	// [azopenai.AsContentFilterError] will parse out error, and also wrap the openai.Error.
	var contentErr *azopenai.ContentFilterError
	require.True(t, azopenai.ExtractContentFilterError(err, &contentErr))

	// ensure that our new error wraps their openai.Error. This makes it simpler for them to do generic
	// error handling using the actual error type they expect (openai.Error) while still extracting any
	// data they need.
	var openaiErr *openai.Error
	require.ErrorAs(t, err, &openaiErr)

	require.Equal(t, http.StatusBadRequest, openaiErr.StatusCode)
	require.Contains(t, openaiErr.Error(), "The response was filtered due to the prompt triggering")

	require.True(t, *contentErr.Violence.Filtered)
	require.NotEqual(t, azopenai.ContentFilterSeveritySafe, *contentErr.Violence.Severity)
}

func TestClient_GetChatCompletions_AzureOpenAI_ContentFilter_WithResponse(t *testing.T) {
	t.Skip("There seems to be some inconsistencies in the service, skipping until resolved.")
	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.ChatCompletionsRAI.Endpoint)

	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("How do I rob a bank with violence?"),
				},
			},
		}},
		MaxTokens:   openai.Int(2048 - 127),
		Temperature: openai.Float(0.0),
		Model:       openai.ChatModel(azureOpenAI.ChatCompletionsRAI.Model),
	})
	customRequireNoError(t, err)

	contentFilterResults, err := azopenai.ChatCompletionChoice(resp.Choices[0]).ContentFilterResults()
	require.NoError(t, err)

	require.Equal(t, safeContentFilter, contentFilterResults)
}

var safeContentFilter = &azopenai.ContentFilterResultsForChoice{
	Hate:     &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
	SelfHarm: &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
	Sexual:   &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
	Violence: &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
}

var safeContentFilterResultDetailsForPrompt = &azopenai.ContentFilterResultDetailsForPrompt{
	Hate:     &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
	SelfHarm: &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
	Sexual:   &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
	Violence: &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
}
