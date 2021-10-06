// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

const (
	organizationsTenantID   = "organizations"
	developerSignOnClientID = "04b07795-8ddb-461a-bbee-02f9e1bf7b46"
)

// DefaultAzureCredentialOptions contains options for configuring authentication. These options
// may not apply to all credentials in the default chain.
type DefaultAzureCredentialOptions struct {
	// The host of the Azure Active Directory authority. The default is AzurePublicCloud.
	// Leave empty to allow overriding the value from the AZURE_AUTHORITY_HOST environment variable.
	AuthorityHost AuthorityHost
	// HTTPClient sets the transport for making HTTP requests
	// Leave this as nil to use the default HTTP transport
	HTTPClient policy.Transporter
	// Retry configures the built-in retry policy behavior
	Retry policy.RetryOptions
	// Telemetry configures the built-in telemetry policy behavior
	Telemetry policy.TelemetryOptions
	// Logging configures the built-in logging policy behavior.
	Logging policy.LogOptions
}

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

	envCred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{AuthorityHost: options.AuthorityHost,
		HTTPClient: options.HTTPClient,
		Logging:    options.Logging,
		Retry:      options.Retry,
		Telemetry:  options.Telemetry,
	})
	if err == nil {
		creds = append(creds, envCred)
	} else {
		errMsg += err.Error()
	}

	msiCred, err := NewManagedIdentityCredential(&ManagedIdentityCredentialOptions{HTTPClient: options.HTTPClient,
		Logging:   options.Logging,
		Telemetry: options.Telemetry,
	})
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
