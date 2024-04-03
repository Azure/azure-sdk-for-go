//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

func TestDeviceCodeCredential_GetTokenInvalidCredentials(t *testing.T) {
	cred, err := NewDeviceCodeCredential(nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred.client.noCAE = fakePublicClient{err: errors.New("invalid credentials")}
	_, err = cred.GetToken(context.Background(), testTRO)
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
	cred.client.noCAE = fakePublicClient{
		dc: public.DeviceCode{
			Result: public.DeviceCodeResult{
				Message:         expected.Message,
				UserCode:        expected.UserCode,
				VerificationURL: expected.VerificationURL,
			},
		},
	}
	_, err = cred.GetToken(expectedCtx, testTRO)
	if err == nil {
		t.Fatal("expected an error")
	}
	if expected := fmt.Sprintf("%s: %s", credNameDeviceCode, success); err.Error() != expected {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestDeviceCodeCredential_Live(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode && !runManualTests {
		t.Skipf("set %s to run this test", azidentityRunManualTests)
	}
	for _, test := range []struct {
		clientID, desc, tenantID string
		opts                     DeviceCodeCredentialOptions
	}{
		{
			desc: "default options",
		},
		{
			desc: "instance discovery disabled",
			opts: DeviceCodeCredentialOptions{DisableInstanceDiscovery: true, TenantID: liveSP.tenantID},
		},
		{
			desc: "optional tenant",
			opts: DeviceCodeCredentialOptions{TenantID: liveSP.tenantID},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			o, stop := initRecording(t)
			defer stop()
			test.opts.ClientOptions = o
			if recording.GetRecordMode() == recording.PlaybackMode {
				test.opts.UserPrompt = func(ctx context.Context, m DeviceCodeMessage) error { return nil }
			}
			cred, err := NewDeviceCodeCredential(&test.opts)
			if err != nil {
				t.Fatal(err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestDeviceCodeCredentialADFS_Live(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode && !runManualTests {
		t.Skipf("set %s to run this test", azidentityRunManualTests)
	}
	if adfsLiveSP.clientID == "" {
		t.Skip("set ADFS_SP_* environment variables to run this test")
	}
	o, stop := initRecording(t)
	defer stop()
	o.Cloud.ActiveDirectoryAuthorityHost = adfsAuthority
	opts := DeviceCodeCredentialOptions{
		ClientID:                 adfsLiveUser.clientID,
		ClientOptions:            o,
		DisableInstanceDiscovery: true,
		TenantID:                 "adfs",
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		opts.UserPrompt = func(ctx context.Context, m DeviceCodeMessage) error { return nil }
	}
	cred, err := NewDeviceCodeCredential(&opts)
	if err != nil {
		t.Fatal(err)
	}
	testGetTokenSuccess(t, cred, adfsScope)
}
