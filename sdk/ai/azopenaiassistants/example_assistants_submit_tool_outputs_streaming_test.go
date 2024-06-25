//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func Example_assistantsUsingSubmitToolOutputsStreaming() {
	azureOpenAIKey := os.Getenv("AOAI_ASSISTANTS_KEY")

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_ASSISTANTS_ENDPOINT")

	if azureOpenAIKey == "" || azureOpenAIEndpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

	client, err := azopenaiassistants.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	assistantName := fmt.Sprintf("your-assistant-name-%d", time.Now().UnixNano())

	// This describes a function that OpenAI will "call" when it wants to answer
	// questions about the current weather.

	// First, let's create an assistant.
	createAssistantResp, err := client.CreateAssistant(context.TODO(), azopenaiassistants.CreateAssistantBody{
		Name:           &assistantName,
		DeploymentName: to.Ptr("gpt-4-1106-preview"),
		Instructions:   to.Ptr("You are an AI assistant that answers questions about the weather using functions like get_current_weather."),
		Tools: []azopenaiassistants.ToolDefinitionClassification{
			&azopenaiassistants.FunctionToolDefinition{
				Function: &azopenaiassistants.FunctionDefinition{
					Name: to.Ptr("get_current_weather"),
					Parameters: map[string]any{
						"required": []string{"location"},
						"type":     "object",
						"properties": map[string]any{
							"location": map[string]any{
								"type":        "string",
								"description": "The city and state, e.g. San Francisco, CA",
							},
							"unit": map[string]any{
								"type": "string",
								"enum": []string{"celsius", "fahrenheit"},
							},
						},
					},
				},
			},
		},
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	assistantID := createAssistantResp.ID

	// cleanup the assistant after this example. Remove this if you want to re-use the assistant.
	defer func() {
		_, err := client.DeleteAssistant(context.TODO(), *assistantID, nil)

		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}
	}()

	// 1/3: First we'll run the thread - this run will stop when it needs us to evaluate a function and
	// submit tool outputs.
	var submitToolOutputsAction *azopenaiassistants.SubmitToolOutputsAction
	var threadID string
	var runID string

	{
		resp, err := client.CreateThreadAndRunStream(context.TODO(), azopenaiassistants.CreateAndRunThreadBody{
			AssistantID:    assistantID,
			DeploymentName: &assistantsModel,
			Thread: &azopenaiassistants.CreateThreadBody{
				Messages: []azopenaiassistants.CreateMessageBody{
					{
						Content: to.Ptr("What's the weather like in Boston, MA, in celsius?"),
						Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
					},
				},
			},
		}, nil)

		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		lastThreadRun, err := processToolOutputsStream(resp.Stream)

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		submitToolOutputsAction = lastThreadRun.RequiredAction.(*azopenaiassistants.SubmitToolOutputsAction)
		threadID = *lastThreadRun.ThreadID
		runID = *lastThreadRun.ID
	}

	type weatherFuncArgs struct {
		Location string
		Unit     string
	}

	var toolOutputToBeSubmitted []azopenaiassistants.ToolOutput

	// 2/3: now we unpack the function arguments and call our function, locally
	{
		funcToolCall := submitToolOutputsAction.SubmitToolOutputs.ToolCalls[0].(*azopenaiassistants.RequiredFunctionToolCall)

		var funcArgs *weatherFuncArgs

		err := json.Unmarshal([]byte(*funcToolCall.Function.Arguments), &funcArgs)

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		err = json.Unmarshal([]byte(*funcToolCall.Function.Arguments), &funcArgs)

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		log.Printf("Function we need to call, with args: %s(%q, %q)", *funcToolCall.Function.Name, funcArgs.Location, funcArgs.Unit)

		// TODO: take the parsed arguments and call into your own function implementation
		// For this example we'll just return a hardcoded answer.
		toolOutputToBeSubmitted = []azopenaiassistants.ToolOutput{
			{
				Output:     to.Ptr("26C"),
				ToolCallID: funcToolCall.ID,
			},
		}
	}

	// 3/3: now we'll submit the outputs and continue streaming results.
	{
		resp, err := client.SubmitToolOutputsToRunStream(context.TODO(), threadID, runID, azopenaiassistants.SubmitToolOutputsToRunBody{
			ToolOutputs: toolOutputToBeSubmitted,
		}, nil)

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		if _, err = processToolOutputsStream(resp.Stream); err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}
	}

	// Output:
}

// processToolOutputsStream continually processes events from the stream.
// If action is required this function returns the relevant ThreadRun data.
func processToolOutputsStream(stream *azopenaiassistants.EventReader[azopenaiassistants.StreamEvent]) (*azopenaiassistants.ThreadRun, error) {
	defer stream.Close()

	var threadRunRequiresAction *azopenaiassistants.ThreadRun

	// process the streaming responses
	for {
		event, err := stream.Read()

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Stream has ended normally")
				return threadRunRequiresAction, nil
			}

			return threadRunRequiresAction, err
		}

		// NOTE: for this example we're handling a small subset of the events that are
		// streaming in. See [azopenaiassistants.AssistantStreamEvent] for the full list of available events.
		switch event.Reason {
		// Assistant events
		case azopenaiassistants.AssistantStreamEventThreadCreated:
			log.Printf("(%s): %s", event.Reason, *event.Event.(*azopenaiassistants.AssistantThread).ID)

		// Thread events
		case azopenaiassistants.AssistantStreamEventThreadRunCreated:
			log.Printf("(%s): %s", event.Reason, *event.Event.(*azopenaiassistants.ThreadRun).ID)
		case azopenaiassistants.AssistantStreamEventThreadRunInProgress:
			log.Printf("(%s): %s", event.Reason, *event.Event.(*azopenaiassistants.ThreadRun).ID)
		case azopenaiassistants.AssistantStreamEventThreadRunCompleted:
			log.Printf("(%s): %s", event.Reason, *event.Event.(*azopenaiassistants.ThreadRun).ID)
		case azopenaiassistants.AssistantStreamEventThreadRunFailed: // failure
			threadRun := event.Event.(*azopenaiassistants.ThreadRun)

			code := ""
			message := ""

			if threadRun.LastError != nil {
				if threadRun.LastError.Code != nil {
					code = *threadRun.LastError.Code
				}

				if threadRun.LastError.Message != nil {
					message = *threadRun.LastError.Message
				}
			}

			return nil, fmt.Errorf("(%s): code: %s, message: %s", event.Reason, code, message)

		//
		// The assistant needs you to call a function so it can continue processing.
		//
		// We need to preserve this so we can submit the tool outputs using either
		// SubmitToolOutputsToRunStream or SubmitToolOutputsToRun.
		///
		case azopenaiassistants.AssistantStreamEventThreadRunRequiresAction:
			log.Printf("(%s): %s", event.Reason, *event.Event.(*azopenaiassistants.ThreadRun).ID)
			threadRunRequiresAction = event.Event.(*azopenaiassistants.ThreadRun)

		// Message events
		case azopenaiassistants.AssistantStreamEventThreadMessageCreated:
			log.Printf("(%s): %s", event.Reason, *event.Event.(*azopenaiassistants.ThreadMessage).ID)
		case azopenaiassistants.AssistantStreamEventThreadMessageDelta:
			messageChunk := event.Event.(*azopenaiassistants.MessageDeltaChunk)
			log.Printf("(%s): %s", event.Reason, *event.Event.(*azopenaiassistants.MessageDeltaChunk).ID)

			for _, c := range messageChunk.Delta.Content {
				switch actualContent := c.(type) {
				case *azopenaiassistants.MessageDeltaImageFileContent:
					log.Printf("  Image: %#v", *actualContent.ImageFile)
				case *azopenaiassistants.MessageDeltaTextContentObject:
					log.Printf("  %q", *actualContent.Text.Value)
				}
			}
		case azopenaiassistants.AssistantStreamEventThreadMessageCompleted:
			log.Printf("(%s): %s", event.Reason, *event.Event.(*azopenaiassistants.ThreadMessage).ID)
		}
	}
}
