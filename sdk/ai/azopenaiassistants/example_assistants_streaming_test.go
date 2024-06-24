//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
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

func Example_assistantsStreaming() {
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

	// First, let's create an assistant.
	createAssistantResp, err := client.CreateAssistant(context.Background(), azopenaiassistants.CreateAssistantBody{
		Name:           &assistantName,
		DeploymentName: to.Ptr("gpt-4-1106-preview"),
		Instructions:   to.Ptr("You are an AI assistant that can write code to help answer math questions."),
		Tools: []azopenaiassistants.ToolDefinitionClassification{
			&azopenaiassistants.CodeInterpreterToolDefinition{},
			// others...
			// &azopenaiassistants.FunctionToolDefinition{}
			// &azopenaiassistants.RetrievalToolDefinition{}
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

	resp, err := client.CreateThreadAndRunStream(context.Background(), azopenaiassistants.CreateAndRunThreadBody{
		AssistantID:    assistantID,
		DeploymentName: &assistantsModel,
		Instructions:   to.Ptr("You're a helpful assistant that provides answers on questions about beetles."),
		Thread: &azopenaiassistants.CreateThreadBody{
			Messages: []azopenaiassistants.CreateMessageBody{
				{
					Content: to.Ptr("Tell me about the origins of the humble beetle in two sentences."),
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

	// process the streaming responses
	for {
		event, err := resp.Stream.Read()

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Stream has ended normally")
				break
			}

			// TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}

		// NOTE: for this example we're handling a small subset of the events that are
		// streaming in. See [AssistantStreamEvent] for the full list of available events.
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
			log.Printf("(%s): %#v", event.Reason, *threadRun.LastError)

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
