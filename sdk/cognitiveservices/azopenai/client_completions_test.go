//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/stretchr/testify/require"
)

func TestClient_GetCompletions_AzureOpenAI(t *testing.T) {
	cred, err := azopenai.NewKeyCredential(azureOpenAI.APIKey)
	require.NoError(t, err)

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAI.Endpoint, cred, newClientOptionsForTest(t))
	require.NoError(t, err)

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
		Prompt:       []string{"What is Azure OpenAI?"},
		MaxTokens:    to.Ptr(int32(2048 - 127)),
		Temperature:  to.Ptr(float32(0.0)),
		DeploymentID: deploymentID,
	}, nil)
	require.NoError(t, err)

	want := azopenai.GetCompletionsResponse{
		Completions: azopenai.Completions{
			Choices: []azopenai.Choice{
				{
					Text:         to.Ptr("\n\nAzure OpenAI is a platform from Microsoft that provides access to OpenAI's artificial intelligence (AI) technologies. It enables developers to build, train, and deploy AI models in the cloud. Azure OpenAI provides access to OpenAI's powerful AI technologies, such as GPT-3, which can be used to create natural language processing (NLP) applications, computer vision models, and reinforcement learning models."),
					Index:        to.Ptr(int32(0)),
					FinishReason: to.Ptr(azopenai.CompletionsFinishReason("stop")),
					LogProbs:     nil,
				},
			},
			Usage: &azopenai.CompletionsUsage{
				CompletionTokens: to.Ptr(int32(85)),
				PromptTokens:     to.Ptr(int32(6)),
				TotalTokens:      to.Ptr(int32(91)),
			},
		},
	}

	want.ID = resp.Completions.ID
	want.Created = resp.Completions.Created

	require.Equal(t, want, resp)
}
