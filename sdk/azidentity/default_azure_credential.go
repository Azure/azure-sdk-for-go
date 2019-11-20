// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

// -------------------- NOTES ------------------------------
/*
Currently all of the languages implement the DefaultAzureCredential as an abstraction over the ChainedTokenCredential
DAC only calls env cred and msi cred (sdks with msal include shared token cache)
There is no guarantee that the credential type used will not change

*/
import (
	"errors"

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
	// Scopes []string
	// EnvCredential *ClientSecretCredential
	// MSICredential *ManagedIdentityCredential // A pointer
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

	if options == nil {
		// CP: these could be set up once in an init?
		envCred, err := NewEnvironmentCredential(nil)
		if err != nil {
			if errors.Is(err, &AuthenticationFailedError{}) {
				return err
			}
			errCount++
		}
		msiCred, err := NewManagedIdentityCredential("", nil)
		if err != nil {
			// CP: this check shouldnt be necessary since we arent authenticating here, just creating the credential? Unless we wrap all other unexpected errors in an auth failed error
			if errors.Is(err, &AuthenticationFailedError{}) {
				return err
			}
			errCount++
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

	if !options.ExcludeEnvironmentCredential {
		envCred, err := NewEnvironmentCredential(nil)
		if err != nil {
			if errors.Is(err, &AuthenticationFailedError{}) {
				return err
			}
			errCount++
		}
		c.sources = append(c.sources, envCred)
	}

	if !options.ExcludeMSICredential {
		msiCred, err := NewManagedIdentityCredential("", nil)
		if err != nil {
			if errors.Is(err, &AuthenticationFailedError{}) {
				return err
			}
			errCount++
		}
		c.sources = append(c.sources, msiCred)
	}

	return nil
}
