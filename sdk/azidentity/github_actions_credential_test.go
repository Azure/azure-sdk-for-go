// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const githubActionsOIDCHost = "oidc.actions.githubusercontent.com"

func setGitHubActionsEnvironment(t *testing.T) {
	t.Helper()
	setEnvironmentVariables(t, map[string]string{
		azureAuthorityHost:         "",
		azureClientID:              fakeClientID,
		azureTenantID:              fakeTenantID,
		azureRegionalAuthorityName: "",
		actionsIDTokenRequestURL:   "https://" + githubActionsOIDCHost + "/oidc?api-version=2.0&existing=a%2Bb&audience=old&audience=other",
		actionsIDTokenRequestToken: "github-request-token",
	})
}

type gitHubActionsTransport func(*http.Request) (*http.Response, error)

func (g gitHubActionsTransport) Do(r *http.Request) (*http.Response, error) {
	return g(r)
}

func TestGitHubActionsCredential_Configuration(t *testing.T) {
	for _, name := range []string{azureTenantID, azureClientID, actionsIDTokenRequestURL, actionsIDTokenRequestToken} {
		for _, unset := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/unset=%t", name, unset), func(t *testing.T) {
				setGitHubActionsEnvironment(t)
				t.Setenv(name, "")
				if unset {
					require.NoError(t, os.Unsetenv(name))
				}
				cred, err := NewGitHubActionsCredential(nil)
				require.Error(t, err)
				assert.Nil(t, cred)
				assert.Contains(t, err.Error(), name)
			})
		}
	}
	for _, endpoint := range []string{
		":invalid-secret", "/relative-secret", "https:///missing-host", "http://example.com/secret",
		"https://user:secret@example.com", "https://example.com/%zz", "https://example.com?secret=%zz",
	} {
		t.Run(endpoint, func(t *testing.T) {
			setGitHubActionsEnvironment(t)
			t.Setenv(actionsIDTokenRequestURL, endpoint)
			cred, err := NewGitHubActionsCredential(nil)
			require.ErrorContains(t, err, actionsIDTokenRequestURL)
			assert.Nil(t, cred)
			assert.NotContains(t, err.Error(), endpoint)
		})
	}
	t.Run("invalid tenant", func(t *testing.T) {
		setGitHubActionsEnvironment(t)
		t.Setenv(azureTenantID, badTenantID)
		_, err := NewGitHubActionsCredential(nil)
		require.ErrorIs(t, err, errInvalidTenantID)
	})
	for _, authority := range []string{"https://login.example.com", "https://", "http://login.microsoftonline.com", "://invalid"} {
		t.Run(authority, func(t *testing.T) {
			setGitHubActionsEnvironment(t)
			t.Setenv(azureAuthorityHost, authority)
			cred, err := NewGitHubActionsCredential(nil)
			require.Error(t, err)
			assert.Nil(t, cred)
		})
	}
}

func TestGitHubActionsCredential_GetToken(t *testing.T) {
	for _, test := range []struct {
		name, authority, envAuthority, audience string
	}{
		{name: "default", audience: "api://AzureADTokenExchange"},
		{name: "public", authority: cloud.AzurePublic.ActiveDirectoryAuthorityHost, audience: "api://AzureADTokenExchange"},
		{name: "china", authority: cloud.AzureChina.ActiveDirectoryAuthorityHost, audience: "api://AzureADTokenExchangeChina"},
		{name: "US government", authority: cloud.AzureGovernment.ActiveDirectoryAuthorityHost, audience: "api://AzureADTokenExchangeUSGov"},
		{name: "France government", authority: "https://login.sovcloud-identity.fr", audience: "api://AzureADTokenExchangeFrance"},
		{name: "Germany government", authority: "https://login.sovcloud-identity.de", audience: "api://AzureADTokenExchangeGermany"},
		{name: "Singapore government", authority: "https://login.sovcloud-identity.sg", audience: "api://AzureADTokenExchangeGovSG"},
		{name: "environment authority", envAuthority: cloud.AzureChina.ActiveDirectoryAuthorityHost, audience: "api://AzureADTokenExchangeChina"},
		{name: "option overrides environment", authority: cloud.AzureGovernment.ActiveDirectoryAuthorityHost, envAuthority: cloud.AzureChina.ActiveDirectoryAuthorityHost, audience: "api://AzureADTokenExchangeUSGov"},
		{name: "case insensitive", authority: "https://LOGIN.MICROSOFTONLINE.COM/", audience: "api://AzureADTokenExchange"},
	} {
		t.Run(test.name, func(t *testing.T) {
			setGitHubActionsEnvironment(t)
			t.Setenv(azureAuthorityHost, test.envAuthority)
			oidcRequests, tokenRequests := 0, 0
			sts := mockSTS{tokenRequestCallback: func(r *http.Request) *http.Response {
				tokenRequests++
				assert.Equal(t, "github-assertion", r.PostForm.Get("client_assertion"))
				assert.Equal(t, fakeClientID, r.PostForm.Get("client_id"))
				assert.Equal(t, "client_credentials", r.PostForm.Get("grant_type"))
				assert.Contains(t, r.URL.Path, "/"+fakeTenantID+"/")
				return nil
			}}
			transport := gitHubActionsTransport(func(r *http.Request) (*http.Response, error) {
				if r.URL.Hostname() != githubActionsOIDCHost {
					return sts.Do(r)
				}
				oidcRequests++
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/oidc", r.URL.Path)
				assert.Equal(t, "Bearer github-request-token", r.Header.Get("Authorization"))
				assert.Equal(t, []string{test.audience}, r.URL.Query()["audience"])
				assert.Equal(t, "2.0", r.URL.Query().Get("api-version"))
				assert.Equal(t, "a+b", r.URL.Query().Get("existing"))
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"value":"github-assertion"}`))}, nil
			})
			cred, err := NewGitHubActionsCredential(&GitHubActionsCredentialOptions{
				ClientOptions: azcore.ClientOptions{
					Cloud:     cloud.Configuration{ActiveDirectoryAuthorityHost: test.authority},
					Transport: transport,
				},
			})
			require.NoError(t, err)
			// Configuration is captured at construction, not reread for each assertion.
			t.Setenv(actionsIDTokenRequestURL, "invalid")
			t.Setenv(actionsIDTokenRequestToken, "different")
			for range 2 {
				token, err := cred.GetToken(ctx, testTRO)
				require.NoError(t, err)
				assert.Equal(t, tokenValue, token.Token)
				assert.False(t, token.ExpiresOn.IsZero())
			}
			assert.Equal(t, 1, oidcRequests, "cached access tokens shouldn't request a new assertion")
			assert.Equal(t, 1, tokenRequests)
		})
	}
}

func TestGitHubActionsCredential_OIDCError(t *testing.T) {
	for _, test := range []struct {
		name, body string
		status     int
		err        error
	}{
		{name: "unauthorized", status: http.StatusUnauthorized, body: "permission denied"},
		{name: "server error", status: http.StatusInternalServerError, body: "unavailable"},
		{name: "transport error", err: errors.New("connection failed")},
		{name: "malformed JSON", body: `{"value":"secret-assertion"`},
		{name: "wrong value type", body: `{"value":{"secret-assertion":true}}`},
		{name: "missing value", body: `{"other":"secret-assertion"}`},
		{name: "empty value", body: `{"value":""}`},
		{name: "whitespace value", body: `{"value":" \n"}`},
		{name: "null value", body: `{"value":null}`},
		{name: "null response", body: `null`},
		{name: "empty response"},
	} {
		t.Run(test.name, func(t *testing.T) {
			setGitHubActionsEnvironment(t)
			sts := mockSTS{}
			transport := gitHubActionsTransport(func(r *http.Request) (*http.Response, error) {
				if r.URL.Hostname() != githubActionsOIDCHost {
					return sts.Do(r)
				}
				if test.err != nil {
					return nil, test.err
				}
				status := test.status
				if status == 0 {
					status = http.StatusOK
				}
				return &http.Response{
					StatusCode: status,
					Status:     http.StatusText(status),
					Body:       io.NopCloser(strings.NewReader(test.body)),
				}, nil
			})
			cred, err := NewGitHubActionsCredential(&GitHubActionsCredentialOptions{
				ClientOptions: azcore.ClientOptions{Transport: transport, Retry: policy.RetryOptions{MaxRetries: -1}},
			})
			require.NoError(t, err)
			token, err := cred.GetToken(ctx, testTRO)
			require.Error(t, err)
			assert.Empty(t, token.Token)
			var authErr *AuthenticationFailedError
			require.ErrorAs(t, err, &authErr)
			assert.Contains(t, err.Error(), credNameGitHubActions)
			assert.NotContains(t, err.Error(), "secret-assertion")
			assert.NotContains(t, err.Error(), "github-request-token")
			if test.status != 0 {
				require.NotNil(t, authErr.RawResponse)
				assert.Equal(t, test.status, authErr.RawResponse.StatusCode)
			} else {
				assert.Nil(t, authErr.RawResponse)
			}
		})
	}
}

func TestGitHubActionsCredential_Options(t *testing.T) {
	setGitHubActionsEnvironment(t)
	tokenRequests, oidcRequests, discoveryRequests := 0, 0, 0
	sts := mockSTS{tenant: "other-tenant", tokenRequestCallback: func(r *http.Request) *http.Response {
		tokenRequests++
		assert.Contains(t, r.PostForm.Get("claims"), "CP1")
		if tokenRequests == 2 {
			assert.Contains(t, r.PostForm.Get("claims"), "test-claim")
		}
		assert.Contains(t, r.URL.Path, "/other-tenant/")
		return nil
	}}
	options := GitHubActionsCredentialOptions{
		AdditionallyAllowedTenants: []string{"other-tenant"},
		DisableInstanceDiscovery:   true,
		Cache: internal.NewCache(func(bool) (cache.ExportReplace, error) {
			return &testCache{}, nil
		}),
		ClientOptions: azcore.ClientOptions{Transport: gitHubActionsTransport(func(r *http.Request) (*http.Response, error) {
			if strings.Contains(r.URL.Path, "/discovery/instance") {
				discoveryRequests++
			}
			if r.URL.Hostname() == githubActionsOIDCHost {
				oidcRequests++
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"value":"github-assertion"}`))}, nil
			}
			return sts.Do(r)
		})},
	}
	tro := policy.TokenRequestOptions{Scopes: testTRO.Scopes, TenantID: "other-tenant", EnableCAE: true}
	for range 2 {
		cred, err := NewGitHubActionsCredential(&options)
		require.NoError(t, err)
		token, err := cred.GetToken(ctx, tro)
		require.NoError(t, err)
		assert.Equal(t, tokenValue, token.Token)
		_, err = cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: testTRO.Scopes, TenantID: "not-allowed"})
		require.ErrorContains(t, err, "AdditionallyAllowedTenants")
	}
	assert.Equal(t, 0, discoveryRequests)
	assert.Equal(t, 1, tokenRequests, "credential instances should share the persistent cache")
	assert.Equal(t, 1, oidcRequests)

	// Claims challenges must bypass the cached token.
	cred, err := NewGitHubActionsCredential(&options)
	require.NoError(t, err)
	tro.Claims = `{"access_token":{"test-claim":null}}`
	_, err = cred.GetToken(ctx, tro)
	require.NoError(t, err)
	assert.Equal(t, 2, tokenRequests)
	assert.Equal(t, 2, oidcRequests)
}

func TestGitHubActionsCredential_CanceledContext(t *testing.T) {
	setGitHubActionsEnvironment(t)
	cred, err := NewGitHubActionsCredential(&GitHubActionsCredentialOptions{
		ClientOptions: azcore.ClientOptions{Transport: gitHubActionsTransport(func(r *http.Request) (*http.Response, error) {
			return nil, r.Context().Err()
		})},
	})
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	_, err = cred.getAssertion(ctx)
	require.ErrorContains(t, err, context.Canceled.Error())
}
