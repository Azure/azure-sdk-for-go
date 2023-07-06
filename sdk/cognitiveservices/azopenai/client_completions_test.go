//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"log"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestClient_GetCompletions(t *testing.T) {
	type args struct {
		ctx          context.Context
		deploymentID string
		body         azopenai.CompletionsOptions
		options      *azopenai.GetCompletionsOptions
	}
	cred, err := azopenai.NewKeyCredential(apiKey)
	require.NoError(t, err)

	client, err := azopenai.NewClientWithKeyCredential(endpoint, cred, completionsModelDeployment, newClientOptionsForTest(t))
	if err != nil {
		log.Fatalf("%v", err)
	}
	tests := []struct {
		name    string
		client  *azopenai.Client
		args    args
		want    azopenai.GetCompletionsResponse
		wantErr bool
	}{
		{
			name:   "chatbot",
			client: client,
			args: args{
				ctx:          context.TODO(),
				deploymentID: completionsModelDeployment,
				body: azopenai.CompletionsOptions{
					Prompt:      []string{"What is Azure OpenAI?"},
					MaxTokens:   to.Ptr(int32(2048 - 127)),
					Temperature: to.Ptr(float32(0.0)),
				},
				options: nil,
			},
			want: azopenai.GetCompletionsResponse{
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
			opts := cmpopts.IgnoreFields(azopenai.Completions{}, "Created", "ID")
			if diff := cmp.Diff(tt.want.Completions, got.Completions, opts); diff != "" {
				t.Errorf("Client.GetCompletions(): -want, +got:\n%s", diff)
			}
		})
	}
}
