// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// NewEnvironmentCredential creates an instance of the ClientSecretCredential type and reads client secret details from environment variables.
// options: The options used to configure the management of the requests sent to the Azure Active Directory service.
func NewEnvironmentCredential(options *TokenCredentialOptions) (*ClientSecretCredential, error) {
	tenantID := os.Getenv("AZURE_TENANT_ID")
	log := azcore.Log()
	if tenantID == "" {
		err := &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_TENANT_ID"}
		msg := fmt.Sprintf("Azure Identity => ERROR in NewEnvironmentCredential() call: %s", err.Error())
		log.Write(azcore.LogError, msg)
		return nil, err
	}

	clientID := os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		err := &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_CLIENT_ID"}
		msg := fmt.Sprintf("Azure Identity => ERROR in NewEnvironmentCredential() call: %s", err.Error())
		log.Write(azcore.LogError, msg)
		return nil, err
	}

	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if clientSecret == "" {
		err := &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_CLIENT_SECRET"}
		msg := fmt.Sprintf("Azure Identity => ERROR in NewEnvironmentCredential() call: %s", err.Error())
		log.Write(azcore.LogError, msg)
		return nil, err
	}
	log.Write(LogCredential, "Azure Identity => NewEnvironmentCredential() invoking ClientSecretCredential")
	return NewClientSecretCredential(tenantID, clientID, clientSecret, options)
}
