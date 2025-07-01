//go:build go1.21
// +build go1.21

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
	"github.com/stretchr/testify/require"
)

func TestClient_ResponsesTextGeneration(t *testing.T) {
	client := newStainlessTestClient(t, azureOpenAI.Assistants.Endpoint)
	model := azureOpenAI.Assistants.Model

	resp, err := client.Responses.New(
		context.TODO(),
		responses.ResponseNewParams{
			Model: model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("Define and explain the concept of catastrophic forgetting?"),
			},
		},
	)
	customRequireNoError(t, err)

	require.Equal(t, model, resp.Model)

	// Verify there's some text content in the output
	var hasTextContent bool
	for _, output := range resp.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					hasTextContent = true
					require.NotEmpty(t, content.Text)
				}
			}
		}
	}
	require.True(t, hasTextContent, "Response should contain text content output with message type and output_text content type")
}

func TestClient_ResponsesChaining(t *testing.T) {
	client := newStainlessTestClient(t, azureOpenAI.Assistants.Endpoint)

	// Disable the sanitizer for the response ID to allow chaining
	err := recording.RemoveRegisteredSanitizers([]string{"AZSDK3430"}, getRecordingOptions(t))
	if err != nil {
		t.Fatalf("Failed to remove registered sanitizers: %v", err)
	}

	model := azureOpenAI.Assistants.Model

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
	customRequireNoError(t, err)
	require.NotEmpty(t, firstResponse.ID)

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
	customRequireNoError(t, err)

	// Verify there's some text content in the second response
	var hasTextContent bool
	for _, output := range secondResponse.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					hasTextContent = true
					require.NotEmpty(t, content.Text)
				}
			}
		}
	}
	require.True(t, hasTextContent, "Second response should contain text content")
}

func TestClient_ResponsesStreaming(t *testing.T) {
	client := newStainlessTestClient(t, azureOpenAI.Assistants.Endpoint)
	model := azureOpenAI.Assistants.Model

	stream := client.Responses.NewStreaming(
		context.TODO(),
		responses.ResponseNewParams{
			Model: model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("Write a brief description of artificial intelligence"),
			},
		},
	)

	var combinedOutput string

	for stream.Next() {
		event := stream.Current()
		if event.Type == "response.output_text.delta" {
			combinedOutput += event.Delta.OfString
		}
	}

	require.NoError(t, stream.Err())
	require.NotEmpty(t, combinedOutput)

	// Close the stream and verify there is no error on closing
	err := stream.Close()
	require.NoError(t, err, "Stream close should not produce an error")
}

func TestClient_ResponsesFunctionCalling(t *testing.T) {
	client := newStainlessTestClient(t, azureOpenAI.Assistants.Endpoint)
	model := azureOpenAI.Assistants.Model

	// Disable the sanitizer for the response ID to allow chaining
	err := recording.RemoveRegisteredSanitizers([]string{"AZSDK3430"}, getRecordingOptions(t))
	if err != nil {
		t.Fatalf("Failed to remove registered sanitizers: %v", err)
	}

	// Disable the sanitizer for the function name
	err = recording.RemoveRegisteredSanitizers([]string{"AZSDK3493"}, getRecordingOptions(t))
	if err != nil {
		t.Fatalf("Failed to remove registered sanitizers: %v", err)
	}

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
	customRequireNoError(t, err)
	require.NotEmpty(t, resp.ID)

	// Process the response to find function calls
	var functionCallID string
	var functionName string
	var functionArgs string

	for _, output := range resp.Output {
		if output.Type == "function_call" {
			functionCallID = output.CallID
			functionName = output.Name
			functionArgs = output.Arguments
			break
		}
	}

	// Check if the function call was detected
	require.NotEmpty(t, functionCallID, "Function call ID should not be empty")
	require.Contains(t, functionArgs, "San Francisco", "Arguments should contain San Francisco")

	require.Equal(t, "get_weather", functionName, "Function name should be get_weather")

	// If a function call was found, provide the function output back to the model
	functionOutput := `{"temperature": "72 degrees", "condition": "sunny"}`
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
							Output: functionOutput,
						},
					},
				},
			},
		},
	)
	customRequireNoError(t, err)

	// Check if there's a final text response that uses the function output
	var finalResponse string
	for _, output := range secondResp.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					finalResponse = content.Text
					break
				}
			}
		}
	}

	require.NotEmpty(t, finalResponse, "Final response should not be empty")
	require.Contains(t, finalResponse, "72 degrees", "Final response should include function output")
}

func TestClient_ResponsesImageInput(t *testing.T) {
	client := newStainlessTestClient(t, azureOpenAI.Assistants.Endpoint)
	model := azureOpenAI.Assistants.Model

	// Load the sample image file of two deer
	imageBytes, err := os.ReadFile("testdata/sampleimage_two_deers.jpg")
	require.NoError(t, err)

	// Create a base64 encoded data URL for the image
	encodedImage := base64.StdEncoding.EncodeToString(imageBytes)
	dataURL := fmt.Sprintf("data:image/jpeg;base64,%s", encodedImage)

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
										Text: "What can you see in this image? Describe it briefly.",
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

	customRequireNoError(t, err)

	// Check the response for image description
	var imageDescription string
	for _, output := range resp.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					imageDescription = content.Text
					break
				}
			}
		}
	}

	require.NotEmpty(t, imageDescription, "Image description should not be empty")
}

func TestClient_ResponsesReasoning(t *testing.T) {
	client := newStainlessTestClient(t, azureOpenAI.Reasoning.Endpoint)
	model := azureOpenAI.Reasoning.Model

	// Create a response with reasoning enabled
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
	customRequireNoError(t, err)

	// Check the response for reasoning steps
	var solution string
	for _, output := range resp.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" {
					solution = content.Text
					break
				}
			}
		}
	}

	require.NotEmpty(t, solution, "Solution should not be empty")
}
