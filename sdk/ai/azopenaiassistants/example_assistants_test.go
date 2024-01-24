//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func Example_assistants() {
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

	// First, let's create an assistant.
	createAssistantResp, err := client.CreateAssistant(context.Background(), azopenaiassistants.AssistantCreationBody{
		Name:           &assistantName,
		DeploymentName: to.Ptr("gpt-4-1106-preview"),
		Instructions:   to.Ptr("You are a personal math tutor. Write and run code to answer math questions."),
		//FileIDs:        []string{*uploadedPythonFile.ID},
		Tools: []azopenaiassistants.ToolDefinitionClassification{
			// &azopenaiassistants.CodeInterpreterToolDefinition{},
			// others...
			// &azopenaiassistants.FunctionToolDefinition{}
			// &azopenaiassistants.RetrievalToolDefinition{}
		},
	}, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	assistantID := createAssistantResp.ID

	// cleanup the assistant after this example. Remove this if you want to re-use the assistant.
	defer func() {
		_, err := client.DeleteAssistant(context.TODO(), *assistantID, nil)

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}
	}()

	// Now we'll create a thread. The thread is where you will add messages, which can later
	// be evaluated using a Run. A thread can be re-used by multiple Runs.
	createThreadResp, err := client.CreateThread(context.Background(), azopenaiassistants.AssistantThreadCreationOptions{}, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	threadID := createThreadResp.ID

	assistantCtx, stopAssistant := context.WithCancel(context.TODO())

	callIdx := -1

	// This is just a simplified example of how you could handle a conversation - `assistantMessages` are the messages that
	// are responses from the assistant, and you return messages from here that are then added to the conversation.
	handleConversation := func(assistantMessages []azopenaiassistants.ThreadMessage) []azopenaiassistants.CreateMessageBody {
		callIdx++

		printAssistantMessages(assistantMessages)

		// For this example we'll just synthesize some responses so we can get a conversation and eventually just end the
		// conversation.
		switch callIdx {
		case 0:
			text := "Can you help me find the y intercept for y = x +4?"
			fmt.Fprintf(os.Stderr, "[ME] %s\n", text)

			return []azopenaiassistants.CreateMessageBody{
				{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: &text},
			}
		case 1:
			text := "Can you explain it with a Python program?"
			fmt.Fprintf(os.Stderr, "[ME] %s\n", text)

			return []azopenaiassistants.CreateMessageBody{
				{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: &text},
			}
		case 2:
			text := "Can you give me the result if that Python program had 'x' set to 10"
			fmt.Fprintf(os.Stderr, "[ME] %s\n", text)

			return []azopenaiassistants.CreateMessageBody{
				{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: &text},
			}
		default:
			stopAssistant()
		}
		return nil
	}

	if err = assistantLoop(assistantCtx, client, *assistantID, *threadID, handleConversation); err != nil {
		// if this is a cancellation error it's just us trying to stop the assistant loop.
		if errors.Is(err, context.Canceled) {
			fmt.Fprintf(os.Stderr, "Assistant stopped cleanly")
		} else {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}
	}

	// Output:
}

// conversationHandler takes responses from an assistant and returns our reply messages. Returns the responses
// based on the contents of assistantMessages
// - assistantMessages - messages that have arrived since our last read of the thread.
type conversationHandler func(assistantMessages []azopenaiassistants.ThreadMessage) []azopenaiassistants.CreateMessageBody

func assistantLoop(ctx context.Context, client *azopenaiassistants.Client,
	assistantID string, threadID string,
	handleConversation conversationHandler) error {
	// from here we'll run in a loop, adding new messages to the conversation and reading the assistants
	// responses.

	var lastAssistantResponses []azopenaiassistants.ThreadMessage

	for {
		yourResponses := handleConversation(lastAssistantResponses)

		var lastMessageID *string

		for _, yourResponse := range yourResponses {
			// Add some messages to the thread. We will use Run the thread later to evaluate these and to get
			// responses from the assistant.
			createMessageResp, err := client.CreateMessage(context.Background(), threadID, yourResponse, nil)

			if err != nil {
				return err
			}

			// we'll always track the final message ID in the thread - when we pull responses we can be more efficient
			// and only grab what's new.
			lastMessageID = createMessageResp.ID
		}

		createRunResp, err := client.CreateRun(context.Background(), threadID, azopenaiassistants.CreateRunBody{
			AssistantID: &assistantID,
		}, nil)

		if err != nil {
			return err
		}

		runID := *createRunResp.ID

		if err := pollRunEnd(ctx, client, threadID, runID); err != nil {
			return err
		}

		log.Printf("====> %s\n", *lastMessageID)

		lastAssistantResponses = nil

		// get all the messages that were added after our most recently added message.
		listMessagesPager := client.NewListMessagesPager(threadID, &azopenaiassistants.ListMessagesOptions{
			After: lastMessageID,
			Order: to.Ptr(azopenaiassistants.ListSortOrderAscending),
		})

		for listMessagesPager.More() {
			page, err := listMessagesPager.NextPage(context.Background())

			if err != nil {
				return err
			}

			lastAssistantResponses = append(lastAssistantResponses, page.Data...)
		}
	}
}

func printAssistantMessages(threadMessages []azopenaiassistants.ThreadMessage) {
	// print out the response contents for debugging.
	for _, response := range threadMessages {
		for _, content := range response.Content {

			switch v := content.(type) {
			case *azopenaiassistants.MessageImageFileContent:
				fmt.Fprintf(os.Stderr, "[ASSISTANT] %s: Image response, file ID: %s\n", *response.ID, *v.ImageFile.FileID.FileID)
			case *azopenaiassistants.MessageTextContent:
				fmt.Fprintf(os.Stderr, "[ASSISTANT] %s: Text response: %s\n", *response.ID, *v.Text.Value)
			}
		}
	}
}

func pollRunEnd(ctx context.Context, client *azopenaiassistants.Client, threadID string, runID string) error {
	for {
		lastGetRunResp, err := client.GetRun(context.Background(), threadID, runID, nil)

		if err != nil {
			return err
		}

		if *lastGetRunResp.Status != azopenaiassistants.RunStatusQueued && *lastGetRunResp.Status != azopenaiassistants.RunStatusInProgress {
			return nil
		}

		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
