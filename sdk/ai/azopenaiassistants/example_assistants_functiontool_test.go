//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func Example_assistantsUsingFunctionTool() {
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
		log.Fatalf("ERROR: %s", err)
	}

	assistantName := fmt.Sprintf("your-assistant-name-%d", time.Now().UnixNano())

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	getWeatherFunctionTool := &azopenaiassistants.FunctionToolDefinition{
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
	}

	// First, let's create an assistant.
	createAssistantResp, err := client.CreateAssistant(context.Background(), azopenaiassistants.CreateAssistantBody{
		Name:           &assistantName,
		DeploymentName: to.Ptr("gpt-4-1106-preview"),
		Instructions:   to.Ptr("You are a personal math tutor. Write and run code to answer math questions."),
		Tools: []azopenaiassistants.ToolDefinitionClassification{
			// Defines a function with a signature like this:
			// get_current_weather(location string, unit string)
			getWeatherFunctionTool,
		},
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// create our thread and a run for our question "What's the weather like in Boston, MA, in celsius?"
	question := "What's the weather like in Boston, MA, in celsius?"
	fmt.Fprintf(os.Stderr, "Asking the question: '%s'\n", question)

	var lastMessageID, runID, threadID string
	{
		fmt.Fprintf(os.Stderr, "Creating our thread\n")
		createThreadResp, err := client.CreateThread(context.Background(), azopenaiassistants.CreateThreadBody{}, nil)

		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		threadID = *createThreadResp.ID

		msgResponse, err := client.CreateMessage(context.Background(), threadID, azopenaiassistants.CreateMessageBody{
			Content: &question,
			Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
		}, nil)

		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		lastMessageID = *msgResponse.ID

		// run the thread
		fmt.Fprintf(os.Stderr, "Creating our run for thread %s\n", *createThreadResp.ID)

		createRunResp, err := client.CreateRun(context.Background(), *createThreadResp.ID, azopenaiassistants.CreateRunBody{
			AssistantID:  createAssistantResp.ID,
			Instructions: to.Ptr("Use functions to answer questions, when possible."),
		}, nil)

		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		runID = *createRunResp.ID
	}

	fmt.Fprintf(os.Stderr, "Waiting for the Run status to indicate it needs tool outputs\n")
	// NOTE: see https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants#example-package-AssistantsConversationLoop
	// for the pollUntilRunEnds function.
	lastResp, err := pollUntilRunEnds(context.Background(), client, threadID, runID)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Our question that we asked requires a tool to answer.
	// *lastResp.Status == azopenaiassistants.RunStatusRequiresAction
	fmt.Fprintf(os.Stderr, "Got run status %s\n", *lastResp.Status)
	fmt.Fprintf(os.Stderr, "Check the response for information we need to submit tool outputs\n")

	var weatherFuncArgs *WeatherFuncArgs
	var funcToolCall *azopenaiassistants.RequiredFunctionToolCall
	{
		// the arguments we need to run the next thing are inside of the run
		submitToolsAction, ok := lastResp.RequiredAction.(*azopenaiassistants.SubmitToolOutputsAction)

		if !ok {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: did not get an azopenaiassistants.SubmitToolOutputsAction as our required action")
		}

		tmpFuncToolCall, ok := submitToolsAction.SubmitToolOutputs.ToolCalls[0].(*azopenaiassistants.RequiredFunctionToolCall)

		if !ok || *tmpFuncToolCall.Function.Name != "get_current_weather" {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: did not get an azopenaiassistants.RequiredFunctionToolCall as our required action, or got an incorrect function name")
		}

		if err := json.Unmarshal([]byte(*tmpFuncToolCall.Function.Arguments), &weatherFuncArgs); err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		funcToolCall = tmpFuncToolCall
	}

	// now call the function (OpenAI is providing the arguments)
	fmt.Fprintf(os.Stderr, "Call our own function to get the weather for %s, in %s\n", weatherFuncArgs.Location, weatherFuncArgs.Unit)
	{
		// submit our outputs from evaluating the tool
		_, err := client.SubmitToolOutputsToRun(context.Background(), threadID, runID, azopenaiassistants.SubmitToolOutputsToRunBody{
			ToolOutputs: []azopenaiassistants.ToolOutput{
				{
					Output:     to.Ptr("0C"),
					ToolCallID: funcToolCall.ID,
				},
			},
		}, nil)

		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}
	}

	// the run will restart now, we just need to wait until it finishes
	// NOTE: see https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants#example-package-AssistantsConversationLoop
	// for the pollUntilRunEnds function.
	lastResp, err = pollUntilRunEnds(context.Background(), client, threadID, runID)

	if err != nil || *lastResp.Status != azopenaiassistants.RunStatusCompleted {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Fprintf(os.Stderr, "Get responses from the assistant, based on our tool outputs\n")

	// NOTE: see https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants#example-package-AssistantsConversationLoop
	// for the getLatestMessages function.
	latestMessages, err := getLatestMessages(context.Background(), client, threadID, &lastMessageID)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	if len(latestMessages) > 0 {
		// Prints something like:
		// [ASSISTANT] <id>: Text response: The current weather in Boston, MA, measured in Celsius, is 0°C.
		// NOTE: see https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants#example-package-AssistantsConversationLoop
		// for the printAssistantMessages function.
		err = printAssistantMessages(context.Background(), client, latestMessages)

		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}
	}

	// DisabledOutput:
}

type WeatherFuncArgs struct {
	Location string
	Unit     string
}
