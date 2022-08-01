//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func TestUsernamePasswordCredential_InvalidTenantID(t *testing.T) {
	cred, err := NewUsernamePasswordCredential(badTenantID, fakeClientID, "username", "password", nil)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestUsernamePasswordCredential_GetTokenSuccess(t *testing.T) {
	cred, err := NewUsernamePasswordCredential(fakeTenantID, fakeClientID, "username", "password", nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred.client = fakePublicClient{}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestUsernamePasswordCredential_Live(t *testing.T) {
	o, stop := initRecording(t)
	defer stop()
	opts := UsernamePasswordCredentialOptions{ClientOptions: o}
	cred, err := NewUsernamePasswordCredential(liveUser.tenantID, developerSignOnClientID, liveUser.username, liveUser.password, &opts)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	testGetTokenSuccess(t, cred)
}

func TestUsernamePasswordCredential_InvalidPasswordLive(t *testing.T) {
	o, stop := initRecording(t)
	defer stop()
	opts := UsernamePasswordCredentialOptions{ClientOptions: o}
	cred, err := NewUsernamePasswordCredential(liveUser.tenantID, developerSignOnClientID, liveUser.username, "invalid password", &opts)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if !reflect.ValueOf(tk).IsZero() {
		t.Fatal("expected a zero value AccessToken")
	}
	var e *AuthenticationFailedError
	if !errors.As(err, &e) {
		t.Fatal("expected AuthenticationFailedError")
	}
	if e.RawResponse == nil {
		t.Fatal("expected a non-nil RawResponse")
	}
	expectedErrorMsg := "UsernamePasswordCredential authentication failed\nPOST https://localhost:5001/fake-tenant/oauth2/v2.0/token\n--------------------------------------------------------------------------------\nRESPONSE 400 Bad Request\n--------------------------------------------------------------------------------\n{\n  \"error\": \"invalid_grant\",\n  \"error_description\": \"AADSTS50126: Error validating credentials due to invalid username or password.\\r\\nTrace ID: be3e1080-9e7f-418d-9c32-dc42cdf2d400\\r\\nCorrelation ID: e4caccbc-1ab2-43f8-a4b4-c74b9a94c87a\\r\\nTimestamp: 2022-01-24 22:41:16Z\",\n  \"error_codes\": [\n    50126\n  ],\n  \"timestamp\": \"2022-01-24 22:41:16Z\",\n  \"trace_id\": \"be3e1080-9e7f-418d-9c32-dc42cdf2d400\",\n  \"correlation_id\": \"e4caccbc-1ab2-43f8-a4b4-c74b9a94c87a\",\n  \"error_uri\": \"https://login.microsoftonline.com/error?code=50126\"\n}\n--------------------------------------------------------------------------------\nTo troubleshoot, visit https://aka.ms/azsdk/go/identity/troubleshoot#username-password"
	if e.Error() != expectedErrorMsg {
		t.Fatal("unexpected error message")
	}
}
