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

func ExampleClient_GetChatCompletions_reasoningModels() {
	azureOpenAIKey := os.Getenv("AOAI_CHAT_COMPLETIONS_API_KEY")
	modelDeploymentID := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL") // Ex: "o1"

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_CHAT_COMPLETIONS_ENDPOINT")

	if azureOpenAIKey == "" || modelDeploymentID == "" || azureOpenAIEndpoint == "" {
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

	opts := azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestUserMessage{
				Content: azopenai.NewChatRequestUserMessageContent("I need to plan a trip to Europe for 10 days, visiting Paris, Rome, and London. Create a possible itinerary, including travel times and estimated costs."),
			},
		},
		DeploymentName:  &modelDeploymentID,
		ReasoningEffort: to.Ptr(azopenai.ReasoningEffortValueLow),
	}

	resp, err := client.GetChatCompletions(context.Background(), opts, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	fmt.Fprintln(os.Stderr, "\nLow reasoning")
	fmt.Fprintf(os.Stderr, "  Completion tokens: %v\n", *resp.Usage.CompletionTokens)
	fmt.Fprintf(os.Stderr, "  Reasoning tokens: %v\n", *resp.Usage.CompletionTokensDetails.ReasoningTokens)
	fmt.Fprintf(os.Stderr, "  Prompt tokens: %v\n", *resp.Usage.PromptTokens)
	fmt.Fprintf(os.Stderr, "  Total tokens: %v\n", *resp.Usage.TotalTokens)

	opts2 := azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestUserMessage{
				Content: azopenai.NewChatRequestUserMessageContent("I need to plan a trip to Europe for 10 days, visiting Paris, Rome, and London. Create a possible itinerary, including travel times and estimated costs."),
			},
		},
		DeploymentName:  &modelDeploymentID,
		ReasoningEffort: to.Ptr(azopenai.ReasoningEffortValueHigh),
	}

	resp2, err := client.GetChatCompletions(context.Background(), opts2, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	fmt.Fprintln(os.Stderr, "\nHigh reasoning")
	fmt.Fprintf(os.Stderr, "  Completion tokens: %v\n", *resp2.Usage.CompletionTokens)
	fmt.Fprintf(os.Stderr, "  Reasoning tokens: %v\n", *resp2.Usage.CompletionTokensDetails.ReasoningTokens)
	fmt.Fprintf(os.Stderr, "  Prompt tokens: %v\n", *resp2.Usage.PromptTokens)
	fmt.Fprintf(os.Stderr, "  Total tokens: %v\n", *resp2.Usage.TotalTokens)

	// Output:
}
