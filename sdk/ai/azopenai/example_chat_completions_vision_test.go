// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/openai/openai-go"
)

func Example_vision() {
	if !CheckRequiredEnvVars("AOAI_VISION_MODEL", "AOAI_VISION_ENDPOINT") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	model := os.Getenv("AOAI_VISION_MODEL") // ex: gpt-4o"
	endpoint := os.Getenv("AOAI_VISION_ENDPOINT")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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

	// Output:
}
