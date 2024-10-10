// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"encoding/json"
)

// EmbeddingItem - Representation of a single embeddings relatedness comparison.
type EmbeddingItem struct {
	// List of embeddings value for the input prompt. These represent a measurement of the vector-based relatedness
	// of the provided input when when [EmbeddingEncodingFormatFloat] is specified.
	Embedding []float32

	// EmbeddingBase64 represents the Embeddings when [EmbeddingEncodingFormatBase64] is specified.
	EmbeddingBase64 string

	// REQUIRED; Index of the prompt to which the EmbeddingItem corresponds.
	Index *int32

	// The object type which is always 'embedding'.
	Object string
}

func deserializeEmbeddingsArray(msg json.RawMessage, embeddingItem *EmbeddingItem) error {
	if len(msg) == 0 {
		return nil
	}

	if msg[0] == '"' && len(msg) > 2 && msg[len(msg)-1] == '"' {
		var s = string(msg)
		embeddingItem.EmbeddingBase64 = s[1 : len(s)-1]
		return nil
	}

	return json.Unmarshal(msg, &embeddingItem.Embedding)
}
