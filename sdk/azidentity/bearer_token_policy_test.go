// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	accessTokenRespError      = `{"error": "invalid_client","error_description": "Invalid client secret is provided.","error_codes": [0],"timestamp": "2019-12-01 19:00:00Z","trace_id": "2d091b0","correlation_id": "a999","error_uri": "https://login.contoso.com/error?code=0"}`
	accessTokenRespSuccess    = `{"access_token": "` + tokenValue + `", "expires_in": 3600}`
	accessTokenRespMalformed  = `{"access_token": 0, "expires_in": 3600}`
	accessTokenRespShortLived = `{"access_token": "` + tokenValue + `", "expires_in": 0}`
)

func defaultTestPipeline(srv azcore.Transport, cred azcore.Credential, scope string) azcore.Pipeline {
	retryOpts := azcore.RetryOptions{
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    50 * time.Millisecond,
	}
	return azcore.NewPipeline(
		srv,
		azcore.NewRetryPolicy(&retryOpts),
		cred.NewAuthenticationPolicy(azcore.AuthenticationOptions{TokenRequest: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewLogPolicy(nil))
}

func TestBearerPolicy_SuccessGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, scope)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
	const expectedToken = bearerTokenPrefix + tokenValue
	if token := resp.Request.Header.Get(headerAuthorization); token != expectedToken {
		t.Fatalf("expected token '%s', got '%s'", expectedToken, token)
	}
}

func TestBearerPolicy_CredentialFailGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, scope)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pipeline.Do(req)
	var afe *AuthenticationFailedError
	if !errors.As(err, &afe) {
		t.Fatalf("unexpected error type %v", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestBearerTokenPolicy_TokenExpired(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespShortLived)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespShortLived)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespShortLived)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, scope)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

// with https scheme enabled we get an auth failed error which let's us test the is not retriable error
func TestRetryPolicy_NonRetriable(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, scope)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	var afe *AuthenticationFailedError
	if !errors.As(err, &afe) {
		t.Fatalf("unexpected error type %v", err)
	}
}

func TestRetryPolicy_HTTPRequest(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	options := ClientSecretCredentialOptions{}
	options.AuthorityHost = srv.URL()
	options.HTTPClient = srv
	cred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, scope)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	var afe *AuthenticationFailedError
	if !errors.As(err, &afe) {
		t.Fatalf("unexpected error type %v", err)
	}
}

func TestBearerPolicy_GetTokenFailsNoDeadlock(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, &ClientSecretCredentialOptions{
		HTTPClient:    srv,
		AuthorityHost: srv.URL(),
		Retry: azcore.RetryOptions{
			// use a negative try timeout to trigger a deadline exceeded error causing GetToken() to fail
			TryTimeout:    -1 * time.Nanosecond,
			MaxRetryDelay: 500 * time.Millisecond,
			RetryDelay:    50 * time.Millisecond,
			MaxRetries:    3,
		}})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, scope)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := pipeline.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestBearerTokenWithAuxiliaryTenants(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	headerResult := "Bearer new_token, Bearer new_token, Bearer new_token"
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse()
	options := ClientSecretCredentialOptions{
		AuthorityHost: srv.URL(),
		HTTPClient:    srv,
	}
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	retryOpts := azcore.RetryOptions{
		MaxRetryDelay: 500 * time.Millisecond,
		RetryDelay:    50 * time.Millisecond,
	}
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewRetryPolicy(&retryOpts),
		cred.NewAuthenticationPolicy(
			azcore.AuthenticationOptions{
				TokenRequest: azcore.TokenRequestOptions{
					Scopes: []string{scope},
				},
				AuxiliaryTenants: []string{"tenant1", "tenant2", "tenant3"},
			}),
		azcore.NewLogPolicy(nil))

	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if auxH := resp.Request.Header.Get(headerAuxiliaryAuthorization); auxH != headerResult {
		t.Fatalf("unexpected auxiliary authorization header %s", auxH)
	}
}
