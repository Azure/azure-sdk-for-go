//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
)

func ExampleClient_CreateImage() {
	azureOpenAIKey := os.Getenv("AOAI_API_KEY")

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_ENDPOINT")

	if azureOpenAIKey == "" || azureOpenAIEndpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential, err := azopenai.NewKeyCredential(azureOpenAIKey)

	if err != nil {
		// TODO: handle error
	}

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		// TODO: handle error
	}

	resp, err := client.CreateImage(context.TODO(), azopenai.ImageGenerationOptions{
		Prompt:         to.Ptr("a cat"),
		ResponseFormat: to.Ptr(azopenai.ImageGenerationResponseFormatURL),
	}, nil)

	if err != nil {
		// TODO: handle error
	}

	for _, generatedImage := range resp.Data {
		// the underlying type for the generatedImage is dictated by the value of
		// ImageGenerationOptions.ResponseFormat. In this example we used `azopenai.ImageGenerationResponseFormatURL`,
		// so the underlying type will be ImageLocation.

		resp, err := http.Head(*generatedImage.URL)

		if err != nil {
			// TODO: handle error
		}

		fmt.Fprintf(os.Stderr, "Image generated, HEAD request on URL returned %d\n", resp.StatusCode)
	}

	// Output:
}
