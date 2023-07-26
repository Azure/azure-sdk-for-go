//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/stretchr/testify/require"
)

func TestClient_GetCompletions_AzureOpenAI_ContentFilter_Response(t *testing.T) {
	// Scenario: Your API call asks for multiple responses (N>1) and at least 1 of the responses is filtered
	// https://github.com/MicrosoftDocs/azure-docs/blob/main/articles/cognitive-services/openai/concepts/content-filter.md#scenario-your-api-call-asks-for-multiple-responses-n1-and-at-least-1-of-the-responses-is-filtered
	client := newAzureOpenAIClientForTest(t, azureOpenAI)

	resp, err := client.GetCompletions(context.Background(), azopenai.CompletionsOptions{
		Prompt:       []string{"How do I rob a bank?"},
		MaxTokens:    to.Ptr(int32(2048 - 127)),
		Temperature:  to.Ptr(float32(0.0)),
		DeploymentID: azureOpenAI.Completions,
	}, nil)

	require.Empty(t, resp)
	assertContentFilterError(t, err, false)
}

func TestClient_GetChatCompletions_AzureOpenAI_ContentFilterWithError(t *testing.T) {
	client := newAzureOpenAIClientForTest(t, azureOpenAICanary)

	resp, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatMessage{
			{Role: to.Ptr(azopenai.ChatRoleSystem), Content: to.Ptr("You are a helpful assistant.")},
			{Role: to.Ptr(azopenai.ChatRoleUser), Content: to.Ptr("How do I rob a bank?")},
		},
		MaxTokens:    to.Ptr(int32(2048 - 127)),
		Temperature:  to.Ptr(float32(0.0)),
		DeploymentID: azureOpenAICanary.ChatCompletions,
	}, nil)
	require.Empty(t, resp)
	assertContentFilterError(t, err, true)
}

func TestClient_GetChatCompletions_AzureOpenAI_ContentFilter_WithResponse(t *testing.T) {
	client := newAzureOpenAIClientForTest(t, azureOpenAICanary)

	resp, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatMessage{
			{Role: to.Ptr(azopenai.ChatRoleUser), Content: to.Ptr("How do I cook a bell pepper?")},
		},
		MaxTokens:    to.Ptr(int32(2048 - 127)),
		Temperature:  to.Ptr(float32(0.0)),
		DeploymentID: azureOpenAICanary.ChatCompletions,
	}, nil)

	require.NoError(t, err)

	require.Equal(t, safeContentFilter, resp.ChatCompletions.Choices[0].ContentFilterResults)
}

// assertContentFilterError checks that the content filtering error came back from Azure OpenAI.
func assertContentFilterError(t *testing.T, err error, requireAnnotations bool) {
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, "content_filter", respErr.ErrorCode)

	require.Contains(t, respErr.Error(), "The response was filtered due to the prompt triggering")

	// Azure also returns error information when content filtering happens.
	var contentFilterErr *azopenai.ContentFilterResponseError
	require.ErrorAs(t, err, &contentFilterErr)

	if requireAnnotations {
		require.Equal(t, &azopenai.ContentFilterResultsHate{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)}, contentFilterErr.ContentFilterResults.Hate)
		require.Equal(t, &azopenai.ContentFilterResultsSelfHarm{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)}, contentFilterErr.ContentFilterResults.SelfHarm)
		require.Equal(t, &azopenai.ContentFilterResultsSexual{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)}, contentFilterErr.ContentFilterResults.Sexual)
		require.Equal(t, &azopenai.ContentFilterResultsViolence{Filtered: to.Ptr(true), Severity: to.Ptr(azopenai.ContentFilterSeverityMedium)}, contentFilterErr.ContentFilterResults.Violence)
	}
}

var safeContentFilter = &azopenai.ChatChoiceContentFilterResults{
	Hate:     &azopenai.ContentFilterResultsHate{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
	SelfHarm: &azopenai.ContentFilterResultsSelfHarm{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
	Sexual:   &azopenai.ContentFilterResultsSexual{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
	Violence: &azopenai.ContentFilterResultsViolence{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)},
}
