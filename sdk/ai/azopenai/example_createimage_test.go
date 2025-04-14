// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/openai/openai-go"
)

func Example_createImage() {
	if !CheckRequiredEnvVars("AOAI_DALLE_ENDPOINT", "AOAI_DALLE_MODEL") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	endpoint := os.Getenv("AOAI_DALLE_ENDPOINT")
	model := os.Getenv("AOAI_DALLE_MODEL")

	// Initialize OpenAI client with Azure configurations using token credential
	client, err := CreateOpenAIClientWithToken(endpoint, "2024-12-01-preview")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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
