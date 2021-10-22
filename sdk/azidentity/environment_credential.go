// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// EnvironmentCredentialOptions configures the EnvironmentCredential with optional parameters.
// All zero-value fields will be initialized with their default values.
type EnvironmentCredentialOptions struct {
	azcore.ClientOptions

	// The host of the Azure Active Directory authority. The default is AzurePublicCloud.
	// Leave empty to allow overriding the value from the AZURE_AUTHORITY_HOST environment variable.
	AuthorityHost AuthorityHost
}

// EnvironmentCredential enables authentication to Azure Active Directory using either ClientSecretCredential, ClientCertificateCredential or UsernamePasswordCredential.
// This credential type will check for the following environment variables in the same order as listed:
// - AZURE_TENANT_ID
// - AZURE_CLIENT_ID
// - AZURE_CLIENT_SECRET
// - AZURE_CLIENT_CERTIFICATE_PATH
// - AZURE_USERNAME
// - AZURE_PASSWORD
// NOTE: EnvironmentCredential will stop checking environment variables as soon as it finds enough environment variables to
// create a credential type.
type EnvironmentCredential struct {
	cred azcore.TokenCredential
}

// NewEnvironmentCredential creates an instance that implements the azcore.TokenCredential interface and reads credential details from environment variables.
// If the expected environment variables are not found at this time, then a CredentialUnavailableError will be returned.
// options: The options used to configure the management of the requests sent to Azure Active Directory.
func NewEnvironmentCredential(options *EnvironmentCredentialOptions) (*EnvironmentCredential, error) {
	cp := EnvironmentCredentialOptions{}
	if options != nil {
		cp = *options
	}
	tenantID := os.Getenv("AZURE_TENANT_ID")
	if tenantID == "" {
		err := &CredentialUnavailableError{credentialType: "Environment Credential", message: "Missing environment variable AZURE_TENANT_ID"}
		logCredentialError(err.credentialType, err)
		return nil, err
	}
	clientID := os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		err := &CredentialUnavailableError{credentialType: "Environment Credential", message: "Missing environment variable AZURE_CLIENT_ID"}
		logCredentialError(err.credentialType, err)
		return nil, err
	}
	if clientSecret := os.Getenv("AZURE_CLIENT_SECRET"); clientSecret != "" {
		log.Write(EventCredential, "Azure Identity => NewEnvironmentCredential() invoking ClientSecretCredential")
		cred, err := NewClientSecretCredential(tenantID, clientID, clientSecret, &ClientSecretCredentialOptions{AuthorityHost: cp.AuthorityHost, ClientOptions: cp.ClientOptions})
		if err != nil {
			return nil, err
		}
		return &EnvironmentCredential{cred: cred}, nil
	}
	if certPath := os.Getenv("AZURE_CLIENT_CERTIFICATE_PATH"); certPath != "" {
		log.Write(EventCredential, "Azure Identity => NewEnvironmentCredential() invoking ClientCertificateCredential")
		certData, err := os.ReadFile(certPath)
		if err != nil {
			return nil, &CredentialUnavailableError{credentialType: "Environment Credential", message: "Failed to read certificate file: " + err.Error()}
		}
		cred, err := NewClientCertificateCredential(tenantID, clientID, certData, &ClientCertificateCredentialOptions{AuthorityHost: cp.AuthorityHost, ClientOptions: cp.ClientOptions})
		if err != nil {
			return nil, err
		}
		return &EnvironmentCredential{cred: cred}, nil
	}
	if username := os.Getenv("AZURE_USERNAME"); username != "" {
		if password := os.Getenv("AZURE_PASSWORD"); password != "" {
			log.Write(EventCredential, "Azure Identity => NewEnvironmentCredential() invoking UsernamePasswordCredential")
			cred, err := NewUsernamePasswordCredential(tenantID, clientID, username, password, &UsernamePasswordCredentialOptions{AuthorityHost: cp.AuthorityHost, ClientOptions: cp.ClientOptions})
			if err != nil {
				return nil, err
			}
			return &EnvironmentCredential{cred: cred}, nil
		}
	}
	err := &CredentialUnavailableError{credentialType: "Environment Credential", message: "Missing environment variable AZURE_CLIENT_SECRET or AZURE_CLIENT_CERTIFICATE_PATH or AZURE_USERNAME and AZURE_PASSWORD"}
	logCredentialError(err.credentialType, err)
	return nil, err
}

// GetToken obtains a token from Azure Active Directory, using the underlying credential's GetToken method.
// ctx: Context used to control the request lifetime.
// opts: TokenRequestOptions contains the list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *EnvironmentCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return c.cred.GetToken(ctx, opts)
}

var _ azcore.TokenCredential = (*EnvironmentCredential)(nil)
