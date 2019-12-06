// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	developerSignOnClientID = "04b07795-8ddb-461a-bbee-02f9e1bf7b46"
)

// DefaultTokenCredential provides a default ChainedTokenCredential configuration for applications that will be deployed to Azure.  The following credential
// types will be tried, in order:
// - EnvironmentCredential
// - ManagedIdentityCredential
// Consult the documentation of these credential types for more information on how they attempt authentication.
type DefaultTokenCredential struct {
	sources []azcore.TokenCredential
}

// DefaultTokenCredentialOptions contain information that can configure Default Token Credentials
type DefaultTokenCredentialOptions struct {
	ExcludeEnvironmentCredential bool
	ExcludeMSICredential         bool
}

// NewDefaultTokenCredential provides a default ChainedTokenCredential configuration for applications that will be deployed to Azure.  The following credential
// types will be tried, in order:
// - EnvironmentCredential
// - ManagedIdentityCredential
// Consult the documentation of these credential types for more information on how they attempt authentication.
func NewDefaultTokenCredential(options *DefaultTokenCredentialOptions) (*ChainedTokenCredential, error) {
	cred := &DefaultTokenCredential{}

	if options == nil {
		options = &DefaultTokenCredentialOptions{}
	}

	if !options.ExcludeEnvironmentCredential {
		envCred, err := NewEnvironmentCredential(nil)
		if err == nil {
			cred.sources = append(cred.sources, envCred)
			return NewChainedTokenCredential(cred.sources...)
		}
	}

	if !options.ExcludeMSICredential {
		cred.sources = append(cred.sources, NewManagedIdentityCredential("", nil))
		return NewChainedTokenCredential(cred.sources...)
	}

	return nil, &AuthenticationFailedError{Message: "Failed to find any available credentials. Make sure you are running in a Managed Identity Environment or have set the appropriate environment variables"}
}
