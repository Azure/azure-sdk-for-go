// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestChainedTokenCredential_InstantiateSuccess(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	secCred, err := NewClientSecretCredential(tenantID, clientID, secret, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	envCred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("Could not find appropriate environment credentials")
	}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{secCred, envCred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred != nil {
		if len(cred.sources) != 2 {
			t.Fatalf("Expected 2 sources in the chained token credential, instead found %d", len(cred.sources))
		}
	}
}

func TestChainedTokenCredential_InstantiateFailure(t *testing.T) {
	secCred, err := NewClientSecretCredential(tenantID, clientID, secret, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = NewChainedTokenCredential([]azcore.TokenCredential{secCred, nil}, nil)
	if err == nil {
		t.Fatalf("Expected an error for sending a nil credential in the chain")
	}
	_, err = NewChainedTokenCredential([]azcore.TokenCredential{}, nil)
	if err == nil {
		t.Fatalf("Expected an error for not sending any credential sources")
	}
}

func TestChainedTokenCredential_GetTokenSuccess(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	secCred, err := NewClientSecretCredential(tenantID, clientID, secret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	envCred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{
		ClientOptions: azcore.ClientOptions{Transport: srv},
		AuthorityHost: AuthorityHost(srv.URL()),
	})
	if err != nil {
		t.Fatalf("Failed to create environment credential: %v", err)
	}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{secCred, envCred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none")
	}
	if tk.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if tk.ExpiresOn.IsZero() {
		t.Fatalf("Received an incorrect time in the response")
	}
}

func TestChainedTokenCredential_GetTokenFail(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	secCred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{secCred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	var authErr AuthenticationFailedError
	if !errors.As(err, &authErr) {
		t.Fatalf("Expected AuthenticationFailedError, received %T", err)
	}
	if len(err.Error()) == 0 {
		t.Fatalf("Did not create an appropriate error message")
	}
}

func TestChainedTokenCredential_GetTokenWithUnavailableCredentialInChain(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendError(newCredentialUnavailableError("MockCredential", "Mocking a credential unavailable error"))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	secCred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	// The chain has the same credential twice, since it doesn't matter what credential we add to the chain as long as it is not a nil credential.
	// Most credentials will not be instantiated if the conditions do not exist to allow them to be used, thus returning a
	// CredentialUnavailable error from the constructor. In order to test the CredentialUnavailable functionality for
	// ChainedTokenCredential we have to mock with two valid credentials, but the first will fail since the first response queued
	// in the test server is a CredentialUnavailable error.
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{secCred, secCred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none")
	}
	if tk.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if tk.ExpiresOn.IsZero() {
		t.Fatalf("Received an incorrect time in the response")
	}
}

func TestChainedTokenCredential_ChecksThatSuccessfulCredentialIsSet(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	secCred, err := NewClientSecretCredential(tenantID, clientID, secret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	envCred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{
		ClientOptions: azcore.ClientOptions{Transport: srv},
		AuthorityHost: AuthorityHost(srv.URL()),
	})
	if err != nil {
		t.Fatalf("Failed to create environment credential: %v", err)
	}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{secCred, envCred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none")
	}
	if tk.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if tk.ExpiresOn.IsZero() {
		t.Fatalf("Received an incorrect time in the response")
	}
	if cred.successfulCredential == nil {
		t.Fatalf("The successful credential was not assigned")
	}
	if cred.successfulCredential != secCred {
		t.Fatalf("The successful credential should have been the secret credential")
	}
}

/**
 * Helps count the number of times a credential is called.
 */
type TestCountPolicy struct{ count int }

/**
 * Helps count the number of times a credential is called.
 */
func (p *TestCountPolicy) Do(req *policy.Request) (*http.Response, error) {
	p.count += 1
	return req.Next()
}

func TestChainedTokenCredential_RepeatedGetTokenWithSuccessfulCredential(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	srv, close := mock.NewTLSServer()
	defer close()

	secretCountPolicy := &TestCountPolicy{}
	environmentCountPolicy := &TestCountPolicy{}

	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.PerCallPolicies = []policy.Policy{secretCountPolicy}
	options.Transport = srv
	secCred, err := NewClientSecretCredential(tenantID, clientID, secret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	envCred, err := NewEnvironmentCredential(&EnvironmentCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Transport:       srv,
			PerCallPolicies: []policy.Policy{environmentCountPolicy},
		},
		AuthorityHost: AuthorityHost(srv.URL()),
	})
	if err != nil {
		t.Fatalf("Failed to create environment credential: %v", err)
	}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{secCred, envCred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none")
	}
	if tk.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if tk.ExpiresOn.IsZero() {
		t.Fatalf("Received an incorrect time in the response")
	}
	if cred.successfulCredential == nil {
		t.Fatalf("The successful credential was not assigned")
	}
	if cred.successfulCredential != secCred {
		t.Fatalf("The successful credential should have been the secret credential")
	}
	if secretCountPolicy.count != 1 {
		t.Fatalf("The secret credential policies should have been triggered once")
	}
	if environmentCountPolicy.count != 0 {
		t.Fatalf("The environment credential policies should not have been triggered")
	}
	tk2, err2 := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err2 != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none. Error: %v", err2)
	}
	if tk2.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if tk2.ExpiresOn.IsZero() {
		t.Fatalf("Received an incorrect time in the response")
	}
	if secretCountPolicy.count != 2 {
		t.Fatalf("The secret credential policies should have been triggered twice")
	}
	if environmentCountPolicy.count != 0 {
		t.Fatalf("The environment credential policies should not have been triggered")
	}
}

// A credential that always throws a CredentialUnavailableError
type UnavailableCredential struct {
	callCount int
}

func NewUnavailableCredential() (*UnavailableCredential, error) {
	return &UnavailableCredential{callCount: 0}, nil
}
func (c *UnavailableCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (token *azcore.AccessToken, err error) {
	c.callCount += 1
	return nil, newCredentialUnavailableError("UnavailableCredential", "Expected CredentialUnavailableError")
}

func TestChainedTokenCredential_RepeatedGetTokenWithSuccessfulCredentialWithRetrySources(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	srv, close := mock.NewTLSServer()
	defer close()

	secretCountPolicy := &TestCountPolicy{}

	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))

	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.PerCallPolicies = []policy.Policy{secretCountPolicy}
	options.Transport = srv

	unavailableCred, _ := NewUnavailableCredential()
	secCred, err := NewClientSecretCredential(tenantID, clientID, secret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}

	// Backwards order: envCred first, secCred later, to check that envCred is always called when RetrySources is set to true.
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{unavailableCred, secCred}, &ChainedTokenCredentialOptions{RetrySources: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none. Error: %v", err)
	}
	if tk.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if tk.ExpiresOn.IsZero() {
		t.Fatalf("Received an incorrect time in the response")
	}
	if cred.successfulCredential == nil {
		t.Fatalf("The successful credential was not assigned")
	}
	if cred.successfulCredential != secCred {
		t.Fatalf("The successful credential should have been the secret credential")
	}
	if secretCountPolicy.count != 1 {
		t.Fatalf("The secret credential policies should have been triggered once")
	}
	if unavailableCred.callCount != 1 {
		t.Fatalf("The environment credential policies should have been triggered once")
	}
	tk2, err2 := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err2 != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none. Error: %v", err2)
	}
	if tk2.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if tk2.ExpiresOn.IsZero() {
		t.Fatalf("Received an incorrect time in the response")
	}
	if secretCountPolicy.count != 2 {
		t.Fatalf("The secret credential policies should have been triggered twice")
	}
	if unavailableCred.callCount != 2 {
		t.Fatalf("The environment credential policies should have been triggered twice")
	}
}

func TestBearerPolicy_ChainedTokenCredential(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to initialize environment variables. Received: %v", err)
	}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.Transport = srv
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	chainedCred, err := NewChainedTokenCredential([]azcore.TokenCredential{cred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pipeline := defaultTestPipeline(srv, chainedCred, scope)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected an empty error but receive: %v", err)
	}
}
