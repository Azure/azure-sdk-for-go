// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

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
	makeSimpleRequest(client, model)
}

func Example_usingClientSecretCredential() {
	endpoint := os.Getenv("AOAI_ENDPOINT")
	model := os.Getenv("AOAI_MODEL")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	if endpoint == "" || model == "" || tenantID == "" || clientID == "" || clientSecret == "" {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not running example.")
		return
	}

	// ClientSecretCredential is used when you have a service principal
	// with client ID and client secret
	credential, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-08-01-preview"),
		azure.WithTokenCredential(credential),
	)

	// Use the client with service principal credentials
	makeSimpleRequest(client, model)
}

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
	makeSimpleRequest(client, model)
}

func Example_usingInteractiveBrowserCredential() {
	endpoint := os.Getenv("AOAI_ENDPOINT")
	model := os.Getenv("AOAI_MODEL")

	if endpoint == "" || model == "" {
		fmt.Fprintf(os.Stderr, "Environment variables are not set, not running example.")
		return
	}

	// InteractiveBrowserCredential authenticates a user by opening the default browser
	// to the Azure login page and waiting for the user to complete the login process
	//
	// Optional configurations can be specified using InteractiveBrowserCredentialOptions:
	// options := &azidentity.InteractiveBrowserCredentialOptions{
	//     TenantID: "<tenant-id>",                    // Specify a tenant for authentication
	//     ClientID: "<client-id>",                    // Use a custom client ID
	//     RedirectURL: "http://localhost",            // Custom redirect URL
	//     LoginHint: "user@contoso.com",             // Pre-fill username field
	// }
	credential, err := azidentity.NewInteractiveBrowserCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-08-01-preview"),
		azure.WithTokenCredential(credential),
	)

	// Use the client with interactive browser credentials
	makeSimpleRequest(client, model)
}

// Helper function to make a simple request to Azure OpenAI
func makeSimpleRequest(client *openai.Client, model string) {
	chatParams := openai.ChatCompletionNewParams{
		Model:     openai.F(model),
		MaxTokens: openai.Int(100),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionMessageParam{
				Role:    openai.F(openai.ChatCompletionMessageParamRoleUser),
				Content: openai.F[any]("Say hello!"),
			},
		}),
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
