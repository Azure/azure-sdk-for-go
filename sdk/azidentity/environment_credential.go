// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// NewEnvironmentCredential creates an instance that implements the azcore.TokenCredential interface and reads credential details from environment variables.
// If the expected environment variables are not found at this time, then a CredentialUnavailableError will be returned.
// options: The options used to configure the management of the requests sent to Azure Active Directory.
func NewEnvironmentCredential(options *TokenCredentialOptions) (azcore.TokenCredential, error) {
	tenantID := os.Getenv("AZURE_TENANT_ID")
	if tenantID == "" {
		err := &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_TENANT_ID"}
		logCredentialError(err.CredentialType, err)
		return nil, err
	}
	clientID := os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		err := &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_CLIENT_ID"}
		logCredentialError(err.CredentialType, err)
		return nil, err
	}
	if clientSecret := os.Getenv("AZURE_CLIENT_SECRET"); clientSecret != "" {
		azcore.Log().Write(LogCredential, "Azure Identity => NewEnvironmentCredential() invoking ClientSecretCredential")
		return NewClientSecretCredential(tenantID, clientID, clientSecret, options)
	}
	if clientCertificate := os.Getenv("AZURE_CLIENT_CERTIFICATE_PATH"); clientCertificate != "" {
		azcore.Log().Write(LogCredential, "Azure Identity => NewEnvironmentCredential() invoking ClientCertificateCredential")
		return NewClientCertificateCredential(tenantID, clientID, clientCertificate, nil, options)
	}
	if username := os.Getenv("AZURE_USERNAME"); username != "" {
		if password := os.Getenv("AZURE_PASSWORD"); password != "" {
			azcore.Log().Write(LogCredential, "Azure Identity => NewEnvironmentCredential() invoking UsernamePasswordCredential")
			return NewUsernamePasswordCredential(tenantID, clientID, username, password, options)
		}
	}
	err := &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_CLIENT_SECRET or AZURE_CLIENT_CERTIFICATE_PATH or AZURE_USERNAME and AZURE_PASSWORD"}
	logCredentialError(err.CredentialType, err)
	return nil, err

}
