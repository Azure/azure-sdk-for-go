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

func TestManagedIdentityCredential_GetTokenInCloudShellLive(t *testing.T) {
	msiEndpoint := os.Getenv("MSI_ENDPOINT")
	if len(msiEndpoint) == 0 {
		t.Skip()
	}
	msiCred := NewManagedIdentityCredential(clientID, newDefaultManagedIdentityOptions())
	_, err := msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
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
	srv.AppendResponse(mock.WithBody([]byte(`{"access_token": "ey0....", "expires_in": 3600}`)))
	testURL := srv.URL()
	_ = os.Setenv("MSI_ENDPOINT", testURL.String())
	msiCred := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
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
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
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
	srv.AppendResponse(mock.WithBody([]byte(`{"access_token": "ey0....", "expires_in": 3600}`)))
	testURL := srv.URL()
	_ = os.Setenv("MSI_ENDPOINT", testURL.String())
	_ = os.Setenv("MSI_SECRET", "secret")
	msiCred := NewManagedIdentityCredential(clientID, &ManagedIdentityCredentialOptions{HTTPClient: srv})
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
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
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
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
// 		srv.AppendResponse(mock.WithBody([]byte(`{"access_token": "ey0....", "expires_in": 3600}`)))
// 		msiCred := NewManagedIdentityCredential("", newDefaultManagedIdentityOptions())
// 		msiCred.client = newTestManagedIdentityClient(&ManagedIdentityCredentialOptions{HTTPClient: srv})
// 		_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
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
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}
