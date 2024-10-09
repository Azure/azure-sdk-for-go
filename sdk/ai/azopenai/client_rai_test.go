//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

// RAI == "responsible AI". This part of the API provides content filtering and
// classification of the failures into categories like Hate, Violence, etc...

func TestClient_GetChatCompletions_AzureOpenAI_ContentFilter_WithResponse(t *testing.T) {
	client := newTestClient(t, azureOpenAI.ChatCompletionsRAI.Endpoint)

	resp, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("How do I cook a bell pepper?")},
		},
		MaxTokens:      to.Ptr(int32(2048 - 127)),
		Temperature:    to.Ptr(float32(0.0)),
		DeploymentName: &azureOpenAI.ChatCompletionsRAI.Model,
	}, nil)
	customRequireNoError(t, err, true)

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
		require.NotNil(t, contentFilterErr.ContentFilterResults)

		require.Equal(t, &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)}, contentFilterErr.ContentFilterResults.Hate)
		require.Equal(t, &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)}, contentFilterErr.ContentFilterResults.SelfHarm)
		require.Equal(t, &azopenai.ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(azopenai.ContentFilterSeveritySafe)}, contentFilterErr.ContentFilterResults.Sexual)
		require.Equal(t, &azopenai.ContentFilterResult{Filtered: to.Ptr(true), Severity: to.Ptr(azopenai.ContentFilterSeverityMedium)}, contentFilterErr.ContentFilterResults.Violence)
	}
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
