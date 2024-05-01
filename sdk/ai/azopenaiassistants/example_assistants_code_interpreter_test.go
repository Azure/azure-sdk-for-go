//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func Example_assistantsUsingCodeInterpreter() {
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
		Instructions:   to.Ptr("You are an AI assistant that can write code to help answer math questions."),
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
	createThreadResp, err := client.CreateThread(context.TODO(), azopenaiassistants.AssistantThreadCreationOptions{}, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	threadID := *createThreadResp.ID

	// Add a user question to the thread
	ourQuestion, err := client.CreateMessage(context.TODO(), threadID, azopenaiassistants.CreateMessageBody{
		Content: to.Ptr("I need to solve the equation `3x + 11 = 14`. Can you help me?"),
		Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Fprintf(os.Stderr, "[USER] I need to solve the equation `3x + 11 = 14`. Can you help me?\n")

	// Run the thread and wait (using pollRunEnd) until it completes.
	threadRun, err := client.CreateRun(context.TODO(), threadID, azopenaiassistants.CreateRunBody{
		AssistantID: assistantID,
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Wait till the assistant has responded
	if err := pollRunEnd(context.TODO(), client, threadID, *threadRun.ID); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// retrieve any messages added after we asked our question.
	listMessagesPager := client.NewListMessagesPager(threadID, &azopenaiassistants.ListMessagesOptions{
		After: ourQuestion.ID,
		Order: to.Ptr(azopenaiassistants.ListSortOrderAscending),
	})

	for listMessagesPager.More() {
		page, err := listMessagesPager.NextPage(context.Background())

		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		for _, threadMessage := range page.Data {
			for _, content := range threadMessage.Content {
				switch v := content.(type) {
				case *azopenaiassistants.MessageTextContent:
					fmt.Fprintf(os.Stderr, "[ASSISTANT] %s: Text response: %s\n", *threadMessage.ID, *v.Text.Value)
				}
			}
		}
	}

	// Output:
}
