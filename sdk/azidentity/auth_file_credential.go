// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// AuthFileCredential enables authentication to Azure Active Directory using configuration information stored Azure SDK Auth File.
type AuthFileCredential struct {
	filePath        string
	credential      azcore.TokenCredential
	authfileoptions TokenCredentialOptions
}

// NewAuthFileCredential creates an instance of the AuthFileCredential class based on information specified in the SDK Authentication file.
func NewAuthFileCredential(filePath string, options *TokenCredentialOptions) (*AuthFileCredential, error) {
	options, err := options.setDefaultValues()
	if err != nil {
		return nil, err
	}
	cred, err := NewClientSecretCredentialFromAuthenticationFile(filePath, options)
	if err != nil {
		return nil, &AuthenticationFailedError{msg: "Error parsing SDK Auth File", inner: err}
	}
	return &AuthFileCredential{filePath: filePath, credential: cred, authfileoptions: *options}, nil
}

// GetToken obtains a token from Azure Active Directory using the client specified in the SDK Authentication file.
func (c *AuthFileCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return c.credential.GetToken(ctx, opts)
}

// AuthenticationPolicy implements the azcore.Credential interface and returns the policy associated with the AuthFileCredential type.
func (c *AuthFileCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}

// NewClientSecretCredentialFromAuthenticationFile creates a ClientSecrentCredential initialized from the specified Authentication file.
func NewClientSecretCredentialFromAuthenticationFile(filePath string, options *TokenCredentialOptions) (*ClientSecretCredential, error) {
	// Read the SDK Authentication file into memory
	authData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Build the ClientSecretCredential from the Authentication File data
	token := struct {
		clientID                       string `json:"clientId"`
		clientSecret                   string `json:"clientSecret"`
		subscriptionID                 string `json:"subscriptionId"`
		tenantID                       string `json:"tenantId"`
		activeDirectoryEndpointURL     string `json:"activeDirectoryEndpointUrl"`
		resourceManagerEndpointURL     string `json:"resourceManagerEndpointUrl"`
		activeDirectoryGraphResourceID string `json:"activeDirectoryGraphResourceId"`
		sqlManagementEndpointURL       string `json:"sqlManagementEndpointUrl"`
		galleryEndpointURL             string `json:"galleryEndpointUrl"`
		managementEndpointURL          string `json:"managementEndpointUrl"`
	}{}

	if err = json.Unmarshal(authData, &token); err != nil {
		return nil, err
	}

	// Parse string activeDirectoryEndpointUrl to a Url.
	var tco = *options // Make a copy of the passed-inoption to not change the customer's options
	tco.AuthorityHost, err = url.Parse(token.activeDirectoryEndpointURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ActiveDirectoryEndpointUrl in auth file: %w", err)
	}
	return NewClientSecretCredential(token.tenantID, token.clientID, token.clientSecret, &tco)
}
