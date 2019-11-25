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
	err := cred.getDefaultTokenCredentialChain(options)
	if err != nil {
		return nil, err
	}
	return NewChainedTokenCredential(cred.sources...)
}

func (c *DefaultTokenCredential) getDefaultTokenCredentialChain(options *DefaultTokenCredentialOptions) error {
	errCount := 0
	var envCred *ClientSecretCredential
	var msiCred *ManagedIdentityCredential
	var err error

	if options == nil {
		envCred, err = NewEnvironmentCredential(nil)
		if err != nil {
			errCount++
		}
		msiCred, err = NewManagedIdentityCredential("", nil)
		if err != nil {
			errCount++
		}
	}

	if !options.ExcludeEnvironmentCredential {
		envCred, err = NewEnvironmentCredential(nil)
		if err != nil {
			errCount++
		}
	}

	if !options.ExcludeMSICredential {
		msiCred, err = NewManagedIdentityCredential("", nil)
		if err != nil {
			errCount++
		}
	}

	if errCount > 1 {
		return &AuthenticationFailedError{Message: "Failed to find any available credentials. Make sure you are running in a Managed Identity Environment or have set the appropriate environment variables"}
	}

	if envCred != nil {
		c.sources = append(c.sources, envCred)
	}

	if msiCred != nil {
		c.sources = append(c.sources, msiCred)
	}

	return nil
}
