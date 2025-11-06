// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
)

// Example_embeddings demonstrates how to generate text embeddings using Azure OpenAI's embedding models.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Convert text input into numerical vector representations
// - Process the embedding vectors from the response
// - Handle embedding results for semantic analysis
//
// The example uses environment variables for configuration:
// - AOAI_EMBEDDINGS_MODEL: The deployment name of your embedding model (e.g., text-embedding-ada-002)
// - AOAI_EMBEDDINGS_ENDPOINT: Your Azure OpenAI endpoint URL
//
// Text embeddings are useful for:
// - Semantic search and information retrieval
// - Text classification and clustering
// - Content recommendation systems
// - Document similarity analysis
// - Natural language understanding tasks
func Example_embeddings() {
	if !CheckRequiredEnvVars("AOAI_EMBEDDINGS_MODEL", "AOAI_EMBEDDINGS_ENDPOINT") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	model := os.Getenv("AOAI_EMBEDDINGS_MODEL") // eg. "text-embedding-ada-002"
	endpoint := os.Getenv("AOAI_EMBEDDINGS_ENDPOINT")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	// Call the embeddings API
	resp, err := client.Embeddings.New(context.TODO(), openai.EmbeddingNewParams{
		Model: openai.EmbeddingModel(model),
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String("The food was delicious and the waiter..."),
		},
	})

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	for i, embed := range resp.Data {
		// embed.Embedding contains the embeddings for this input index
		fmt.Fprintf(os.Stderr, "Got embeddings for input %d with embedding length: %d\n", i, len(embed.Embedding))
	}
}
