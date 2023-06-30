// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/stretchr/testify/require"
)

func TestClient_GetEmbeddings_InvalidModel(t *testing.T) {
	cred, err := azopenai.NewKeyCredential(apiKey)
	require.NoError(t, err)

	chatClient, err := azopenai.NewClientWithKeyCredential(endpoint, cred, "thisdoesntexist", newClientOptionsForTest(t))
	require.NoError(t, err)

	_, err = chatClient.GetEmbeddings(context.Background(), azopenai.EmbeddingsOptions{}, nil)

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, "DeploymentNotFound", respErr.ErrorCode)
}

func TestClient_OpenAI_GetEmbeddings(t *testing.T) {
	client := newOpenAIClientForTest(t)
	modelID := "text-similarity-curie-001"
	testGetEmbeddings(t, client, modelID)
}

func TestClient_GetEmbeddings(t *testing.T) {
	// model deployment points to `text-similarity-curie-001`
	deploymentID := "embedding"

	cred, err := azopenai.NewKeyCredential(apiKey)
	require.NoError(t, err)

	client, err := azopenai.NewClientWithKeyCredential(endpoint, cred, deploymentID, newClientOptionsForTest(t))
	require.NoError(t, err)

	testGetEmbeddings(t, client, deploymentID)
}

func testGetEmbeddings(t *testing.T, client *azopenai.Client, modelOrDeploymentID string) {
	type args struct {
		ctx          context.Context
		deploymentID string
		body         azopenai.EmbeddingsOptions
		options      *azopenai.GetEmbeddingsOptions
	}

	tests := []struct {
		name    string
		client  *azopenai.Client
		args    args
		want    azopenai.GetEmbeddingsResponse
		wantErr bool
	}{
		{
			name:   "Embeddings",
			client: client,
			args: args{
				ctx:          context.TODO(),
				deploymentID: modelOrDeploymentID,
				body: azopenai.EmbeddingsOptions{
					Input: []*string{to.Ptr("\"Your text string goes here\"")},
					Model: &modelOrDeploymentID,
				},
				options: nil,
			},
			want: azopenai.GetEmbeddingsResponse{
				azopenai.Embeddings{
					Data:  []*azopenai.EmbeddingItem{},
					Usage: &azopenai.EmbeddingsUsage{},
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
