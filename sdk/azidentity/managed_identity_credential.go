// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const defaultSuffix = "/.default"

// ManagedIdentityCredentialOptions contains parameters that can be used to configure a Managed Identity Credential
type ManagedIdentityCredentialOptions struct {
	// HTTPClient sets the transport for making HTTP requests.
	// Leave this as nil to use the default HTTP transport.
	HTTPClient azcore.Transport

	// LogOptions configures the built-in request logging policy behavior.
	LogOptions azcore.RequestLogOptions

	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry azcore.TelemetryOptions
}

func (m *ManagedIdentityCredentialOptions) setDefaultValues() *ManagedIdentityCredentialOptions {
	if m == nil {
		m = defaultMSIOpts
	}
	return m
}

// ManagedIdentityCredential attempts authentication using a managed identity that has been assigned to the deployment environment. This authentication type works in Azure VMs,
// App Service and Azure Functions applications, as well as inside of Azure Cloud Shell. More information about configuring managed identities can be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview
type ManagedIdentityCredential struct {
	clientID string
	client   *managedIdentityClient
}

// NewManagedIdentityCredential creates an instance of the ManagedIdentityCredential capable of authenticating a resource with a managed identity.
// clientID: The client id to authenticate for a user assigned managed identity.  More information on user assigned managed identities cam be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview#how-a-user-assigned-managed-identity-works-with-an-azure-vm
// options: Options that allow to configure the management of the requests sent to the Azure Active Directory service.
func NewManagedIdentityCredential(clientID string, options *ManagedIdentityCredentialOptions) (*ManagedIdentityCredential, error) {
	// Create a new Managed Identity Client with default options
	client := newManagedIdentityClient(options)
	// Create a context that will timeout after 500 milliseconds (that is the amount of time designated to find out if the IMDS endpoint is available)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(client.imdsAvailableTimeoutMS)*time.Millisecond)
	defer cancelFunc()
	msiType, err := client.getMSIType(ctx)
	// If there is an error that means that the code is not running in a Managed Identity environment
	if err != nil {
		return nil, &CredentialUnavailableError{CredentialType: "Managed Identity Credential", Message: "Please make sure you are running in a managed identity environment, such as a VM, Azure Functions, Cloud Shell, etc..."}
	}
	// Assign the msiType discovered onto the client
	client.msiType = msiType
	return &ManagedIdentityCredential{clientID: clientID, client: client}, nil
}

// GetToken obtains an AccessToken from the Managed Identity service if available.
// scopes: The list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls, or a default AccessToken if no managed identity is available.
func (c *ManagedIdentityCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, c.clientID, opts.Scopes)
}

// AuthenticationPolicy implements the azcore.Credential interface on ManagedIdentityCredential.
// Please note: the TokenRequestOptions included in AuthenticationPolicyOptions must be a slice of resources in this case and not scopes
func (c *ManagedIdentityCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	// The following code will remove the /.default suffix from any scopes passed into the method since ManagedIdentityCredentials expect a resource string instead of a scope string
	for i, s := range options.Options.Scopes {
		if strings.HasSuffix(s, defaultSuffix) {
			options.Options.Scopes[i] = s[:len(s)-len(defaultSuffix)]
		}
	}
	return newBearerTokenPolicy(c, options)
}
