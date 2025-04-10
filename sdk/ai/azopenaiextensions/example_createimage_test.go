// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

func Example_createImage() {
	endpoint := os.Getenv("AOAI_DALLE_ENDPOINT")
	model := os.Getenv("AOAI_DALLE_MODEL")

	if endpoint == "" || model == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	// Create a token credential using Azure Identity
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	// Initialize OpenAI client with Azure configurations using token credential
	client := openai.NewClient(
		azure.WithTokenCredential(cred),
		azure.WithEndpoint(endpoint, "2024-12-01-preview"),
	)

	resp, err := client.Images.Generate(context.TODO(), openai.ImageGenerateParams{
		Prompt:         "a cat",
		Model:          openai.ImageModel(model),
		ResponseFormat: openai.ImageGenerateParamsResponseFormatURL,
		Size:           openai.ImageGenerateParamsSize1024x1024,
	})

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	for _, generatedImage := range resp.Data {
		resp, err := http.Head(generatedImage.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		_ = resp.Body.Close()
		fmt.Fprintf(os.Stderr, "Image generated, HEAD request on URL returned %d\n", resp.StatusCode)
	}
}
