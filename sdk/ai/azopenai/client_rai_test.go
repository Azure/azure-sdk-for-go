//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
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
