//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"context"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestClient_GetChatCompletions(t *testing.T) {
	type args struct {
		ctx          context.Context
		deploymentID string
		body         ChatCompletionsOptions
		options      *ClientGetChatCompletionsOptions
	}
	cred := KeyCredential{APIKey: apiKey}
	deploymentID := "gpt-35-turbo"
	chatClient, err := NewClientWithKeyCredential(endpoint, cred, deploymentID, newClientOptionsForTest(t))
	if err != nil {
		log.Fatalf("%v", err)
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    ClientGetChatCompletionsResponse
		wantErr bool
	}{
		{
			name:   "ChatCompletions",
			client: chatClient,
			args: args{
				ctx:          context.TODO(),
				deploymentID: "gpt-35-turbo",
				body: ChatCompletionsOptions{
					Messages: []*ChatMessage{
						{
							Role:    to.Ptr(ChatRole("user")),
							Content: to.Ptr("Count to 100, with a comma between each number and no newlines. E.g., 1, 2, 3, ..."),
						},
					},
					MaxTokens:   to.Ptr(int32(1024)),
					Temperature: to.Ptr(float32(0.0)),
				},
				options: nil,
			},
			want: ClientGetChatCompletionsResponse{
				ChatCompletions: ChatCompletions{
					Choices: []*ChatChoice{
						{
							Message: &ChatChoiceMessage{
								Role:    to.Ptr(ChatRole("assistant")),
								Content: to.Ptr("1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100."),
							},
							Index:        to.Ptr(int32(0)),
							FinishReason: to.Ptr(CompletionsFinishReason("stop")),
						},
					},
					Usage: &CompletionsUsage{
						CompletionTokens: to.Ptr(int32(299)),
						PromptTokens:     to.Ptr(int32(37)),
						TotalTokens:      to.Ptr(int32(336)),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetChatCompletions(tt.args.ctx, tt.args.body, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetChatCompletions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opts := cmpopts.IgnoreFields(ChatCompletions{}, "Created", "ID")
			if diff := cmp.Diff(tt.want.ChatCompletions, got.ChatCompletions, opts); diff != "" {
				t.Errorf("Client.GetCompletions(): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestClient_GetChatCompletions_InvalidModel(t *testing.T) {
	cred := KeyCredential{APIKey: apiKey}
	chatClient, err := NewClientWithKeyCredential(endpoint, cred, "thisdoesntexist", newClientOptionsForTest(t))
	require.NoError(t, err)

	_, err = chatClient.GetChatCompletions(context.Background(), ChatCompletionsOptions{
		Messages: []*ChatMessage{
			{
				Role:    to.Ptr(ChatRole("user")),
				Content: to.Ptr("Count to 100, with a comma between each number and no newlines. E.g., 1, 2, 3, ..."),
			},
		},
		MaxTokens:   to.Ptr(int32(1024)),
		Temperature: to.Ptr(float32(0.0)),
	}, nil)

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, "DeploymentNotFound", respErr.ErrorCode)
}

func TestClient_GetEmbeddings_InvalidModel(t *testing.T) {
	cred := KeyCredential{APIKey: apiKey}
	chatClient, err := NewClientWithKeyCredential(endpoint, cred, "thisdoesntexist", newClientOptionsForTest(t))
	require.NoError(t, err)

	_, err = chatClient.GetEmbeddings(context.Background(), EmbeddingsOptions{}, nil)

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, "DeploymentNotFound", respErr.ErrorCode)
}

func TestClient_GetCompletions(t *testing.T) {
	type args struct {
		ctx          context.Context
		deploymentID string
		body         CompletionsOptions
		options      *ClientGetCompletionsOptions
	}
	cred := KeyCredential{APIKey: apiKey}
	deploymentID := "text-davinci-003"
	client, err := NewClientWithKeyCredential(endpoint, cred, deploymentID, newClientOptionsForTest(t))
	if err != nil {
		log.Fatalf("%v", err)
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    ClientGetCompletionsResponse
		wantErr bool
	}{
		{
			name:   "chatbot",
			client: client,
			args: args{
				ctx:          context.TODO(),
				deploymentID: deploymentID,
				body: CompletionsOptions{
					Prompt:      []*string{to.Ptr("What is Azure OpenAI?")},
					MaxTokens:   to.Ptr(int32(2048 - 127)),
					Temperature: to.Ptr(float32(0.0)),
				},
				options: nil,
			},
			want: ClientGetCompletionsResponse{
				Completions: Completions{
					Choices: []*Choice{
						{
							Text:         to.Ptr("\n\nAzure OpenAI is a platform from Microsoft that provides access to OpenAI's artificial intelligence (AI) technologies. It enables developers to build, train, and deploy AI models in the cloud. Azure OpenAI provides access to OpenAI's powerful AI technologies, such as GPT-3, which can be used to create natural language processing (NLP) applications, computer vision models, and reinforcement learning models."),
							Index:        to.Ptr(int32(0)),
							FinishReason: to.Ptr(CompletionsFinishReason("stop")),
							Logprobs:     nil,
						},
					},
					Usage: &CompletionsUsage{
						CompletionTokens: to.Ptr(int32(85)),
						PromptTokens:     to.Ptr(int32(6)),
						TotalTokens:      to.Ptr(int32(91)),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetCompletions(tt.args.ctx, tt.args.body, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetCompletions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opts := cmpopts.IgnoreFields(Completions{}, "Created", "ID")
			if diff := cmp.Diff(tt.want.Completions, got.Completions, opts); diff != "" {
				t.Errorf("Client.GetCompletions(): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestClient_GetEmbeddings(t *testing.T) {
	type args struct {
		ctx          context.Context
		deploymentID string
		body         EmbeddingsOptions
		options      *ClientGetEmbeddingsOptions
	}
	deploymentID := "embedding"
	cred := KeyCredential{APIKey: apiKey}
	client, err := NewClientWithKeyCredential(endpoint, cred, deploymentID, newClientOptionsForTest(t))
	if err != nil {
		log.Fatalf("%v", err)
	}

	tests := []struct {
		name    string
		client  *Client
		args    args
		want    ClientGetEmbeddingsResponse
		wantErr bool
	}{
		{
			name:   "Embeddings",
			client: client,
			args: args{
				ctx:          context.TODO(),
				deploymentID: "embedding",
				body: EmbeddingsOptions{
					Input: []byte("\"Your text string goes here\""),
					Model: to.Ptr("text-similarity-curie-001"),
				},
				options: nil,
			},
			want: ClientGetEmbeddingsResponse{
				Embeddings{
					Data:  []*EmbeddingItem{},
					Usage: &EmbeddingsUsage{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetEmbeddings(tt.args.ctx, tt.args.body, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetEmbeddings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.Embeddings.Data[0].Embedding) != 4096 {
				t.Errorf("Client.GetEmbeddings() len(Data) want 4096, got %d", len(got.Embeddings.Data))
				return
			}
		})
	}
}
