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

// Example_usingDefaultAzureCredential demonstrates how to authenticate with Azure OpenAI using Azure Active Directory credentials.
// This example shows how to:
// - Create an Azure OpenAI client using DefaultAzureCredential
// - Configure authentication options with tenant ID
// - Make a simple request to test the authentication
//
// The example uses environment variables for configuration:
// - AOAI_ENDPOINT: Your Azure OpenAI endpoint URL
// - AOAI_MODEL: The deployment name of your model
// - AZURE_TENANT_ID: Your Azure tenant ID
// - AZURE_CLIENT_ID: (Optional) Your Azure client ID
// - AZURE_CLIENT_SECRET: (Optional) Your Azure client secret
//
// DefaultAzureCredential supports multiple authentication methods including:
// - Environment variables
// - Managed Identity
// - Azure CLI credentials
func Example_usingDefaultAzureCredential() {
	if !CheckRequiredEnvVars("AOAI_ENDPOINT", "AOAI_MODEL") {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not running example.")
		return
	}

	endpoint := os.Getenv("AOAI_ENDPOINT")
	model := os.Getenv("AOAI_MODEL")
	tenantID := os.Getenv("AZURE_TENANT_ID")

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

// Example_usingManagedIdentityCredential demonstrates how to authenticate with Azure OpenAI using Managed Identity.
// This example shows how to:
// - Create an Azure OpenAI client using ManagedIdentityCredential
// - Support both system-assigned and user-assigned managed identities
// - Make authenticated requests without storing credentials
//
// The example uses environment variables for configuration:
// - AOAI_ENDPOINT: Your Azure OpenAI endpoint URL
// - AOAI_MODEL: The deployment name of your model
//
// Managed Identity is ideal for:
// - Azure services (VMs, App Service, Azure Functions, etc.)
// - Azure DevOps pipelines with the Azure DevOps service connection
// - CI/CD scenarios where you want to avoid storing secrets
// - Production workloads requiring secure, credential-free authentication
func Example_usingManagedIdentityCredential() {
	endpoint := os.Getenv("AOAI_ENDPOINT")
	model := os.Getenv("AOAI_MODEL")

	if endpoint == "" || model == "" {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not running example.")
		return
	}

	var credential *azidentity.ManagedIdentityCredential
	var err error

	// Use system assigned managed identity
	credential, err = azidentity.NewManagedIdentityCredential(nil)

	// When using User Assigned Managed Identity use this instead and pass your client id in the options
	// clientID := azidentity.ClientID("abcd1234-...")
	// opts := azidentity.ManagedIdentityCredentialOptions{ID: clientID}
	// cred, err := azidentity.NewManagedIdentityCredential(&opts)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-08-01-preview"),
		azure.WithTokenCredential(credential),
	)

	// Use the client with managed identity credentials
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
