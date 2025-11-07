// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/openai/openai-go/v3"
)

// Example_usingAzureOnYourData demonstrates how to use Azure OpenAI's Azure-On-Your-Data feature.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Configure an Azure Cognitive Search data source
// - Send a chat completion request with data source integration
// - Process Azure-specific response data including citations and content filtering results
//
// The example uses environment variables for configuration:
// - AOAI_OYD_ENDPOINT: Your Azure OpenAI endpoint URL
// - AOAI_OYD_MODEL: The deployment name of your model
// - COGNITIVE_SEARCH_API_ENDPOINT: Your Azure Cognitive Search endpoint
// - COGNITIVE_SEARCH_API_INDEX: The name of your search index
//
// Azure-On-Your-Data enables you to enhance chat completions with information from your
// own data sources, allowing for more contextual and accurate responses based on your content.
func Example_usingAzureOnYourData() {
	if !CheckRequiredEnvVars("AOAI_OYD_ENDPOINT", "AOAI_OYD_MODEL",
		"COGNITIVE_SEARCH_API_ENDPOINT", "COGNITIVE_SEARCH_API_INDEX") {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not \nrunning example.")
		return
	}

	endpoint := os.Getenv("AOAI_OYD_ENDPOINT")
	model := os.Getenv("AOAI_OYD_MODEL")
	cognitiveSearchEndpoint := os.Getenv("COGNITIVE_SEARCH_API_ENDPOINT")
	cognitiveSearchIndexName := os.Getenv("COGNITIVE_SEARCH_API_INDEX")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	chatParams := openai.ChatCompletionNewParams{
		Model:     openai.ChatModel(model),
		MaxTokens: openai.Int(512),
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("What does the OpenAI package do?"),
				},
			},
		}},
	}

	// There are other types of data sources available. Examples:
	//
	// - AzureCosmosDBChatExtensionConfiguration
	// - AzureMachineLearningIndexChatExtensionConfiguration
	// - AzureSearchChatExtensionConfiguration
	// - PineconeChatExtensionConfiguration
	//
	// See the definition of [AzureChatExtensionConfigurationClassification] for a full list.
	azureSearchDataSource := &azopenai.AzureSearchChatExtensionConfiguration{
		Parameters: &azopenai.AzureSearchChatExtensionParameters{
			Endpoint:       &cognitiveSearchEndpoint,
			IndexName:      &cognitiveSearchIndexName,
			Authentication: &azopenai.OnYourDataSystemAssignedManagedIdentityAuthenticationOptions{},
		},
	}

	resp, err := client.Chat.Completions.New(
		context.TODO(),
		chatParams,
		azopenai.WithDataSources(azureSearchDataSource),
	)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	for _, chatChoice := range resp.Choices {
		// Azure-specific response data can be extracted using helpers, like [azopenai.ChatCompletionChoice].
		azureChatChoice := azopenai.ChatCompletionChoice(chatChoice)
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
		azureChatCompletionMsg := azopenai.ChatCompletionMessage(chatChoice.Message)
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

	fmt.Fprintf(os.Stderr, "Example complete\n")
}

// Example_usingEnhancements demonstrates how to use Azure OpenAI's enhanced features.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Configure chat completion enhancements like grounding
// - Process Azure-specific response data including content filtering
// - Handle message context and citations
//
// The example uses environment variables for configuration:
// - AOAI_OYD_ENDPOINT: Your Azure OpenAI endpoint URL
// - AOAI_OYD_MODEL: The deployment name of your model
//
// Azure OpenAI enhancements provide additional capabilities beyond standard OpenAI features,
// such as improved grounding and content filtering for more accurate and controlled responses.
func Example_usingEnhancements() {
	if !CheckRequiredEnvVars("AOAI_OYD_ENDPOINT", "AOAI_OYD_MODEL") {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not \nrunning example.")
		return
	}

	endpoint := os.Getenv("AOAI_OYD_ENDPOINT")
	model := os.Getenv("AOAI_OYD_MODEL")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	chatParams := openai.ChatCompletionNewParams{
		Model:     openai.ChatModel(model),
		MaxTokens: openai.Int(512),
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("What does the OpenAI package do?"),
				},
			},
		}},
	}

	resp, err := client.Chat.Completions.New(
		context.TODO(),
		chatParams,
		azopenai.WithEnhancements(azopenai.AzureChatEnhancementConfiguration{
			Grounding: &azopenai.AzureChatGroundingEnhancementConfiguration{
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
		// Azure-specific response data can be extracted using helpers, like [azopenai.ChatCompletionChoice].
		azureChatChoice := azopenai.ChatCompletionChoice(chatChoice)
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
		azureChatCompletionMsg := azopenai.ChatCompletionMessage(chatChoice.Message)
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

	fmt.Fprintf(os.Stderr, "Example complete\n")
}
