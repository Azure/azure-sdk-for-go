// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeserializeEmbeddingsArray(t *testing.T) {
	t.Run("IsBase64", func(t *testing.T) {
		rawJSON := "\"eLoFPtiY3b7NSTQ9ob4DPzM/jj1KG7y+sJRfvlov8z76nIy8qPytvg==\""

		var ei EmbeddingItem
		err := deserializeEmbeddingsArray(json.RawMessage(rawJSON), &ei)
		require.NoError(t, err)

		require.Equal(t, rawJSON[1:len(rawJSON)-1], ei.EmbeddingBase64)
	})

	t.Run("InvalidBase64", func(t *testing.T) {
		rawJSON := "\"hello\""

		var ei EmbeddingItem
		err := deserializeEmbeddingsArray(json.RawMessage(rawJSON), &ei)
		require.NoError(t, err)
		require.Equal(t, rawJSON[1:len(rawJSON)-1], ei.EmbeddingBase64)
	})

	t.Run("IsFloats", func(t *testing.T) {
		rawJSON := "[\n        0.13059413,\n        -0.43280673,\n        0.044015694,\n        0.5146275,\n        0.06945648,\n        -0.3673957,\n        -0.21834064,\n        0.47497064,\n        -0.017164696,\n        -0.33981824\n      ]"
		expected := []float32{
			0.13059413,
			-0.43280673,
			0.044015694,
			0.5146275,
			0.06945648,
			-0.3673957,
			-0.21834064,
			0.47497064,
			-0.017164696,
			-0.33981824,
		}

		var ei EmbeddingItem
		err := deserializeEmbeddingsArray(json.RawMessage(rawJSON), &ei)
		require.NoError(t, err)
		require.Equal(t, expected, ei.Embedding)
	})
}
