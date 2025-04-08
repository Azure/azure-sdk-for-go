//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestGetCompletionsStream(t *testing.T) {
	testFn := func(t *testing.T, epm endpointWithModel) {
		body := azopenai.CompletionsStreamOptions{
			Prompt:         []string{"What is Azure OpenAI?"},
			MaxTokens:      to.Ptr(int32(2048)),
			Temperature:    to.Ptr(float32(0.0)),
			DeploymentName: &epm.Model,
		}

		client := newTestClient(t, epm.Endpoint)

		response, err := client.GetCompletionsStream(context.TODO(), body, nil)
		customRequireNoError(t, err, true)

		if err != nil {
			t.Errorf("Client.GetCompletionsStream() error = %v", err)
			return
		}

		reader := response.CompletionsStream
		defer reader.Close()

		var sb strings.Builder
		var eventCount int

		for {
			completion, err := reader.Read()

			if errors.Is(err, io.EOF) {
				break
			}

			if completion.PromptFilterResults != nil {
				require.Equal(t, []azopenai.ContentFilterResultsForPrompt{
					{PromptIndex: to.Ptr[int32](0), ContentFilterResults: safeContentFilterResultDetailsForPrompt},
				}, completion.PromptFilterResults)
			}

			eventCount++

			if err != nil {
				t.Errorf("reader.Read() error = %v", err)
				return
			}

			if len(completion.Choices) > 0 {
				sb.WriteString(*completion.Choices[0].Text)
			}
		}
		got := sb.String()

		require.NotEmpty(t, got)

		// there's no strict requirement of how the response is streamed so just
		// choosing something that's reasonable but will be lower than typical usage
		// (which is usually somewhere around the 80s).
		require.GreaterOrEqual(t, eventCount, 50)
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, azureOpenAI.Completions)
	})

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, openAI.Completions)
	})
}

func TestClient_GetCompletions_Error(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip()
	}

	doTest := func(t *testing.T, model string) {
		client := newBogusAzureOpenAIClient(t)

		streamResp, err := client.GetCompletionsStream(context.Background(), azopenai.CompletionsStreamOptions{
			Prompt:         []string{"What is Azure OpenAI?"},
			MaxTokens:      to.Ptr(int32(2048 - 127)),
			Temperature:    to.Ptr(float32(0.0)),
			DeploymentName: &model,
		}, nil)
		require.Empty(t, streamResp)
		assertResponseIsError(t, err)
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		doTest(t, azureOpenAI.Completions.Model)
	})

	t.Run("OpenAI", func(t *testing.T) {
		doTest(t, openAI.Completions.Model)
	})
}
