//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleClient_GetImageGenerations() {
	azureOpenAIKey := os.Getenv("AOAI_DALLE_API_KEY")

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_DALLE_ENDPOINT")

	azureDeployment := os.Getenv("AOAI_DALLE_MODEL")

	if azureOpenAIKey == "" || azureOpenAIEndpoint == "" || azureDeployment == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	resp, err := client.GetImageGenerations(context.TODO(), azopenai.ImageGenerationOptions{
		Prompt:         to.Ptr("a cat"),
		ResponseFormat: to.Ptr(azopenai.ImageGenerationResponseFormatURL),
		DeploymentName: &azureDeployment,
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	for _, generatedImage := range resp.Data {
		// the underlying type for the generatedImage is dictated by the value of
		// ImageGenerationOptions.ResponseFormat. In this example we used `azopenai.ImageGenerationResponseFormatURL`,
		// so the underlying type will be ImageLocation.

		resp, err := http.Head(*generatedImage.URL)

		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		_ = resp.Body.Close()
		fmt.Fprintf(os.Stderr, "Image generated, HEAD request on URL returned %d\n", resp.StatusCode)
	}

	// Output:
}
