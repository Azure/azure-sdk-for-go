// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	msiScope                   = "https://storage.azure.com"
	appServiceTokenSuccessResp = `{"access_token": "new_token", "expires_on": "09/14/2017 00:00:00 PM +00:00", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
	expiresOnIntResp           = `{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": "1560974028", "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
)

func clearEnvVars(envVars ...string) {
	for _, ev := range envVars {
		_ = os.Setenv(ev, "")
	}
}
func TestManagedIdentityCredential_GetTokenInCloudShellLive(t *testing.T) {
	if len(os.Getenv("MSI_ENDPOINT")) == 0 {
		t.Skip()
	}
	msiCred, err := NewManagedIdentityCredential(clientID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInCloudShellMock(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	msiCred, err := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInCloudShellMockFail(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	msiCred, err := NewManagedIdentityCredential("", &ManagedIdentityCredentialOptions{HTTPClient: srv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceV20170901Mock(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceTokenSuccessResp)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	msiCred, err := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
	if tk.Token != "new_token" {
		t.Fatalf("Did not receive the correct token. Expected \"new_token\", Received: %s", tk.Token)
	}
	if msiCred.client.msiType != msiTypeAppServiceV20170901 {
		t.Fatalf("Failed to detect the correct MSI Environment. Expected: %d, Received: %d", msiTypeAppServiceV20170901, msiCred.client.msiType)
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceV20190801Mock(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceTokenSuccessResp)))
	_ = os.Setenv("IDENTITY_ENDPOINT", srv.URL())
	_ = os.Setenv("IDENTITY_HEADER", "header")
	defer clearEnvVars("IDENTITY_ENDPOINT", "IDENTITY_HEADER")
	msiCred, err := NewManagedIdentityCredential("", &ManagedIdentityCredentialOptions{HTTPClient: srv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
	if tk.Token != "new_token" {
		t.Fatalf("Did not receive the correct token. Expected \"new_token\", Received: %s", tk.Token)
	}
	if msiCred.client.msiType != msiTypeAppServiceV20190801 {
		t.Fatalf("Failed to detect the correct MSI Environment. Expected: %d, Received: %d", msiTypeAppServiceV20190801, msiCred.client.msiType)
	}
}

func TestManagedIdentityCredential_CreateAppServiceAuthRequestV20170901(t *testing.T) {
	// setting a dummy value for MSI_ENDPOINT in order to be able to get a ManagedIdentityCredential type in order
	// to test App Service authentication request creation.
	_ = os.Setenv("IDENTITY_ENDPOINT", "somevalue")
	_ = os.Setenv("IDENTITY_HEADER", "header")
	defer clearEnvVars("IDENTITY_ENDPOINT", "IDENTITY_HEADER")
	cred, err := NewManagedIdentityCredential(clientID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createAuthRequest(context.Background(), "", []string{msiScope})
	if err != nil {
		t.Fatal(err)
	}
	if req.Request.Header.Get("X-IDENTITY-HEADER") != "header" {
		t.Fatalf("Unexpected value for secret header")
	}
	reqQueryParams, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		t.Fatalf("Unable to parse App Service request query params: %v", err)
	}
	if reqQueryParams["api-version"][0] != "2019-08-01" {
		t.Fatalf("Unexpected App Service API version")
	}
	if reqQueryParams["resource"][0] != msiScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
}

func TestManagedIdentityCredential_CreateAppServiceAuthRequestV20190801(t *testing.T) {
	// setting a dummy value for MSI_ENDPOINT in order to be able to get a ManagedIdentityCredential type in order
	// to test App Service authentication request creation.
	_ = os.Setenv("MSI_ENDPOINT", "somevalue")
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	cred, err := NewManagedIdentityCredential(clientID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createAuthRequest(context.Background(), "", []string{msiScope})
	if err != nil {
		t.Fatal(err)
	}
	if req.Request.Header.Get("secret") != "secret" {
		t.Fatalf("Unexpected value for secret header")
	}
	reqQueryParams, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		t.Fatalf("Unable to parse App Service request query params: %v", err)
	}
	if reqQueryParams["api-version"][0] != "2017-09-01" {
		t.Fatalf("Unexpected App Service API version")
	}
	if reqQueryParams["resource"][0] != msiScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
}

func TestManagedIdentityCredential_CreateAccessTokenExpiresOnInt(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(expiresOnIntResp)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	msiCred, err := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceMockFail(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	_ = os.Setenv("MSI_SECRET", "secret")
	defer clearEnvVars("MSI_ENDPOINT", "MSI_SECRET")
	msiCred, err := NewManagedIdentityCredential("", &ManagedIdentityCredentialOptions{HTTPClient: srv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

// func TestManagedIdentityCredential_GetTokenIMDSMock(t *testing.T) {
// 	timeout := time.After(5 * time.Second)
// 	done := make(chan bool)
// 	go func() {
// 		err := resetEnvironmentVarsForTest()
// 		if err != nil {
// 			t.Fatalf("Unable to set environment variables")
// 		}
// 		srv, close := mock.NewServer()
// 		defer close()
// 		srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
// 		msiCred := NewManagedIdentityCredential("", &ManagedIdentityCredentialOptions{HTTPClient: srv})
// 		_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
// 		if err == nil {
// 			t.Fatalf("Cannot run IMDS test in this environment")
// 		}
// 		time.Sleep(550 * time.Millisecond)
// 		done <- true
// 	}()

// 	select {
// 	case <-timeout:
// 		t.Fatal("Test didn't finish in time")
// 	case <-done:
// 	}
// }

func TestManagedIdentityCredential_NewManagedIdentityCredentialFail(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	_ = os.Setenv("MSI_ENDPOINT", "https://t .com")
	defer clearEnvVars("MSI_ENDPOINT")
	cred, err := NewManagedIdentityCredential("", &ManagedIdentityCredentialOptions{HTTPClient: srv})
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestBearerPolicy_ManagedIdentityCredential(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	cred, err := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(nil),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{msiScope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatalf("Expected an empty error but receive: %v", err)
	}
}

func TestManagedIdentityCredential_GetTokenUnexpectedJSON(t *testing.T) {
	resetEnvironmentVarsForTest()
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespMalformed)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	msiCred, err := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err == nil {
		t.Fatalf("Expected a JSON marshal error but received nil")
	}
}

func TestManagedIdentityCredential_CreateIMDSAuthRequest(t *testing.T) {
	// setting a dummy value for MSI_ENDPOINT in order to be able to get a ManagedIdentityCredential type in order
	// to test IMDS authentication request creation.
	_ = os.Setenv("MSI_ENDPOINT", "somevalue")
	defer clearEnvVars("MSI_ENDPOINT")
	cred, err := NewManagedIdentityCredential(clientID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cred.client.endpoint = imdsEndpoint
	req, err := cred.client.createIMDSAuthRequest(context.Background(), []string{msiScope})
	if err != nil {
		t.Fatal(err)
	}
	if req.Request.Header.Get(azcore.HeaderMetadata) != "true" {
		t.Fatalf("Unexpected value for Content-Type header")
	}
	reqQueryParams, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		t.Fatalf("Unable to parse IMDS query params: %v", err)
	}
	if reqQueryParams["api-version"][0] != imdsAPIVersion {
		t.Fatalf("Unexpected IMDS API version")
	}
	if reqQueryParams["resource"][0] != msiScope {
		t.Fatalf("Unexpected resource in resource query param")
	}
	if u := req.Request.URL.String(); !strings.HasPrefix(u, imdsEndpoint) {
		t.Fatalf("Unexpected default authority host %s", u)
	}
	if req.Request.URL.Scheme != "http" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestManagedIdentityCredential_GetTokenEnvVar(t *testing.T) {
	resetEnvironmentVarsForTest()
	err := os.Setenv("AZURE_CLIENT_ID", "test_client_id")
	if err != nil {
		t.Fatal(err)
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	_ = os.Setenv("MSI_ENDPOINT", srv.URL())
	defer clearEnvVars("MSI_ENDPOINT")
	msiCred, err := NewManagedIdentityCredential("", &ManagedIdentityCredentialOptions{HTTPClient: srv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	at, err := msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
	if at.Token != "new_token" {
		t.Fatalf("Did not receive the correct access token")
	}
}
