// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
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

func TestClient_GetEmbeddings_embeddingsFormat(t *testing.T) {
	testFn := func(t *testing.T, tv testVars, dimension int32) {
		client := newTestClient(t, tv.Endpoint)

		arg := azopenai.EmbeddingsOptions{
			Input:          []string{"hello"},
			EncodingFormat: to.Ptr(azopenai.EmbeddingEncodingFormatBase64),
			DeploymentName: &tv.TextEmbedding3Small,
		}

		if dimension > 0 {
			arg.Dimensions = &dimension
		}

		base64Resp, err := client.GetEmbeddings(context.Background(), arg, nil)
		require.NoError(t, err)

		require.NotEmpty(t, base64Resp.Data)
		require.Empty(t, base64Resp.Data[0].Embedding)
		embeddings := deserializeBase64Embeddings(t, base64Resp.Data[0])

		// sanity checks - we deserialized everything and didn't create anything impossible.
		for _, v := range embeddings {
			require.True(t, v <= 1.0 && v >= -1.0)
		}

		arg2 := azopenai.EmbeddingsOptions{
			Input:          []string{"hello"},
			DeploymentName: &tv.TextEmbedding3Small,
		}

		if dimension > 0 {
			arg2.Dimensions = &dimension
		}

		floatResp, err := client.GetEmbeddings(context.Background(), arg2, nil)
		require.NoError(t, err)

		require.NotEmpty(t, floatResp.Data)
		require.NotEmpty(t, floatResp.Data[0].Embedding)

		require.Equal(t, len(floatResp.Data[0].Embedding), len(embeddings))

		// This works "most of the time" but it's non-deterministic since two separate calls don't always
		// produce the exact same data. Leaving it here in case you want to do some rough checks later.
		// require.Equal(t, floatResp.Data[0].Embedding[0:dimension], base64Resp.Data[0].Embedding[0:dimension])
	}

	for _, dim := range []int32{0, 1, 10, 100} {
		t.Run(fmt.Sprintf("AzureOpenAI(dimensions=%d)", dim), func(t *testing.T) {
			testFn(t, azureOpenAI, dim)
		})

		t.Run(fmt.Sprintf("OpenAI(dimensions=%d)", dim), func(t *testing.T) {
			testFn(t, openAI, dim)
		})
	}
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

func deserializeBase64Embeddings(t *testing.T, ei azopenai.EmbeddingItem) []float32 {
	destBytes, err := base64.StdEncoding.DecodeString(ei.EmbeddingBase64)
	require.NoError(t, err)

	floats := make([]float32, len(destBytes)/4) // it's a binary serialization of float32s.
	var reader = bytes.NewReader(destBytes)

	err = binary.Read(reader, binary.LittleEndian, floats)
	require.NoError(t, err)

	return floats
}
