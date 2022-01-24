// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const Assertion = "SomeJWTToken"

func TestClientAssertionCredential_InvalidTenantID(t *testing.T) {
	cred, err := NewClientAssertionCredential(badTenantID, fakeClientID, Assertion, nil)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestClientAssertionCredential_GetTokenSuccess(t *testing.T) {
	cred, err := NewClientAssertionCredential(fakeTenantID, fakeClientID, Assertion, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred.client = fakeConfidentialClient{}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
}

//
// TODO - this test needs to configure Workload Identity Federated token on AzureAD against tests are run to
//func TestClientAssertionCredential_Live(t *testing.T) {
//	opts, stop := initRecording(t)
//	defer stop()
//	o := ClientAssertionCredentialOptions{ClientOptions: opts}
//	cred, err := NewClientAssertionCredential(liveSP.tenantID, liveSP.clientID, liveSP.clientAssertion, &o)
//	if err != nil {
//		t.Fatalf("failed to construct credential: %v", err)
//	}
//	testGetTokenSuccess(t, cred)
//}

func TestClientAssertionCredential_InvalidAssertionLive(t *testing.T) {
	opts, stop := initRecording(t)
	defer stop()
	o := ClientAssertionCredentialOptions{ClientOptions: opts}
	cred, err := NewClientAssertionCredential(liveSP.tenantID, liveSP.clientID, "invalid Assertion", &o)
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if tk != nil {
		t.Fatal("GetToken returned a token")
	}
	var e AuthenticationFailedError
	if !errors.As(err, &e) {
		t.Fatal("expected AuthenticationFailedError")
	}
	if e.RawResponse == nil {
		t.Fatal("expected a non-nil RawResponse")
	}
}
