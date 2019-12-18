// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	msiScope                   = "https://storage.azure.com"
	appServiceTokenSuccessResp = `{"access_token": "new_token", "expires_on": "09/14/2017 00:00:00 PM +00:00", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
	expiresOnIntResp           = `{"access_token": "new_token", "refresh_token": "", "expires_in": "", "expires_on": "1560974028", "not_before": "1560970130", "resource": "https://vault.azure.net", "token_type": "Bearer"}`
)

func TestManagedIdentityCredential_GetTokenInCloudShellLive(t *testing.T) {
	msiEndpoint := os.Getenv("MSI_ENDPOINT")
	if len(msiEndpoint) == 0 {
		t.Skip()
	}
	msiCred := NewManagedIdentityCredential(clientID, nil)
	_, err := msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInCloudShellMock(t *testing.T) {
	err := resetEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to set environment variables")
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	testURL := srv.URL()
	_ = os.Setenv("MSI_ENDPOINT", testURL.String())
	msiCred := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInCloudShellMockFail(t *testing.T) {
	err := resetEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to set environment variables")
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	testURL := srv.URL()
	_ = os.Setenv("MSI_ENDPOINT", testURL.String())
	msiCred := NewManagedIdentityCredential("", &ManagedIdentityCredentialOptions{HTTPClient: srv})
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceMock(t *testing.T) {
	err := resetEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to set environment variables")
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(appServiceTokenSuccessResp)))
	testURL := srv.URL()
	_ = os.Setenv("MSI_ENDPOINT", testURL.String())
	_ = os.Setenv("MSI_SECRET", "secret")
	msiCred := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_CreateAccessTokenExpiresOnInt(t *testing.T) {
	err := resetEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to set environment variables")
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(expiresOnIntResp)))
	testURL := srv.URL()
	_ = os.Setenv("MSI_ENDPOINT", testURL.String())
	_ = os.Setenv("MSI_SECRET", "secret")
	msiCred := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to retrieve a token")
	}
}

func TestManagedIdentityCredential_GetTokenInAppServiceMockFail(t *testing.T) {
	err := resetEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to set environment variables")
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	testURL := srv.URL()
	_ = os.Setenv("MSI_ENDPOINT", testURL.String())
	_ = os.Setenv("MSI_SECRET", "secret")
	msiCred := NewManagedIdentityCredential("", &ManagedIdentityCredentialOptions{HTTPClient: srv})
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

func TestManagedIdentityCredential_GetTokenUnknownFail(t *testing.T) {
	err := resetEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to set environment variables")
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	_ = os.Setenv("MSI_ENDPOINT", "https://t .com")
	msiCred := NewManagedIdentityCredential("", &ManagedIdentityCredentialOptions{HTTPClient: srv})
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestBearerPolicy_ManagedIdentityCredential(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	testURL := srv.URL()
	_ = os.Setenv("MSI_ENDPOINT", testURL.String())
	cred := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{msiScope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	req := pipeline.NewRequest(http.MethodGet, srv.URL())
	_, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("Expected an empty error but receive: %v", err)
	}
}

func TestManagedIdentityCredential_GetTokenUnexpectedJSON(t *testing.T) {
	err := resetEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to set environment variables")
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespMalformed)))
	testURL := srv.URL()
	_ = os.Setenv("MSI_ENDPOINT", testURL.String())
	msiCred := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{msiScope}})
	if err == nil {
		t.Fatalf("Expected a JSON marshal error but received nil")
	}
}

func TestManagedIdentityCredential_CreateIMDSAuthRequest(t *testing.T) {
	cred := NewManagedIdentityCredential(clientID, nil)
	cred.client.endpoint = imdsURL
	req := cred.client.createIMDSAuthRequest([]string{msiScope})
	if req.Request.Header.Get(azcore.HeaderMetadata) != "true" {
		t.Fatalf("Unexpected value for Content-Type header")
	}
	// queryValues, err :=
	// if req.Request.URL.Query["api-version"] != imdsAPIVersion {
	// 	t.Fatalf("Unexpected IMDS API version")
	// }
	// if reqQueryParams["resource"][0] != msiScope {
	// 	t.Fatalf("Unexpected resource in resource query param")
	// }
	// if req.Request.URL.Host != imdsURL.Host {
	// 	t.Fatalf("Unexpected default authority host")
	// }
	// if req.Request.URL.Scheme != "https" {
	// 	t.Fatalf("Wrong request scheme")
	// }
}
