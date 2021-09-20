// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"golang.org/x/net/http2"
)

func TestInteractiveBrowserCredential_InvalidTenantID(t *testing.T) {
	options := InteractiveBrowserCredentialOptions{}
	options.TenantID = badTenantID
	cred, err := NewInteractiveBrowserCredential(&options)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
	var errType *CredentialUnavailableError
	if !errors.As(err, &errType) {
		t.Fatalf("Did not receive a CredentialUnavailableError. Received: %t", err)
	}
}

func TestInteractiveBrowserCredential_CreateWithNilOptions(t *testing.T) {
	cred, err := NewInteractiveBrowserCredential(nil)
	if err != nil {
		t.Fatalf("Failed to create interactive browser credential: %v", err)
	}
	if cred.client.authorityHost != string(AzurePublicCloud) {
		t.Fatalf("Wrong authority host set. Expected: %s, Received: %s", AzurePublicCloud, cred.client.authorityHost)
	}
	if cred.options.ClientID != developerSignOnClientID {
		t.Fatalf("Wrong clientID set. Expected: %s, Received: %s", developerSignOnClientID, cred.options.ClientID)
	}
	if cred.options.TenantID != organizationsTenantID {
		t.Fatalf("Wrong tenantID set. Expected: %s, Received: %s", organizationsTenantID, cred.options.TenantID)
	}
}

func TestInteractiveBrowserCredential_GetTokenSuccess(t *testing.T) {
	srv, close := mock.NewTLSServer(mock.WithHTTP2Enabled(true))
	defer close()
	tr := &http.Transport{}
	if err := http2.ConfigureTransport(tr); err != nil {
		t.Fatalf("Failed to configure http2 transport: %v", err)
	}
	tr.TLSClientConfig.InsecureSkipVerify = true
	client := &http.Client{Transport: tr}
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := InteractiveBrowserCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.HTTPClient = client
	cred, err := NewInteractiveBrowserCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	authCodeReceiver = func(ctx context.Context, authorityHost string, opts *InteractiveBrowserCredentialOptions, scopes []string) (*interactiveConfig, error) {
		return &interactiveConfig{
			authCode:    "12345",
			redirectURI: srv.URL(),
		}, nil
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"https://storage.azure.com/.default"}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
	if tk.Token != "new_token" {
		t.Fatal("Received unexpected token")
	}
}

func TestInteractiveBrowserCredential_SetPort(t *testing.T) {
	srv, close := mock.NewTLSServer(mock.WithHTTP2Enabled(true))
	defer close()
	tr := &http.Transport{}
	if err := http2.ConfigureTransport(tr); err != nil {
		t.Fatalf("Failed to configure http2 transport: %v", err)
	}
	tr.TLSClientConfig.InsecureSkipVerify = true
	client := &http.Client{Transport: tr}
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := InteractiveBrowserCredentialOptions{}
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.HTTPClient = client
	options.Port = 8080
	cred, err := NewInteractiveBrowserCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	authCodeReceiver = func(ctx context.Context, authorityHost string, opts *InteractiveBrowserCredentialOptions, scopes []string) (*interactiveConfig, error) {
		if opts.Port != 8080 {
			t.Fatalf("Did not receive the correct port. Expected: %v, Received: %v", 8080, opts.Port)
		}
		return &interactiveConfig{
			authCode:    "12345",
			redirectURI: srv.URL(),
		}, nil
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"https://storage.azure.com/.default"}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
	if tk.Token != "new_token" {
		t.Fatal("Received unexpected token")
	}
}

func TestInteractiveBrowserCredential_GetTokenInvalidCredentials(t *testing.T) {
	srv, close := mock.NewTLSServer(mock.WithHTTP2Enabled(true))
	defer close()
	tr := &http.Transport{}
	if err := http2.ConfigureTransport(tr); err != nil {
		t.Fatalf("Failed to configure http2 transport: %v", err)
	}
	tr.TLSClientConfig.InsecureSkipVerify = true
	client := &http.Client{Transport: tr}
	srv.SetResponse(mock.WithBody([]byte(accessTokenRespError)), mock.WithStatusCode(http.StatusUnauthorized))
	options := InteractiveBrowserCredentialOptions{}
	options.ClientSecret = wrongSecret
	options.AuthorityHost = AuthorityHost(srv.URL())
	options.HTTPClient = client
	cred, err := NewInteractiveBrowserCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	authCodeReceiver = func(ctx context.Context, authorityHost string, opts *InteractiveBrowserCredentialOptions, scopes []string) (*interactiveConfig, error) {
		return &interactiveConfig{
			authCode:    "12345",
			redirectURI: srv.URL(),
		}, nil
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
	var respError *AADAuthenticationFailedError
	if !errors.As(authFailed.Unwrap(), &respError) {
		t.Fatalf("Expected: AADAuthenticationFailedError, Received: %T", err)
	}
	if len(respError.Message) == 0 {
		t.Fatalf("Did not receive an error message")
	}
	if len(respError.Description) == 0 {
		t.Fatalf("Did not receive an error description")
	}
	if len(respError.Timestamp) == 0 {
		t.Fatalf("Did not receive a timestamp")
	}
	if len(respError.TraceID) == 0 {
		t.Fatalf("Did not receive a TraceID")
	}
	if len(respError.CorrelationID) == 0 {
		t.Fatalf("Did not receive a CorrelationID")
	}
	if len(respError.URL) == 0 {
		t.Fatalf("Did not receive an error URL")
	}
	if respError.Response == nil {
		t.Fatalf("Did not receive an error response")
	}
}
