// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleClient_GetChatCompletions_vision() {
	azureOpenAIKey := os.Getenv("AOAI_VISION_API_KEY")
	modelDeployment := os.Getenv("AOAI_VISION_MODEL") // ex: gpt-4-vision-preview"

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_VISION_ENDPOINT")

	if azureOpenAIKey == "" || modelDeployment == "" || azureOpenAIEndpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	imageURL := "https://www.bing.com/th?id=OHR.BradgateFallow_EN-US3932725763_1920x1080.jpg"

	content := azopenai.NewChatRequestUserMessageContent([]azopenai.ChatCompletionRequestMessageContentPartClassification{
		&azopenai.ChatCompletionRequestMessageContentPartText{
			Text: to.Ptr("Describe this image"),
		},
		&azopenai.ChatCompletionRequestMessageContentPartImage{
			ImageURL: &azopenai.ChatCompletionRequestMessageContentPartImageURL{
				URL: &imageURL,
			},
		},
	})

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	resp, err := client.GetChatCompletions(ctx, azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestUserMessage{
				Content: content,
			},
		},
		MaxTokens:      to.Ptr[int32](512),
		DeploymentName: to.Ptr(modelDeployment),
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	for _, choice := range resp.Choices {
		if choice.Message != nil && choice.Message.Content != nil {
			// Prints "Result: The image shows two deer standing in a field of tall, autumn-colored ferns"
			fmt.Fprintf(os.Stderr, "Result: %s\n", *choice.Message.Content)
		}
	}

	// Output:
}
