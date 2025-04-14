// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

// This example demonstrates how to use different Azure authentication methods
// with Azure OpenAI Services
func Example_usingDefaultAzureCredential() {
	endpoint := os.Getenv("AOAI_ENDPOINT")
	model := os.Getenv("AOAI_MODEL")
	tenantID := os.Getenv("AZURE_TENANT_ID")

	if endpoint == "" || model == "" {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not running example.")
		return
	}

	// DefaultAzureCredential automatically tries different authentication methods in order:
	// - Environment variables (AZURE_CLIENT_ID, AZURE_CLIENT_SECRET, AZURE_TENANT_ID)
	// - Managed Identity
	// - Azure CLI credentials
	credential, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		TenantID: tenantID,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-08-01-preview"),
		azure.WithTokenCredential(credential),
	)

	// Use the client with default credentials
	makeSimpleRequest(&client, model)
}

// Helper function to make a simple request to Azure OpenAI
func makeSimpleRequest(client *openai.Client, model string) {
	chatParams := openai.ChatCompletionNewParams{
		Model:     openai.ChatModel(model),
		MaxTokens: openai.Int(512),
		Messages: []openai.ChatCompletionMessageParamUnion{{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String("Say hello!"),
				},
			},
		}},
	}

	resp, err := client.Chat.Completions.New(
		context.TODO(),
		chatParams,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	if len(resp.Choices) > 0 {
		fmt.Fprintf(os.Stderr, "Response: %s\n", resp.Choices[0].Message.Content)
	}
}
