// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
)

// Example_vision demonstrates how to use Azure OpenAI's Vision capabilities for image analysis.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Send an image URL to the model for analysis
// - Configure the chat completion request with image content
// - Process the model's description of the image
//
// The example uses environment variables for configuration:
// - AOAI_VISION_MODEL: The deployment name of your vision-capable model (e.g., gpt-4-vision)
// - AOAI_VISION_ENDPOINT: Your Azure OpenAI endpoint URL
// - AZURE_OPENAI_API_VERSION: Azure OpenAI service API version to use. See https://learn.microsoft.com/azure/ai-foundry/openai/api-version-lifecycle?tabs=go for information about API versions.
//
// Vision capabilities are useful for:
// - Image description and analysis
// - Visual question answering
// - Content moderation
// - Accessibility features
// - Image-based search and retrieval
func Example_vision() {
	model := os.Getenv("AOAI_VISION_MODEL") // ex: gpt-4o"
	endpoint := os.Getenv("AOAI_VISION_ENDPOINT")
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

	imageURL := "https://www.bing.com/th?id=OHR.BradgateFallow_EN-US3932725763_1920x1080.jpg"

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfArrayOfContentParts: []openai.ChatCompletionContentPartUnionParam{
							{
								OfText: &openai.ChatCompletionContentPartTextParam{
									Text: "Describe this image",
								},
							},
							{
								OfImageURL: &openai.ChatCompletionContentPartImageParam{
									ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
										URL: imageURL,
									},
								},
							},
						},
					},
				},
			},
		},
		MaxTokens: openai.Int(512),
	})

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		// Prints "Result: The image shows two deer standing in a field of tall, autumn-colored ferns"
		fmt.Fprintf(os.Stderr, "Result: %s\n", resp.Choices[0].Message.Content)
	}
}
