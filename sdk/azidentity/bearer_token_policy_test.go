// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	accessTokenRespError      = `{"error": "invalid_client","error_description": "Invalid client secret is provided.","error_codes": [7000215],"timestamp": "2019-12-01 19:00:00Z","trace_id": "2d091b0","correlation_id": "a999","error_uri": "https://login.contoso.com/error?code=7000215"}`
	accessTokenRespSuccess    = `{"access_token": "new_token", "expires_in": 3600}`
	accessTokenRespMalformed  = `{"access_token": 0, "expires_in": 3600}`
	accessTokenRespShortLived = `{"access_token": "new_token", "expires_in": 0}`
)

func TestBearerPolicy_SuccessGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srvURL := srv.URL()
	cred := NewClientSecretCredential(tenantID, clientID, secret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	_, err := pipeline.Do(context.Background(), azcore.NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
}

func TestBearerPolicy_CredentialFailGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srvURL := srv.URL()
	cred := NewClientSecretCredential(tenantID, clientID, wrongSecret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	_, err := pipeline.Do(context.Background(), azcore.NewRequest(http.MethodGet, srv.URL()))
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
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

	srvURL := srv.URL()
	cred := NewClientSecretCredential(tenantID, clientID, secret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	req := azcore.NewRequest(http.MethodGet, srv.URL())
	_, err := pipeline.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
	_, err = pipeline.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
}

// with https scheme enabled we get an auth failed error which let's us test the is not retriable error
func TestRetryPolicy_IsNotRetriable(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srvURL := srv.URL()
	cred := NewClientSecretCredential(tenantID, clientID, wrongSecret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	_, err := pipeline.Do(context.Background(), azcore.NewRequest(http.MethodGet, srv.URL()))
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestRetryPolicy_HTTPRequest(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srvURL := srv.URL()
	cred := NewClientSecretCredential(tenantID, clientID, wrongSecret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	_, err := pipeline.Do(context.Background(), azcore.NewRequest(http.MethodGet, srv.URL()))
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}
