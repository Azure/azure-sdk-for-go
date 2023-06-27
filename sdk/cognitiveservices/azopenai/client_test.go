//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/stretchr/testify/require"
)

func TestClient_OpenAI_InvalidModel(t *testing.T) {
	chatClient := newOpenAIClientForTest(t)

	_, err := chatClient.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		Messages: []*azopenai.ChatMessage{
			{
				Role:    to.Ptr(azopenai.ChatRoleSystem),
				Content: to.Ptr("hello"),
			},
		},
		Model: to.Ptr("non-existent-model"),
	}, nil)

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusNotFound, respErr.StatusCode)
	require.Contains(t, respErr.Error(), "The model `non-existent-model` does not exist")
}

func newOpenAIClientForTest(t *testing.T) *azopenai.Client {
	if openAIKey == "" {
		t.Skipf("OPENAI_API_KEY not defined, skipping OpenAI public endpoint test")
	}

	cred, err := azopenai.NewKeyCredential(openAIKey)
	require.NoError(t, err)

	chatClient, err := azopenai.NewClientForOpenAI(openAIEndpoint, cred, newClientOptionsForTest(t))
	require.NoError(t, err)

	return chatClient
}
