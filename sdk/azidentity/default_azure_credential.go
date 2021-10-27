// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	organizationsTenantID   = "organizations"
	developerSignOnClientID = "04b07795-8ddb-461a-bbee-02f9e1bf7b46"
)

// DefaultAzureCredentialOptions contains options for configuring authentication. These options
// may not apply to all credentials in the default chain.
type DefaultAzureCredentialOptions struct {
	azcore.ClientOptions

	// The host of the Azure Active Directory authority. The default is AzurePublicCloud.
	// Leave empty to allow overriding the value from the AZURE_AUTHORITY_HOST environment variable.
	AuthorityHost AuthorityHost
	// TenantID identifies the tenant the Azure CLI should authenticate in.
	// Defaults to the CLI's default tenant, which is typically the home tenant of the user logged in to the CLI.
	TenantID string
}

// DefaultAzureCredential is a default credential chain for applications that will be deployed to Azure.
// It combines credentials suitable for deployed applications with credentials suitable in local development.
// It attempts to authenticate with each of these credential types, in the following order:
// - EnvironmentCredential
// - ManagedIdentityCredential
// - AzureCLICredential
// Consult the documentation for these credential types for more information on how they authenticate.
type DefaultAzureCredential struct {
	chain *ChainedTokenCredential
}

// NewDefaultAzureCredential creates a default credential chain for applications that will be deployed to Azure.
func NewDefaultAzureCredential(options *DefaultAzureCredentialOptions) (*DefaultAzureCredential, error) {
	var creds []azcore.TokenCredential
	errMsg := ""

	cp := DefaultAzureCredentialOptions{}
	if options != nil {
		cp = *options
	}

	envCred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{AuthorityHost: cp.AuthorityHost,
		ClientOptions: cp.ClientOptions,
	})
	if err == nil {
		creds = append(creds, envCred)
	} else {
		errMsg += err.Error()
	}

	msiCred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: cp.ClientOptions})
	if err == nil {
		creds = append(creds, msiCred)
	} else {
		errMsg += err.Error()
	}

	cliCred, err := NewAzureCLICredential(&AzureCLICredentialOptions{TenantID: cp.TenantID})
	if err == nil {
		creds = append(creds, cliCred)
	} else {
		errMsg += err.Error()
	}

	if len(creds) == 0 {
		err := errors.New(errMsg)
		logCredentialError("Default Azure Credential", err)
		return nil, err
	}
	chain, err := NewChainedTokenCredential(creds, nil)
	if err != nil {
		return nil, err
	}
	return &DefaultAzureCredential{chain: chain}, nil
}

// GetToken attempts to acquire a token from each of the default chain's credentials, stopping when one provides a token.
func (c *DefaultAzureCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (token *azcore.AccessToken, err error) {
	return c.chain.GetToken(ctx, opts)
}
