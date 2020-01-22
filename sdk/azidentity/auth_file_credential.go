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

type authFileAccesstoken struct {
	ClientId                       string `json:"clientId"`
	ClientSecret                   string `json:"clientSecret"`
	SubscriptionId                 string `json:"subscriptionId"`
	TenantId                       string `json:"tenantId"`
	ActiveDirectoryEndpointUrl     string `json:"activeDirectoryEndpointUrl"`
	ResourceManagerEndpointUrl     string `json:"resourceManagerEndpointUrl"`
	ActiveDirectoryGraphResourceId string `json:"activeDirectoryGraphResourceId"`
	SqlManagementEndpointUrl       string `json:"sqlManagementEndpointUrl"`
	GalleryEndpointUrl             string `json:"galleryEndpointUrl"`
	ManagementEndpointUrl          string `json:"managementEndpointUrl"`
}

// Enables authentication to Azure Active Directory using configuration information stored Azure SDK Auth File.
type AuthFileCredential struct {
	filePath        string
	credential      azcore.TokenCredential
	authfileoptions *TokenCredentialOptions
}

// Creates an instance of the SdkAuthFileCredential class based on information in given SDK Auth file.
func NewAuthFileCredential(filePath string, options *TokenCredentialOptions) (*AuthFileCredential, error) {
	if filePath == "" {
		return nil, fmt.Errorf("The parameter 'filePath' must set a value")
	}

	return &AuthFileCredential{filePath: filePath, authfileoptions: options}, nil
}

// Obtains a token from the Azure Active Directory service, using the specified client detailed specified in the SDK Auth file.
func (c *AuthFileCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	_, err := c.ensureCredential()
	if err != nil {
		return nil, &AuthenticationFailedError{msg: "Error parsing SDK Auth File", inner: err}
	}

	return c.credential.GetToken(ctx, opts)
}

// AuthenticationPolicy implements the azcore.Credential interface on AuthFileCredential.
func (c *AuthFileCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}

// Ensures that credential information is loaded from the SDK Auth file. This method should be called to initialize.
func (c *AuthFileCredential) ensureCredential() (azcore.TokenCredential, error) {
	if c.credential != nil {
		return c.credential, nil
	}

	authData, err := c.parseCredentialsFile(c.filePath)
	if err != nil {
		return nil, err
	}

	c.credential, err = c.buildCredentialForCredentialsFile(authData)
	if err != nil {
		return nil, err
	}

	return c.credential, nil
}

// Parse credential file from local file path to byte.
func (c *AuthFileCredential) parseCredentialsFile(filePath string) ([]byte, error) {
	authData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return authData, nil
}

// Use the SDK auth file info to build a ClientSecretCredential
func (c *AuthFileCredential) buildCredentialForCredentialsFile(authData []byte) (*ClientSecretCredential, error) {
	token := &authFileAccesstoken{}

	err := json.Unmarshal(authData, token)
	if err != nil {
		return nil, err
	}

	clientId := token.ClientId
	clientSecret := token.ClientSecret
	tenantId := token.TenantId
	activeDirectoryEndpointUrl := token.ActiveDirectoryEndpointUrl

	if clientId == "" || clientSecret == "" || tenantId == "" || activeDirectoryEndpointUrl == "" {
		return nil, fmt.Errorf("Malformed Azure SDK Auth file. The file should contain 'clientId', 'clientSecret', 'tenantId' and 'activeDirectoryEndpointUrl' values.")
	}

	// Parse string activeDirectoryEndpointUrl to a Url.
	AuthorityHost, err := url.Parse(activeDirectoryEndpointUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse activeDirectoryEndpointUrl in auth file: %w", err)
	}
	c.authfileoptions.AuthorityHost = AuthorityHost

	return NewClientSecretCredential(tenantId, clientId, clientSecret, c.authfileoptions)
}
