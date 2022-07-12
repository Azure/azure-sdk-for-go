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
}
