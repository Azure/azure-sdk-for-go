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

func Example_assistantsWithConversationLoop() {
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
	createAssistantResp, err := client.CreateAssistant(context.Background(), azopenaiassistants.CreateAssistantBody{
		Name:           &assistantName,
		DeploymentName: to.Ptr("gpt-4-1106-preview"),
		Instructions:   to.Ptr("You are a personal math tutor. Write and run code to answer math questions."),
		Tools: []azopenaiassistants.ToolDefinitionClassification{
			&azopenaiassistants.CodeInterpreterToolDefinition{},
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
	createThreadResp, err := client.CreateThread(context.Background(), azopenaiassistants.CreateThreadBody{}, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	threadID := createThreadResp.ID

	assistantCtx, stopAssistant := context.WithCancel(context.TODO())

	callIdx := -1

	// This is just a simplified example of how you could handle a conversation - `assistantMessages` are the messages that
	// are responses from the assistant, and you return messages from here that are then added to the conversation.
	handleConversation := func(ctx context.Context, assistantMessages []azopenaiassistants.ThreadMessage) ([]azopenaiassistants.CreateMessageBody, error) {
		callIdx++

		if err := printAssistantMessages(ctx, client, assistantMessages); err != nil {
			return nil, err
		}

		// For this example we'll just synthesize some responses, simulating a conversation.
		// In a real application these messages would come from the user, responding to replies
		// from the assistant.
		switch callIdx {
		case 0:
			text := "Can you help me find the y intercept for y = x + 4?"
			fmt.Fprintf(os.Stderr, "[ME] %s\n", text)

			return []azopenaiassistants.CreateMessageBody{
				{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: &text},
			}, nil
		case 1:
			text := "Can you explain it with a Python program?"
			fmt.Fprintf(os.Stderr, "[ME] %s\n", text)

			return []azopenaiassistants.CreateMessageBody{
				{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: &text},
			}, nil
		case 2:
			text := "Can you give me the result if that Python program had 'x' set to 10"
			fmt.Fprintf(os.Stderr, "[ME] %s\n", text)

			return []azopenaiassistants.CreateMessageBody{
				{Role: to.Ptr(azopenaiassistants.MessageRoleUser), Content: &text},
			}, nil
		default:
			stopAssistant()
		}
		return nil, nil
	}

	if err = assistantLoop(assistantCtx, client, *assistantID, *threadID, handleConversation); err != nil {
		// if this is a cancellation error it's just us trying to stop the assistant loop.
		if errors.Is(err, context.Canceled) {
			fmt.Fprintf(os.Stderr, "Assistant stopped cleanly\n")
		} else {
			//  TODO: Update the following line with your application specific error handling logic
			log.Printf("ERROR: %s", err)
			return
		}
	}

	// DisabledOutput:
}

// conversationHandler takes responses from an assistant and returns our reply messages. Returns the responses
// based on the contents of assistantMessages
// - assistantMessages - messages that have arrived since our last read of the thread.
type conversationHandler func(ctx context.Context, assistantMessages []azopenaiassistants.ThreadMessage) ([]azopenaiassistants.CreateMessageBody, error)

func assistantLoop(ctx context.Context, client *azopenaiassistants.Client,
	assistantID string, threadID string,
	handleConversation conversationHandler) error {
	// from here we'll run in a loop, adding new messages to the conversation and reading the assistants
	// responses.

	var lastAssistantResponses []azopenaiassistants.ThreadMessage

	for {
		yourResponses, err := handleConversation(ctx, lastAssistantResponses)

		if err != nil {
			return err
		}

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

		if _, err := pollUntilRunEnds(ctx, client, threadID, runID); err != nil {
			return err
		}

		// get all the messages that were added after our most recently added message.
		lastAssistantResponses, err = getLatestMessages(ctx, client, threadID, lastMessageID)

		if err != nil {
			return err
		}
	}
}

func printAssistantMessages(ctx context.Context, client *azopenaiassistants.Client, threadMessages []azopenaiassistants.ThreadMessage) error {
	// print out the response contents for debugging.
	for _, response := range threadMessages {
		for _, content := range response.Content {
			switch v := content.(type) {
			case *azopenaiassistants.MessageImageFileContent:
				fmt.Fprintf(os.Stderr, "[ASSISTANT] Image response, file ID: %s\n", *v.ImageFile.FileID)

				// Download the contents of the file through the returned reader.
				fileContentResp, err := client.GetFileContent(ctx, *v.ImageFile.FileID, nil)

				if err != nil {
					return err
				}

				contents, err := io.ReadAll(fileContentResp.Content)

				if err != nil {
					return err
				}

				fmt.Fprintf(os.Stderr, "  File contents downloaded, length %d\n", len(contents))
			case *azopenaiassistants.MessageTextContent:
				fmt.Fprintf(os.Stderr, "[ASSISTANT] %s: Text response: %s\n", *response.ID, *v.Text.Value)
			}
		}
	}

	return nil
}

func pollUntilRunEnds(ctx context.Context, client *azopenaiassistants.Client, threadID string, runID string) (azopenaiassistants.GetRunResponse, error) {
	for {
		lastGetRunResp, err := client.GetRun(context.Background(), threadID, runID, nil)

		if err != nil {
			return azopenaiassistants.GetRunResponse{}, err
		}

		switch *lastGetRunResp.Status {
		case azopenaiassistants.RunStatusInProgress, azopenaiassistants.RunStatusQueued:
			// we're either running or about to run so we'll just keep polling for the end.
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return azopenaiassistants.GetRunResponse{}, ctx.Err()
			}
		case azopenaiassistants.RunStatusRequiresAction:
			// The assistant run has stopped because a tool requires you to submit inputs.
			// You can see an example of this in Example_assistantsUsingFunctionTool.
			return lastGetRunResp, nil
		case azopenaiassistants.RunStatusCompleted:
			// The run has completed successfully
			return lastGetRunResp, nil
		case azopenaiassistants.RunStatusFailed:
			// The run has failed. We can use the code and message to give us an idea of why.
			var code, description string

			if lastGetRunResp.LastError != nil && lastGetRunResp.LastError.Code != nil {
				code = *lastGetRunResp.LastError.Code
			}

			if lastGetRunResp.LastError != nil && lastGetRunResp.LastError.Message != nil {
				description = *lastGetRunResp.LastError.Message
			}

			return lastGetRunResp, fmt.Errorf("run failed, code: %s, message: %s", code, description)

		default:
			return azopenaiassistants.GetRunResponse{}, fmt.Errorf("run ended but status was not complete: %s", *lastGetRunResp.Status)
		}
	}
}

// getLatestMessages gets any messages that have occurred since lastMessageID.
// If an error occurs, returns any messages received so far, as well as the error.
func getLatestMessages(ctx context.Context, client *azopenaiassistants.Client, threadID string, lastMessageID *string) ([]azopenaiassistants.ThreadMessage, error) {
	// grab any messages that occurred after our last known message
	listMessagesPager := client.NewListMessagesPager(threadID, &azopenaiassistants.ListMessagesOptions{
		After: lastMessageID,
		Order: to.Ptr(azopenaiassistants.ListSortOrderAscending),
	})

	var all []azopenaiassistants.ThreadMessage

	for listMessagesPager.More() {
		page, err := listMessagesPager.NextPage(ctx)

		if err != nil {
			return all, err
		}

		all = append(all, page.Data...)
	}

	return all, nil
}
