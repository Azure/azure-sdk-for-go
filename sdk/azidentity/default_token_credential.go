// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

const (
	developerSignOnClientID = "04b07795-8ddb-461a-bbee-02f9e1bf7b46"
)

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
	var creds []azcore.TokenCredential
	errMsg := ""

	if options == nil {
		options = &DefaultTokenCredentialOptions{}
	}

	if !options.ExcludeEnvironmentCredential {
		envCred, err := NewEnvironmentCredential(nil)
		if err == nil {
			creds = append(creds, envCred)
		}
		errMsg += "Make sure you have set the environment variables necessary for authentication: AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET."
	}

	if !options.ExcludeMSICredential {
		creds = append(creds, NewManagedIdentityCredential("", nil))
	}

	if len(creds) == 0 {
		return nil, &CredentialUnavailableError{CredentialType: "Default Token Credential", Message: errMsg}
	}
	return NewChainedTokenCredential(creds...)
}
