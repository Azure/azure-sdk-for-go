// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
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
// - AZURE_OPENAI_API_VERSION: Azure OpenAI service API version to use. See https://learn.microsoft.com/azure/ai-foundry/openai/api-version-lifecycle?tabs=go for information about API versions.
//
// Image generation is useful for:
// - Creating custom illustrations and artwork
// - Generating visual content for applications
// - Prototyping design concepts
// - Producing visual aids for documentation
func Example_createImage() {
	endpoint := os.Getenv("AOAI_DALLE_ENDPOINT")
	model := os.Getenv("AOAI_DALLE_MODEL")
	apiVersion := os.Getenv("AZURE_OPENAI_API_VERSION")

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, apiVersion),
		azure.WithTokenCredential(tokenCredential),
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
		resp, err := http.Get(generatedImage.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			}
		}()

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
