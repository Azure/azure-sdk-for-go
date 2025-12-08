// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/require"
)

func TestClient_GetEmbeddings_InvalidModel(t *testing.T) {
	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.Embeddings.Endpoint)

	_, err := client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
		Model: openai.EmbeddingModel("thisdoesntexist"),
	})

	var openaiErr *openai.Error
	require.ErrorAs(t, err, &openaiErr)
	require.Equal(t, http.StatusBadRequest, openaiErr.StatusCode)
}

func TestClient_GetEmbeddings(t *testing.T) {
	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.Embeddings.Endpoint)

	resp, err := client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: []string{"\"Your text string goes here\""},
		},
		Model: openai.EmbeddingModel(azureOpenAI.Embeddings.Model),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Data[0].Embedding)
}

func TestClient_GetEmbeddings_embeddingsFormat(t *testing.T) {
	testFn := func(t *testing.T, epm endpointWithModel, dimension int64) {
		client := newStainlessTestClientWithAzureURL(t, epm.Endpoint)

		arg := openai.EmbeddingNewParams{
			Input: openai.EmbeddingNewParamsInputUnion{
				OfArrayOfStrings: []string{"hello"},
			},
			EncodingFormat: openai.EmbeddingNewParamsEncodingFormatBase64,
			Model:          openai.EmbeddingModel(epm.Model),
		}

		if dimension > 0 {
			arg.Dimensions = openai.Int(dimension)
		}

		base64Resp, err := client.Embeddings.New(context.Background(), arg)
		require.NoError(t, err)

		require.NotEmpty(t, base64Resp.Data)
		require.Empty(t, base64Resp.Data[0].Embedding)

		embeddings := deserializeBase64Embeddings(t, base64Resp.Data[0].JSON.Embedding.Raw())

		// sanity checks - we deserialized everything and didn't create anything impossible.
		for _, v := range embeddings {
			require.True(t, v <= 1.0 && v >= -1.0)
		}

		arg2 := openai.EmbeddingNewParams{
			Input: openai.EmbeddingNewParamsInputUnion{
				OfArrayOfStrings: []string{"hello"},
			},
			Model: openai.EmbeddingModel(epm.Model),
		}

		if dimension > 0 {
			arg2.Dimensions = openai.Int(dimension)
		}

		floatResp, err := client.Embeddings.New(context.Background(), arg2)
		require.NoError(t, err)

		require.NotEmpty(t, floatResp.Data)
		require.NotEmpty(t, floatResp.Data[0].Embedding)

		require.Equal(t, len(floatResp.Data[0].Embedding), len(embeddings))

		// This works "most of the time" but it's non-deterministic since two separate calls don't always
		// produce the exact same data. Leaving it here in case you want to do some rough checks later.
		// require.Equal(t, floatResp.Data[0].Embedding[0:dimension], base64Resp.Data[0].Embedding[0:dimension])
	}

	for _, dim := range []int64{0, 1, 10, 100} {
		t.Run(fmt.Sprintf("AzureOpenAI(dimensions=%d)", dim), func(t *testing.T) {
			testFn(t, azureOpenAI.TextEmbedding3Small, dim)
		})
	}
}

func deserializeBase64Embeddings(t *testing.T, rawJSON string) []float32 {
	var base64Text *string

	err := json.Unmarshal([]byte(rawJSON), &base64Text)
	require.NoError(t, err)

	destBytes, err := base64.StdEncoding.DecodeString(*base64Text)
	require.NoError(t, err)

	floats := make([]float32, len(destBytes)/4) // it's a binary serialization of float32s.
	var reader = bytes.NewReader(destBytes)

	err = binary.Read(reader, binary.LittleEndian, floats)
	require.NoError(t, err)

	return floats
}
