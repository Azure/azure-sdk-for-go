//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	type args struct {
		endpoint   string
		credential azcore.TokenCredential
		options    *azopenai.ClientOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *azopenai.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := azopenai.NewClient(tt.args.endpoint, tt.args.credential, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientWithKeyCredential(t *testing.T) {
	type args struct {
		endpoint   string
		credential *azcore.KeyCredential
		options    *azopenai.ClientOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *azopenai.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := azopenai.NewClientWithKeyCredential(tt.args.endpoint, tt.args.credential, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientWithKeyCredential() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientWithKeyCredential() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCompletionsStream(t *testing.T) {
	testFn := func(t *testing.T, epm endpointWithModel) {
		body := azopenai.CompletionsOptions{
			Prompt:         []string{"What is Azure OpenAI?"},
			MaxTokens:      to.Ptr(int32(2048)),
			Temperature:    to.Ptr(float32(0.0)),
			DeploymentName: &epm.Model,
		}

		client := newTestClient(t, epm.Endpoint)

		response, err := client.GetCompletionsStream(context.TODO(), body, nil)
		skipNowIfThrottled(t, err)
		require.NoError(t, err)

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

			if err == io.EOF {
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

		streamResp, err := client.GetCompletionsStream(context.Background(), azopenai.CompletionsOptions{
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
