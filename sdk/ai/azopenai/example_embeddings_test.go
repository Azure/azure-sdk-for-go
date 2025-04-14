// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

func Example_embeddings() {
	model := os.Getenv("AOAI_EMBEDDINGS_MODEL") // eg. "text-embedding-ada-002"

	endpoint := os.Getenv("AOAI_EMBEDDINGS_ENDPOINT")

	if model == "" || endpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-08-01-preview"),
		azure.WithTokenCredential(tokenCredential),
	)

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

	// Output:
}
