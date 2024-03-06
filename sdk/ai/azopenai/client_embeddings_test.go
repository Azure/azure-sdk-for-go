// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestClient_GetEmbeddings_InvalidModel(t *testing.T) {
	client := newTestClient(t, azureOpenAI.Endpoint)

	_, err := client.GetEmbeddings(context.Background(), azopenai.EmbeddingsOptions{
		DeploymentName: to.Ptr("thisdoesntexist"),
	}, nil)

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, "DeploymentNotFound", respErr.ErrorCode)
}

func TestClient_OpenAI_GetEmbeddings(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping OpenAI tests when attempting to do quick tests")
	}

	client := newOpenAIClientForTest(t)
	testGetEmbeddings(t, client, openAI.Embeddings)
}

func TestClient_GetEmbeddings(t *testing.T) {
	client := newTestClient(t, azureOpenAI.Endpoint)
	testGetEmbeddings(t, client, azureOpenAI.Embeddings)
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
					Input:          []string{"\"Your text string goes here\""},
					DeploymentName: &modelOrDeploymentID,
				},
				options: nil,
			},
			want: azopenai.GetEmbeddingsResponse{
				azopenai.Embeddings{
					Data:  []azopenai.EmbeddingItem{},
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

			require.NotEmpty(t, got.Embeddings.Data[0].Embedding)
		})
	}
}
