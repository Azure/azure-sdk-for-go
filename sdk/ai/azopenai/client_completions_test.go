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

func TestClient_GetCompletions_AzureOpenAI(t *testing.T) {
	client := newTestClient(t, azureOpenAI.Endpoint)
	testGetCompletions(t, client, true)
}

func TestClient_GetCompletions_OpenAI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping OpenAI tests when attempting to do quick tests")
	}

	client := newOpenAIClientForTest(t)
	testGetCompletions(t, client, false)
}

func testGetCompletions(t *testing.T, client *azopenai.Client, isAzure bool) {
	deploymentID := openAI.Completions

	if isAzure {
		deploymentID = azureOpenAI.Completions
	}

	resp, err := client.GetCompletions(context.Background(), azopenai.CompletionsOptions{
		Prompt:         []string{"What is Azure OpenAI?"},
		MaxTokens:      to.Ptr(int32(2048 - 127)),
		Temperature:    to.Ptr(float32(0.0)),
		DeploymentName: &deploymentID,
	}, nil)
	skipNowIfThrottled(t, err)
	require.NoError(t, err)

	// we'll do a general check here - as models change the answers can also change, token usages are different,
	// etc... So we'll just make sure data is coming back and is reasonable.
	require.NotZero(t, *resp.Completions.Usage.PromptTokens)
	require.NotZero(t, *resp.Completions.Usage.CompletionTokens)
	require.NotZero(t, *resp.Completions.Usage.TotalTokens)
	require.Equal(t, int32(0), *resp.Completions.Choices[0].Index)
	require.Equal(t, azopenai.CompletionsFinishReasonStopped, *resp.Completions.Choices[0].FinishReason)

	require.NotEmpty(t, *resp.Completions.Choices[0].Text)

	if isAzure {
		require.Equal(t, safeContentFilter, resp.Completions.Choices[0].ContentFilterResults)
		require.Equal(t, []azopenai.ContentFilterResultsForPrompt{
			{
				PromptIndex:          to.Ptr[int32](0),
				ContentFilterResults: safeContentFilterResultDetailsForPrompt,
			}}, resp.PromptFilterResults)
	}

}
