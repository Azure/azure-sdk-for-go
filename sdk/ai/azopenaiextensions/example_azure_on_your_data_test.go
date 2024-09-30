// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiextensions"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

func Example_usingAzureOnYourData() {
	endpoint := os.Getenv("AOAI_OYD_ENDPOINT")
	model := os.Getenv("AOAI_OYD_MODEL")
	cognitiveSearchEndpoint := os.Getenv("COGNITIVE_SEARCH_API_ENDPOINT") // Ex: https://<your-service>.search.windows.net
	cognitiveSearchIndexName := os.Getenv("COGNITIVE_SEARCH_API_INDEX")

	if endpoint == "" || model == "" || cognitiveSearchEndpoint == "" || cognitiveSearchIndexName == "" {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not \nrunning example.")
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-07-01-preview"),
		azure.WithTokenCredential(tokenCredential),
	)

	chatParams := openai.ChatCompletionNewParams{
		Model:     openai.F(model),
		MaxTokens: openai.Int(512),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionMessageParam{
				Role:    openai.F(openai.ChatCompletionMessageParamRoleUser),
				Content: openai.F[any]("What does the OpenAI package do?"),
			},
		}),
	}

	// There are other types of data sources available. Examples:
	//
	// - AzureCosmosDBChatExtensionConfiguration
	// - AzureMachineLearningIndexChatExtensionConfiguration
	// - AzureSearchChatExtensionConfiguration
	// - PineconeChatExtensionConfiguration
	//
	// See the definition of [AzureChatExtensionConfigurationClassification] for a full list.
	azureSearchDataSource := &azopenaiextensions.AzureSearchChatExtensionConfiguration{
		Parameters: &azopenaiextensions.AzureSearchChatExtensionParameters{
			Endpoint:       &cognitiveSearchEndpoint,
			IndexName:      &cognitiveSearchIndexName,
			Authentication: &azopenaiextensions.OnYourDataSystemAssignedManagedIdentityAuthenticationOptions{},
		},
	}

	resp, err := client.Chat.Completions.New(
		context.TODO(),
		chatParams,
		azopenaiextensions.WithDataSources(azureSearchDataSource),
	)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	for _, chatChoice := range resp.Choices {
		// Azure-specific response data can be extracted using helpers, like [azopenaiextensions.ChatCompletionChoice].
		azureChatChoice := azopenaiextensions.ChatCompletionChoice(chatChoice)
		azureContentFilterResult, err := azureChatChoice.ContentFilterResults()

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		if azureContentFilterResult != nil {
			fmt.Fprintf(os.Stderr, "ContentFilterResult: %#v\n", azureContentFilterResult)
		}

		// there are also helpers for individual types, not just top-level response types.
		azureChatCompletionMsg := azopenaiextensions.ChatCompletionMessage(chatChoice.Message)
		msgContext, err := azureChatCompletionMsg.Context()

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		for _, citation := range msgContext.Citations {
			if citation.Content != nil {
				fmt.Fprintf(os.Stderr, "Citation = %s\n", *citation.Content)
			}
		}

		// the original fields from the type are also still available.
		fmt.Fprintf(os.Stderr, "Content: %s\n", azureChatCompletionMsg.Content)
	}

	fmt.Printf("Example complete\n")

	// Output:
	// Example complete
	//
}

func Example_usingEnhancements() {
	endpoint := os.Getenv("AOAI_OYD_ENDPOINT")
	model := os.Getenv("AOAI_OYD_MODEL")
	cognitiveSearchEndpoint := os.Getenv("COGNITIVE_SEARCH_API_ENDPOINT") // Ex: https://<your-service>.search.windows.net
	cognitiveSearchIndexName := os.Getenv("COGNITIVE_SEARCH_API_INDEX")

	if endpoint == "" || model == "" || cognitiveSearchEndpoint == "" || cognitiveSearchIndexName == "" {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not \nrunning example.")
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-07-01-preview"),
		azure.WithTokenCredential(tokenCredential),
	)

	chatParams := openai.ChatCompletionNewParams{
		Model:     openai.F(model),
		MaxTokens: openai.Int(512),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionMessageParam{
				Role:    openai.F(openai.ChatCompletionMessageParamRoleUser),
				Content: openai.F[any]("What does the OpenAI package do?"),
			},
		}),
	}

	resp, err := client.Chat.Completions.New(
		context.TODO(),
		chatParams,
		azopenaiextensions.WithEnhancements(azopenaiextensions.AzureChatEnhancementConfiguration{
			Grounding: &azopenaiextensions.AzureChatGroundingEnhancementConfiguration{
				Enabled: to.Ptr(true),
			},
		}),
	)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	for _, chatChoice := range resp.Choices {
		// Azure-specific response data can be extracted using helpers, like [azopenaiextensions.ChatCompletionChoice].
		azureChatChoice := azopenaiextensions.ChatCompletionChoice(chatChoice)
		azureContentFilterResult, err := azureChatChoice.ContentFilterResults()

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		if azureContentFilterResult != nil {
			fmt.Fprintf(os.Stderr, "ContentFilterResult: %#v\n", azureContentFilterResult)
		}

		// there are also helpers for individual types, not just top-level response types.
		azureChatCompletionMsg := azopenaiextensions.ChatCompletionMessage(chatChoice.Message)
		msgContext, err := azureChatCompletionMsg.Context()

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return
		}

		for _, citation := range msgContext.Citations {
			if citation.Content != nil {
				fmt.Fprintf(os.Stderr, "Citation = %s\n", *citation.Content)
			}
		}

		// the original fields from the type are also still available.
		fmt.Fprintf(os.Stderr, "Content: %s\n", azureChatCompletionMsg.Content)
	}

	fmt.Printf("Example complete\n")

	// Output:
	// Example complete
	//
}
