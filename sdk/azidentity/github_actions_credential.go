// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	credNameGitHubActions      = "GitHubActionsCredential"
	actionsIDTokenRequestURL   = "ACTIONS_ID_TOKEN_REQUEST_URL"
	actionsIDTokenRequestToken = "ACTIONS_ID_TOKEN_REQUEST_TOKEN"
)

// GitHubActionsCredential authenticates with workload identity federation in a GitHub Actions workflow.
// The workflow must have the id-token: write permission and a federated identity credential configured
// for its Microsoft Entra application or user-assigned managed identity. See [GitHub documentation]
// for more information.
//
// [GitHub documentation]: https://docs.github.com/actions/how-tos/secure-your-work/security-harden-deployments/oidc-in-azure
type GitHubActionsCredential struct {
	cred                  *ClientAssertionCredential
	oidcURL, requestToken string
}

// GitHubActionsCredentialOptions contains optional parameters for GitHubActionsCredential.
type GitHubActionsCredentialOptions struct {
	azcore.ClientOptions

	// AdditionallyAllowedTenants specifies additional tenants for which the credential may acquire tokens.
	// Add the wildcard value "*" to allow the credential to acquire tokens for any tenant in which the
	// application is registered.
	AdditionallyAllowedTenants []string

	// Cache is a persistent cache the credential will use to store the tokens it acquires, making
	// them available to other processes and credential instances. The default, zero value means the
	// credential will store tokens in memory and not share them with any other credential instance.
	Cache Cache

	// DisableInstanceDiscovery determines whether the credential requests Microsoft Entra instance metadata
	// from https://login.microsoft.com before authenticating. Setting this to true will skip this request,
	// making the application responsible for ensuring the configured authority is valid and trustworthy.
	// The authority must still be one of the Azure clouds supported by GitHubActionsCredential.
	DisableInstanceDiscovery bool
}

// NewGitHubActionsCredential constructs a GitHubActionsCredential. It reads AZURE_TENANT_ID and
// AZURE_CLIENT_ID for the identity to authenticate, and ACTIONS_ID_TOKEN_REQUEST_URL and
// ACTIONS_ID_TOKEN_REQUEST_TOKEN as provided by GitHub Actions. All four variables must have nonempty
// values. Pass nil for options to accept defaults.
//
// The credential derives the federation audience from the authority host configured by options.Cloud,
// AZURE_AUTHORITY_HOST, or the Azure Public default, in that order. Supported clouds are Azure Public,
// Azure China, Azure US Government, France Government, Germany Government, and Singapore Government.
func NewGitHubActionsCredential(options *GitHubActionsCredentialOptions) (*GitHubActionsCredential, error) {
	if options == nil {
		options = &GitHubActionsCredentialOptions{}
	}
	tenantID := os.Getenv(azureTenantID)
	if tenantID == "" {
		return nil, fmt.Errorf("no tenant ID specified. Set %s", azureTenantID)
	}
	if !validTenantID(tenantID) {
		return nil, errInvalidTenantID
	}
	clientID := os.Getenv(azureClientID)
	if clientID == "" {
		return nil, fmt.Errorf("no client ID specified. Set %s", azureClientID)
	}
	endpoint := os.Getenv(actionsIDTokenRequestURL)
	requestToken := os.Getenv(actionsIDTokenRequestToken)
	for _, v := range []struct{ name, value string }{
		{actionsIDTokenRequestURL, endpoint},
		{actionsIDTokenRequestToken, requestToken},
	} {
		if v.value == "" {
			return nil, fmt.Errorf("no value for environment variable %s. Run in a GitHub Actions workflow with permissions: id-token: write", v.name)
		}
	}
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil {
		return nil, errors.New("invalid URL for " + actionsIDTokenRequestURL + ": expected an absolute HTTPS URL without user information")
	}
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, errors.New("invalid query parameters in " + actionsIDTokenRequestURL)
	}

	g := GitHubActionsCredential{requestToken: requestToken}
	cred, err := NewClientAssertionCredential(tenantID, clientID, g.getAssertion, &ClientAssertionCredentialOptions{
		AdditionallyAllowedTenants: options.AdditionallyAllowedTenants,
		Cache:                      options.Cache,
		ClientOptions:              options.ClientOptions,
		DisableInstanceDiscovery:   options.DisableInstanceDiscovery,
	})
	if err != nil {
		return nil, err
	}
	authority, err := url.Parse(cred.client.host)
	if err != nil {
		return nil, err
	}
	var audience string
	switch strings.ToLower(authority.Hostname()) {
	case "login.microsoftonline.com":
		audience = "api://AzureADTokenExchange"
	case "login.chinacloudapi.cn":
		audience = "api://AzureADTokenExchangeChina"
	case "login.microsoftonline.us":
		audience = "api://AzureADTokenExchangeUSGov"
	case "login.sovcloud-identity.fr":
		audience = "api://AzureADTokenExchangeFrance"
	case "login.sovcloud-identity.de":
		audience = "api://AzureADTokenExchangeGermany"
	case "login.sovcloud-identity.sg":
		audience = "api://AzureADTokenExchangeGovSG"
	default:
		return nil, fmt.Errorf("%s doesn't support authority host %q", credNameGitHubActions, cred.client.host)
	}
	query.Set("audience", audience)
	u.RawQuery = query.Encode()
	g.oidcURL = u.String()
	cred.client.name = credNameGitHubActions
	g.cred = cred
	return &g, nil
}

// GetToken requests an access token from Microsoft Entra ID. Azure SDK clients call this method automatically.
func (g *GitHubActionsCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, credNameGitHubActions+"."+traceOpGetToken, g.cred.client.azClient.Tracer(), nil)
	defer func() { endSpan(err) }()
	tk, err := g.cred.GetToken(ctx, opts)
	return tk, err
}

func (g *GitHubActionsCredential) getAssertion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.oidcURL, nil)
	if err != nil {
		return "", newAuthenticationFailedError(credNameGitHubActions, "couldn't create OIDC token request: "+err.Error(), nil)
	}
	req.Header.Set("Authorization", "Bearer "+g.requestToken)
	res, err := doForClient(g.cred.client.azClient, req)
	if err != nil {
		return "", newAuthenticationFailedError(credNameGitHubActions, "couldn't send OIDC token request: "+err.Error(), nil)
	}
	if res.StatusCode != http.StatusOK {
		return "", newAuthenticationFailedError(credNameGitHubActions, "unexpected response from the GitHub Actions OIDC endpoint. Check workflow permissions and configuration", res)
	}
	b, err := runtime.Payload(res)
	if err != nil {
		return "", newAuthenticationFailedError(credNameGitHubActions, "couldn't read OIDC response content: "+err.Error(), nil)
	}
	var r struct {
		Value string `json:"value"`
	}
	// Don't include successful response content or JSON errors, which may contain the assertion.
	if err = json.Unmarshal(b, &r); err != nil || strings.TrimSpace(r.Value) == "" {
		return "", newAuthenticationFailedError(credNameGitHubActions, "the GitHub Actions OIDC endpoint returned an invalid token response", nil)
	}
	return r.Value, nil
}

var _ azcore.TokenCredential = (*GitHubActionsCredential)(nil)
