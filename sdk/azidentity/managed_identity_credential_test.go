// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

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
	msiCred := NewManagedIdentityCredential(clientID, newDefaultManagedIdentityOptions())
	msiCred.client = newTestManagedIdentityClient(&ManagedIdentityCredentialOptions{HTTPClient: srv})
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
	msiCred := NewManagedIdentityCredential("", newDefaultManagedIdentityOptions())
	msiCred.client = newTestManagedIdentityClient(&ManagedIdentityCredentialOptions{HTTPClient: srv})
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
	msiCred := NewManagedIdentityCredential(clientID, newDefaultManagedIdentityOptions())
	msiCred.client = newTestManagedIdentityClient(&ManagedIdentityCredentialOptions{HTTPClient: srv})
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
	msiCred := NewManagedIdentityCredential("", newDefaultManagedIdentityOptions())
	msiCred.client = newTestManagedIdentityClient(&ManagedIdentityCredentialOptions{HTTPClient: srv})
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
	// testURL := srv.URL()
	_ = os.Setenv("MSI_ENDPOINT", "https://t .com")
	msiCred := NewManagedIdentityCredential("", newDefaultManagedIdentityOptions())
	msiCred.client = newTestManagedIdentityClient(&ManagedIdentityCredentialOptions{HTTPClient: srv})
	_, err = msiCred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func newTestManagedIdentityClient(options *ManagedIdentityCredentialOptions) *managedIdentityClient {
	options = options.setDefaultValues()
	// TODO document the use of these variables
	return &managedIdentityClient{
		pipeline:               newDefaultTestMSIPipeline(options),
		imdsAPIVersion:         imdsAPIVersion,
		imdsAvailableTimeoutMS: 500,
		msiType:                unknown,
	}
}

// NewDefaultMSIPipeline creates a Pipeline using the specified pipeline options needed
// for a Managed Identity, such as a MSI specific retry policy
func newDefaultTestMSIPipeline(o *ManagedIdentityCredentialOptions) azcore.Pipeline {
	if o.HTTPClient == nil {
		o.HTTPClient = azcore.DefaultHTTPClientTransport()
	}

	// retry policy for MSI is not end-user configurable
	retryOpts := azcore.RetryOptions{
		MaxTries:   5,
		RetryDelay: 2 * time.Second,
		StatusCodes: append(azcore.StatusCodesForRetry[:],
			http.StatusNotFound,
			http.StatusGone,
			// all remaining 5xx
			http.StatusNotImplemented,
			http.StatusHTTPVersionNotSupported,
			http.StatusVariantAlsoNegotiates,
			http.StatusInsufficientStorage,
			http.StatusLoopDetected,
			http.StatusNotExtended,
			http.StatusNetworkAuthenticationRequired),
	}

	return azcore.NewPipeline(
		o.HTTPClient,
		azcore.NewTelemetryPolicy(o.Telemetry),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(retryOpts),
		azcore.NewRequestLogPolicy(o.LogOptions))
}
