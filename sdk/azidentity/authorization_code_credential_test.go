// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

func TestAuthorizationCodeCredential_InvalidTenantID(t *testing.T) {
	cred, err := NewAuthorizationCodeCredential(badTenantID, fakeClientID, "fake_auth_code", "http://localhost", nil)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestAuthorizationCodeCredential_ConfidentialLive(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		t.Skip("this test requires manual recording and can't pass live in CI")
	}
	opts, stop := initRecording(t)
	defer stop()
	o := AuthorizationCodeCredentialOptions{ClientOptions: opts, ClientSecret: liveSP.secret}
	cred, err := NewAuthorizationCodeCredential(liveSP.tenantID, liveSP.clientID, "fake_auth_code", "http://localhost", &o)
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred)
}

func TestAuthorizationCodeCredential_ConfidentialInvalidCodeLive(t *testing.T) {
	opts, stop := initRecording(t)
	defer stop()
	o := AuthorizationCodeCredentialOptions{ClientOptions: opts, ClientSecret: liveSP.secret}
	cred, err := NewAuthorizationCodeCredential(liveSP.tenantID, liveSP.clientID, "fake_auth_code", "http://localhost", &o)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if tk != nil {
		t.Fatal("GetToken returned a token")
	}
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestAuthorizationCodeCredential_PublicLive(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		t.Skip("this test requires manual recording and can't pass live in CI")
	}
	opts, stop := initRecording(t)
	defer stop()
	o := AuthorizationCodeCredentialOptions{ClientOptions: opts}
	cred, err := NewAuthorizationCodeCredential(liveSP.tenantID, liveSP.clientID, "fake_auth_code", "http://localhost", &o)
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred)
}

func TestAuthorizationCodeCredential_PublicInvalidCodeLive(t *testing.T) {
	opts, stop := initRecording(t)
	defer stop()
	o := AuthorizationCodeCredentialOptions{ClientOptions: opts}
	cred, err := NewAuthorizationCodeCredential(liveSP.tenantID, liveSP.clientID, "fake_auth_code", "http://localhost", &o)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if tk != nil {
		t.Fatal("GetToken returned a token")
	}
	if err == nil {
		t.Fatal("expected an error")
	}
}
