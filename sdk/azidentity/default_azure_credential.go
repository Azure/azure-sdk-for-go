// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

const (
	organizationsTenantID   = "organizations"
	developerSignOnClientID = "04b07795-8ddb-461a-bbee-02f9e1bf7b46"
)

// DefaultAzureCredentialOptions contains options for configuring how credentials are acquired.
type DefaultAzureCredentialOptions struct{}

// NewDefaultAzureCredential provides a default ChainedTokenCredential configuration for applications that will be deployed to Azure.  The following credential
// types will be tried, in the following order:
// - EnvironmentCredential
// - ManagedIdentityCredential
// - AzureCLICredential
// Consult the documentation for these credential types for more information on how they attempt authentication.
func NewDefaultAzureCredential(options *DefaultAzureCredentialOptions) (*ChainedTokenCredential, error) {
	var creds []azcore.TokenCredential
	errMsg := ""

	if options == nil {
		options = &DefaultAzureCredentialOptions{}
	}

	envCred, err := NewEnvironmentCredential(nil)
	if err == nil {
		creds = append(creds, envCred)
	} else {
		errMsg += err.Error()
	}

	msiCred, err := NewManagedIdentityCredential("", nil)
	if err == nil {
		creds = append(creds, msiCred)
	} else {
		errMsg += err.Error()
	}

	cliCred, err := NewAzureCLICredential(nil)
	if err == nil {
		creds = append(creds, cliCred)
	} else {
		errMsg += err.Error()
	}

	// if no credentials are added to the slice of TokenCredentials then return a CredentialUnavailableError
	if len(creds) == 0 {
		err := &CredentialUnavailableError{credentialType: "Default Azure Credential", message: errMsg}
		logCredentialError(err.credentialType, err)
		return nil, err
	}
	log.Write(LogCredential, "Azure Identity => NewDefaultAzureCredential() invoking NewChainedTokenCredential()")
	return NewChainedTokenCredential(creds, nil)
}
