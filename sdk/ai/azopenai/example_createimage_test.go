// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/openai/openai-go/v3"
)

// Example_createImage demonstrates how to generate images using Azure OpenAI's DALL-E model.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Configure image generation parameters including size and format
// - Generate an image from a text prompt
// - Verify the generated image URL is accessible
//
// The example uses environment variables for configuration:
// - AOAI_DALLE_ENDPOINT: Your Azure OpenAI endpoint URL
// - AOAI_DALLE_MODEL: The deployment name of your DALL-E model
//
// Image generation is useful for:
// - Creating custom illustrations and artwork
// - Generating visual content for applications
// - Prototyping design concepts
// - Producing visual aids for documentation
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
		resp, err := http.Get(generatedImage.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// Handle non-200 status code
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		imageData, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		// Save the generated image to a file
		err = os.WriteFile("generated_image.png", imageData, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}
	}
}
