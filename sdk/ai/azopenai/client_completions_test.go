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

func TestClient_GetCompletions(t *testing.T) {
	testFn := func(t *testing.T, epm endpointWithModel) {
		client := newTestClient(t, epm.Endpoint)

		resp, err := client.GetCompletions(context.Background(), azopenai.CompletionsOptions{
			Prompt:         []string{"What is Azure OpenAI?"},
			MaxTokens:      to.Ptr(int32(2048 - 127)),
			Temperature:    to.Ptr(float32(0.0)),
			DeploymentName: &epm.Model,
		}, nil)
		customRequireNoError(t, err, true)

		// we'll do a general check here - as models change the answers can also change, token usages are different,
		// etc... So we'll just make sure data is coming back and is reasonable.
		require.NotZero(t, *resp.Completions.Usage.PromptTokens)
		require.NotZero(t, *resp.Completions.Usage.CompletionTokens)
		require.NotZero(t, *resp.Completions.Usage.TotalTokens)
		require.Equal(t, int32(0), *resp.Completions.Choices[0].Index)
		require.Equal(t, azopenai.CompletionsFinishReasonStopped, *resp.Completions.Choices[0].FinishReason)

		require.NotEmpty(t, *resp.Completions.Choices[0].Text)

		if epm.Endpoint.Azure {
			require.Equal(t, safeContentFilter, resp.Completions.Choices[0].ContentFilterResults)
			require.Equal(t, []azopenai.ContentFilterResultsForPrompt{
				{
					PromptIndex:          to.Ptr[int32](0),
					ContentFilterResults: safeContentFilterResultDetailsForPrompt,
				}}, resp.PromptFilterResults)
		}

	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, azureOpenAI.Completions)
	})

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, openAI.Completions)
	})
}
