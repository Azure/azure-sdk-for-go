//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

func TestDeviceCodeCredential_InvalidTenantID(t *testing.T) {
	options := DeviceCodeCredentialOptions{}
	options.TenantID = badTenantID
	cred, err := NewDeviceCodeCredential(&options)
	if err == nil {
		t.Fatal("Expected an error but received none")
	}
	if cred != nil {
		t.Fatalf("Expected a nil credential value. Received: %v", cred)
	}
}

func TestDeviceCodeCredential_GetTokenInvalidCredentials(t *testing.T) {
	cred, err := NewDeviceCodeCredential(nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred.client = fakePublicClient{err: errors.New("invalid credentials")}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}

func TestDeviceCodeCredential_UserPromptError(t *testing.T) {
	expectedCtx := context.WithValue(context.Background(), struct{}{}, "")
	expected := DeviceCodeMessage{UserCode: "user code", VerificationURL: "http://localhost", Message: "message"}
	success := "it worked"
	options := DeviceCodeCredentialOptions{
		ClientID: fakeClientID,
		TenantID: fakeTenantID,
		UserPrompt: func(ctx context.Context, m DeviceCodeMessage) error {
			if ctx != expectedCtx {
				t.Fatal("UserPrompt received unexpected Context")
			}
			if m.Message != expected.Message {
				t.Fatalf(`unexpected Message "%s"`, m.Message)
			}
			if m.UserCode != expected.UserCode {
				t.Fatalf(`unexpected UserCode "%s"`, m.UserCode)
			}
			if m.VerificationURL != expected.VerificationURL {
				t.Fatalf(`unexpected VerificationURL "%s"`, m.VerificationURL)
			}
			return errors.New(success)
		},
	}
	cred, err := NewDeviceCodeCredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential: %v", err)
	}
	cred.client = fakePublicClient{
		dc: public.DeviceCode{
			Result: public.DeviceCodeResult{
				Message:         expected.Message,
				UserCode:        expected.UserCode,
				VerificationURL: expected.VerificationURL,
			},
		},
	}
	_, err = cred.GetToken(expectedCtx, policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatal("expected an error")
	}
	if err.Error() != success {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestDeviceCodeCredential_Live(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		t.Skip("this test requires manual recording and can't pass live in CI")
	}
	o, stop := initRecording(t)
	defer stop()
	opts := DeviceCodeCredentialOptions{TenantID: liveUser.tenantID, ClientOptions: o}
	if recording.GetRecordMode() == recording.PlaybackMode {
		opts.UserPrompt = func(ctx context.Context, m DeviceCodeMessage) error { return nil }
	}
	cred, err := NewDeviceCodeCredential(&opts)
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred)
}
