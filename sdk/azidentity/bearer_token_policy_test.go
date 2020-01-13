// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	accessTokenRespError      = `{"error": "invalid_client","error_description": "Invalid client secret is provided.","error_codes": [7000215],"timestamp": "2019-12-01 19:00:00Z","trace_id": "2d091b0","correlation_id": "a999","error_uri": "https://login.contoso.com/error?code=7000215"}`
	accessTokenRespSuccess    = `{"access_token": "` + tokenValue + `", "expires_in": 3600}`
	accessTokenAppServSuccess = `{"access_token": "` + tokenValue + `","expires_on": "06/20/2019 02:57:58 +00:00","resource": "https://vault.azure.net","token_type": "Bearer"}`
	accessTokenRespMalformed  = `{"access_token": 0, "expires_in": 3600}`
	accessTokenRespShortLived = `{"access_token": "` + tokenValue + `", "expires_in": 0}`
)

func TestBearerPolicy_SuccessGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srvURL := srv.URL()
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	resp, err := pipeline.Do(context.Background(), azcore.NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("Expected nil error but received one")
	}
	const expectedToken = bearerTokenPrefix + tokenValue
	if token := resp.Request.Header.Get(azcore.HeaderAuthorization); token != expectedToken {
		t.Fatalf("expected token '%s', got '%s'", expectedToken, token)
	}
}

func TestBearerPolicy_CredentialFailGetToken(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srvURL := srv.URL()
	cred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	resp, err := pipeline.Do(context.Background(), azcore.NewRequest(http.MethodGet, srv.URL()))
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

	srvURL := srv.URL()
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	req := azcore.NewRequest(http.MethodGet, srv.URL())
	_, err = pipeline.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	_, err = pipeline.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

// with https scheme enabled we get an auth failed error which let's us test the is not retriable error
func TestRetryPolicy_IsNotRetriable(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srvURL := srv.URL()
	cred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	_, err = pipeline.Do(context.Background(), azcore.NewRequest(http.MethodGet, srv.URL()))
	var afe *AuthenticationFailedError
	if !errors.As(err, &afe) {
		t.Fatalf("unexpected error type %v", err)
	}
}

func TestRetryPolicy_HTTPRequest(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srvURL := srv.URL()
	cred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	_, err = pipeline.Do(context.Background(), azcore.NewRequest(http.MethodGet, srv.URL()))
	var afe *AuthenticationFailedError
	if !errors.As(err, &afe) {
		t.Fatalf("unexpected error type %v", err)
	}
}
