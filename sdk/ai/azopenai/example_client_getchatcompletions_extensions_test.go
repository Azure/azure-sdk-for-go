// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

// With Azure OpenAI you can integrate data you've already uploaded to an Azure Cognitive Search index.
// For more information about this feature see the article "[Azure OpenAI on your data]".
//
// [Azure OpenAI on your data]: https://learn.microsoft.com/azure/ai-services/openai/concepts/use-your-data
func ExampleClient_GetChatCompletions_bringYourOwnDataWithCognitiveSearch() {
	azureOpenAIKey := os.Getenv("AOAI_CHAT_COMPLETIONS_RAI_API_KEY")
	modelDeploymentID := os.Getenv("AOAI_CHAT_COMPLETIONS_RAI_MODEL")

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_RAI_ENDPOINT")

	// Azure Cognitive Search configuration
	searchIndex := os.Getenv("COGNITIVE_SEARCH_API_INDEX")
	searchEndpoint := os.Getenv("COGNITIVE_SEARCH_API_ENDPOINT")
	searchAPIKey := os.Getenv("COGNITIVE_SEARCH_API_KEY")

	if azureOpenAIKey == "" || modelDeploymentID == "" || azureOpenAIEndpoint == "" || searchIndex == "" || searchEndpoint == "" || searchAPIKey == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

	// In Azure OpenAI you must deploy a model before you can use it in your client. For more information
	// see here: https://learn.microsoft.com/azure/cognitive-services/openai/how-to/create-resource
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("What are the differences between Azure Machine Learning and Azure AI services?")},
		},
		MaxTokens: to.Ptr[int32](512),
		AzureExtensionsOptions: []azopenai.AzureChatExtensionConfigurationClassification{
			&azopenai.AzureSearchChatExtensionConfiguration{
				// This allows Azure OpenAI to use an Azure Cognitive Search index.
				//
				// > Because the model has access to, and can reference specific sources to support its responses, answers are not only based on its pretrained knowledge
				// > but also on the latest information available in the designated data source. This grounding data also helps the model avoid generating responses
				// > based on outdated or incorrect information.
				//
				// Quote from here: https://learn.microsoft.com/en-us/azure/ai-services/openai/concepts/use-your-data
				Parameters: &azopenai.AzureSearchChatExtensionParameters{
					Endpoint:  &searchEndpoint,
					IndexName: &searchIndex,
					Authentication: &azopenai.OnYourDataAPIKeyAuthenticationOptions{
						Key: &searchAPIKey,
					},
				},
			},
		},
		DeploymentName: &modelDeploymentID,
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	// Contains contextual information from your Azure chat completion extensions, configured above in `AzureExtensionsOptions`
	msgContext := resp.Choices[0].Message.Context

	fmt.Fprintf(os.Stderr, "Extensions Context (length): %d\n", len(*msgContext.Citations[0].Content))

	fmt.Fprintf(os.Stderr, "ChatRole: %s\nChat content: %s\n",
		*resp.Choices[0].Message.Role,
		*resp.Choices[0].Message.Content,
	)

	// Output:
}
