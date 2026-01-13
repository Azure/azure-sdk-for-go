// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

// Example_responsesApiTextGeneration demonstrates how to use the Azure OpenAI Responses API for text generation.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Send a simple text prompt
// - Process the response
// - Delete the response to clean up
//
// The example uses environment variables for configuration:
// - AZURE_OPENAI_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
// - AZURE_OPENAI_MODEL: The deployment name of your model (e.g., "gpt-4o")
//
// The Responses API is a new stateful API from Azure OpenAI that brings together capabilities
// from chat completions and assistants APIs in a unified experience.
func Example_responsesApiTextGeneration() {
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	model := os.Getenv("AZURE_OPENAI_MODEL")

	// Create a client with token credentials
	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// Create a simple text input
	resp, err := client.Responses.New(
		context.TODO(),
		responses.ResponseNewParams{
			Model: model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("Define and explain the concept of catastrophic forgetting?"),
			},
		},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	// Process the response
	fmt.Fprintf(os.Stderr, "Response ID: %s\n", resp.ID)
	fmt.Fprintf(os.Stderr, "Model: %s\n", resp.Model)

	// Print the text content from the output
	for _, output := range resp.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					fmt.Fprintf(os.Stderr, "Content: %s\n", content.Text)
				}
			}
		}
	}

	// Delete the response to clean up
	err = client.Responses.Delete(
		context.TODO(),
		resp.ID,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR deleting response: %s\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Response deleted successfully\n")
	}

	fmt.Fprintf(os.Stderr, "Example complete\n")
}

// Example_responsesApiChaining demonstrates how to chain multiple responses together
// in a conversation flow using the Azure OpenAI Responses API.
// This example shows how to:
// - Create an initial response
// - Chain a follow-up response using the previous response ID
// - Process both responses
// - Delete both responses to clean up
//
// The example uses environment variables for configuration:
// - AZURE_OPENAI_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
// - AZURE_OPENAI_MODEL: The deployment name of your model (e.g., "gpt-4o")
func Example_responsesApiChaining() {
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	model := os.Getenv("AZURE_OPENAI_MODEL")

	// Create a client with token credentials
	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// Create the first response
	firstResponse, err := client.Responses.New(
		context.TODO(),
		responses.ResponseNewParams{
			Model: model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("Define and explain the concept of catastrophic forgetting?"),
			},
		},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "First response ID: %s\n", firstResponse.ID)

	// Chain a second response using the previous response ID
	secondResponse, err := client.Responses.New(
		context.TODO(),
		responses.ResponseNewParams{
			Model: model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("Explain this at a level that could be understood by a college freshman"),
			},
			PreviousResponseID: openai.String(firstResponse.ID),
		},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "Second response ID: %s\n", secondResponse.ID)

	// Print the text content from the second response
	for _, output := range secondResponse.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					fmt.Fprintf(os.Stderr, "Second response content: %s\n", content.Text)
				}
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Example complete\n")
}

// Example_responsesApiStreaming demonstrates how to use streaming with the Azure OpenAI Responses API.
// This example shows how to:
// - Create a streaming response
// - Process the stream events as they arrive
// - Clean up by deleting the response
//
// The example uses environment variables for configuration:
// - AZURE_OPENAI_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
// - AZURE_OPENAI_MODEL: The deployment name of your model (e.g., "gpt-4o")
func Example_responsesApiStreaming() {
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	model := os.Getenv("AZURE_OPENAI_MODEL")

	// Create a client with token credentials
	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// Create a streaming response
	stream := client.Responses.NewStreaming(
		context.TODO(),
		responses.ResponseNewParams{
			Model: model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("This is a test"),
			},
		},
	)

	// Process the stream
	fmt.Fprintf(os.Stderr, "Streaming response: ")

	for stream.Next() {
		event := stream.Current()
		if event.Type == "response.output_text.delta" {
			fmt.Fprintf(os.Stderr, "%s", event.Delta)
		}
	}

	if stream.Err() != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: %s\n", stream.Err())
		return
	}

	fmt.Fprintf(os.Stderr, "\nExample complete\n")
}

// Example_responsesApiFunctionCalling demonstrates how to use the Azure OpenAI Responses API with function calling.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Define tools (functions) that the model can call
// - Process the response containing function calls
// - Provide function outputs back to the model
// - Delete the responses to clean up
//
// The example uses environment variables for configuration:
// - AZURE_OPENAI_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
// - AZURE_OPENAI_MODEL: The deployment name of your model (e.g., "gpt-4o")
func Example_responsesApiFunctionCalling() {
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	model := os.Getenv("AZURE_OPENAI_MODEL")

	// Create a client with token credentials
	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// Define the get_weather function parameters as a JSON schema
	paramSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type": "string",
			},
		},
		"required": []string{"location"},
	}

	// Create a response with tools (functions)
	resp, err := client.Responses.New(
		context.TODO(),
		responses.ResponseNewParams{
			Model: model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("What's the weather in San Francisco?"),
			},
			Tools: []responses.ToolUnionParam{
				{
					OfFunction: &responses.FunctionToolParam{
						Name:        "get_weather",
						Description: openai.String("Get the weather for a location"),
						Parameters:  paramSchema,
					},
				},
			},
		},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	// Process the response to find function calls
	var functionCallID string
	var functionName string

	for _, output := range resp.Output {
		if output.Type == "function_call" {
			functionCallID = output.CallID
			functionName = output.Name
			fmt.Fprintf(os.Stderr, "Function call detected: %s\n", functionName)
			fmt.Fprintf(os.Stderr, "Function arguments: %s\n", output.Arguments)
		}
	}

	// If a function call was found, provide the function output back to the model
	if functionCallID != "" {
		// In a real application, you would actually call the function
		// Here we're just simulating a response
		var functionOutput string
		if functionName == "get_weather" {
			functionOutput = `{"temperature": "72 degrees", "condition": "sunny"}`
		}

		// Create a second response, providing the function output
		secondResp, err := client.Responses.New(
			context.TODO(),
			responses.ResponseNewParams{
				Model:              model,
				PreviousResponseID: openai.String(resp.ID),
				Input: responses.ResponseNewParamsInputUnion{
					OfInputItemList: []responses.ResponseInputItemUnionParam{
						{
							OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
								CallID: functionCallID,
								Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
									OfString: openai.String(functionOutput),
								},
							},
						},
					},
				},
			},
		)

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR with second response: %s\n", err)
			return
		}

		// Process the final model response after receiving function output
		for _, output := range secondResp.Output {
			if output.Type == "message" {
				for _, content := range output.Content {
					if content.Type == "output_text" {
						fmt.Fprintf(os.Stderr, "Final response: %s\n", content.Text)
					}
				}
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Example complete\n")
}

// Example_responsesApiImageInput demonstrates how to use the Azure OpenAI Responses API with image input.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Fetch an image from a URL and encode it to Base64
// - Send a query with both text and a Base64-encoded image
// - Process the response
//
// The example uses environment variables for configuration:
// - AZURE_OPENAI_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
// - AZURE_OPENAI_MODEL: The deployment name of your model (e.g., "gpt-4o")
//
// Note: This example fetches and encodes an image from a URL because there is a known issue with image url
// based image input. Currently only base64 encoded images are supported.
func Example_responsesApiImageInput() {
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	model := os.Getenv("AZURE_OPENAI_MODEL")

	// Create a client with token credentials
	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// Image URL to fetch and encode, you can also use a local file path
	imageURL := "https://www.bing.com/th?id=OHR.BradgateFallow_EN-US3932725763_1920x1080.jpg"

	// Fetch the image from the URL and encode it to Base64
	httpClient := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := httpClient.Get(imageURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR fetching image: %s\n", err)
		return
	}

	defer func() {
		if err := httpResp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		}
	}()

	imgBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR reading image: %s\n", err)
		return
	}

	// Encode the image to Base64
	base64Image := base64.StdEncoding.EncodeToString(imgBytes)
	fmt.Fprintf(os.Stderr, "Successfully encoded image from URL\n")

	// Determine content type based on image data or response headers
	contentType := httpResp.Header.Get("Content-Type")
	if contentType == "" {
		// Default to jpeg if we can't determine
		contentType = "image/jpeg"
	}

	// Create the data URL for the image
	dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Image)

	// Create a response with the image input
	resp, err := client.Responses.New(
		context.TODO(),
		responses.ResponseNewParams{
			Model: model,
			Input: responses.ResponseNewParamsInputUnion{
				OfInputItemList: []responses.ResponseInputItemUnionParam{
					{
						OfInputMessage: &responses.ResponseInputItemMessageParam{
							Role: "user",
							Content: []responses.ResponseInputContentUnionParam{
								{
									OfInputText: &responses.ResponseInputTextParam{
										Text: "What can you see in this image?",
									},
								},
								{
									OfInputImage: &responses.ResponseInputImageParam{
										ImageURL: openai.String(dataURL),
									},
								},
							},
						},
					},
				},
			},
		},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	// Print the text content from the output
	for _, output := range resp.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					fmt.Fprintf(os.Stderr, "Model's description of the image: %s\n", content.Text)
				}
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Example complete\n")
}

// Example_responsesApiReasoning demonstrates how to use the Azure OpenAI Responses API with reasoning.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Send a complex problem-solving request that requires reasoning
// - Enable the reasoning parameter to get step-by-step thought process
// - Process the response
//
// The example uses environment variables for configuration:
// - AZURE_OPENAI_ENDPOINT: Your Azure OpenAI endpoint URL (ex: "https://yourservice.openai.azure.com")
// - AZURE_OPENAI_MODEL: The deployment name of your model (e.g., "gpt-4o")
func Example_responsesApiReasoning() {
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	model := os.Getenv("AZURE_OPENAI_MODEL")

	// Create a client with token credentials
	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("%s/openai/v1", endpoint)),
		azure.WithTokenCredential(tokenCredential),
	)

	// Create a response with reasoning enabled
	// This will make the model show its step-by-step reasoning
	resp, err := client.Responses.New(
		context.TODO(),
		responses.ResponseNewParams{
			Model: model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("Solve the following problem step by step: If a train travels at 120 km/h and needs to cover a distance of 450 km, how long will the journey take?"),
			},
			Reasoning: openai.ReasoningParam{
				Effort: openai.ReasoningEffortMedium,
			},
		},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	// Print the text content from the output
	for _, output := range resp.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					fmt.Fprintf(os.Stderr, "\nOutput: %s\n", content.Text)
				}
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Example complete\n")
}
