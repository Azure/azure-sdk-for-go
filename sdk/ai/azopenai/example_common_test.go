// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/option"
)

// Default API version to use for Azure OpenAI
const DefaultAPIVersion = "2024-08-01-preview"

// CreateOpenAIClientWithToken creates an OpenAI client with Azure AD token authentication
func CreateOpenAIClientWithToken(endpoint string, apiVersion string) (*openai.Client, error) {
	if apiVersion == "" {
		apiVersion = DefaultAPIVersion
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create default Azure credential: %w", err)
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, apiVersion),
		azure.WithTokenCredential(tokenCredential),
	)

	return &client, nil
}

// CreateOpenAIClientWithKey creates an OpenAI client with API key authentication
func CreateOpenAIClientWithKey(endpoint string, apiKey string, apiVersion string) (*openai.Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey cannot be empty")
	}

	if apiVersion == "" {
		apiVersion = DefaultAPIVersion
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		azure.WithEndpoint(endpoint, apiVersion),
	)

	return &client, nil
}

// GetRequiredEnvVar retrieves an environment variable and returns an error if it's not set
func GetRequiredEnvVar(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", name)
	}
	return value, nil
}

// CheckRequiredEnvVars checks if all required environment variables are set
// Returns true if all variables are set, false otherwise
func CheckRequiredEnvVars(names ...string) bool {
	for _, name := range names {
		if os.Getenv(name) == "" {
			fmt.Fprintf(os.Stderr, "Required environment variable '%s' is not set\n", name)
			return false
		}
	}
	return true
}
