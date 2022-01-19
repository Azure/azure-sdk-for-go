// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// DefaultAzureCredentialOptions contains optional parameters for DefaultAzureCredential.
// These options may not apply to all credentials in the chain.
type DefaultAzureCredentialOptions struct {
	azcore.ClientOptions

	// AuthorityHost is the base URL of an Azure Active Directory authority. Defaults
	// to the value of environment variable AZURE_AUTHORITY_HOST, if set, or AzurePublicCloud.
	AuthorityHost AuthorityHost
	// TenantID identifies the tenant the Azure CLI should authenticate in.
	// Defaults to the CLI's default tenant, which is typically the home tenant of the user logged in to the CLI.
	TenantID string
}

// DefaultAzureCredential is a default credential chain for applications that will deploy to Azure.
// It combines credentials suitable for deployment with credentials suitable for local development.
// It attempts to authenticate with each of these credential types, in the following order, stopping when one provides a token:
// - EnvironmentCredential
// - ManagedIdentityCredential
// - AzureCLICredential
// Consult the documentation for these credential types for more information on how they authenticate.
type DefaultAzureCredential struct {
	chain *ChainedTokenCredential
}

// NewDefaultAzureCredential creates a DefaultAzureCredential.
func NewDefaultAzureCredential(options *DefaultAzureCredentialOptions) (*DefaultAzureCredential, error) {
	var creds []azcore.TokenCredential
	var errorMessages []string

	if options == nil {
		options = &DefaultAzureCredentialOptions{}
	}

	envCred, err := NewEnvironmentCredential(
		&EnvironmentCredentialOptions{AuthorityHost: options.AuthorityHost, ClientOptions: options.ClientOptions},
	)
	if err == nil {
		creds = append(creds, envCred)
	} else {
		errorMessages = append(errorMessages, fmt.Sprintf("EnvironmentCredential: %s", err.Error()))
	}

	msiCred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{ClientOptions: options.ClientOptions})
	if err == nil {
		creds = append(creds, msiCred)
		msiCred.client.imdsTimeout = time.Second
	} else {
		errorMessages = append(errorMessages, fmt.Sprintf("ManagedIdentityCredential: %s", err.Error()))
	}

	cliCred, err := NewAzureCLICredential(&AzureCLICredentialOptions{TenantID: options.TenantID})
	if err == nil {
		creds = append(creds, cliCred)
	} else {
		errorMessages = append(errorMessages, fmt.Sprintf("AzureCLICredential: %s", err.Error()))
	}

	errMsg := "\n"
	for _, msg := range errorMessages {
		errMsg += fmt.Sprintf("\t%s\n", msg)
	}

	if len(creds) == 0 {
		err := errors.New(errMsg)
		logCredentialError("Default Azure Credential", err)
		return nil, err
	}

	if len(errMsg) != 0 {
		log.Writef(EventAuthentication, "Azure Identity => Failed to initialize some credentials on the Default Azure Credential:%s", errMsg)
	}

	chain, err := NewChainedTokenCredential(creds, nil)
	if err != nil {
		return nil, err
	}
	chain.chainedCredentialName = "DefaultAzureCredential"
	return &DefaultAzureCredential{chain: chain}, nil
}

// GetToken obtains a token from Azure Active Directory. This method is called automatically by Azure SDK clients.
// ctx: Context used to control the request lifetime.
// opts: Options for the token request, in particular the desired scope of the access token.
func (c *DefaultAzureCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (token *azcore.AccessToken, err error) {
	return c.chain.GetToken(ctx, opts)
}
