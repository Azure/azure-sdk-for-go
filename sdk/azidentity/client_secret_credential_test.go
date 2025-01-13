//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
)

const fakeSecret = "secret"

func TestClientSecretCredential_Live(t *testing.T) {
	if recording.GetRecordMode() == recording.LiveMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22879")
	}
	for _, disabledID := range []bool{true, false} {
		name := "default options"
		if disabledID {
			name = "instance discovery disabled"
		}
		t.Run(name, func(t *testing.T) {
			opts, stop := initRecording(t)
			defer stop()
			o := ClientSecretCredentialOptions{ClientOptions: opts, DisableInstanceDiscovery: disabledID}
			cred, err := NewClientSecretCredential(liveSP.tenantID, liveSP.clientID, liveSP.secret, &o)
			if err != nil {
				t.Fatalf("failed to construct credential: %v", err)
			}
			testGetTokenSuccess(t, cred)
		})
	}
}

func TestClientSecretCredentialADFS_Live(t *testing.T) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		if adfsLiveSP.clientID == "" || adfsLiveSP.secret == "" || adfsScope == "" {
			t.Skip("set ADFS_SP_* environment variables to run this test live")
		}
	}
	opts, stop := initRecording(t)
	defer stop()
	opts.Cloud.ActiveDirectoryAuthorityHost = adfsAuthority
	o := ClientSecretCredentialOptions{ClientOptions: opts, DisableInstanceDiscovery: true}
	cred, err := NewClientSecretCredential("adfs", adfsLiveSP.clientID, adfsLiveSP.secret, &o)
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	testGetTokenSuccess(t, cred, adfsScope)
}

func TestClientSecretCredential_InvalidSecretLive(t *testing.T) {
	opts, stop := initRecording(t)
	defer stop()
	o := ClientSecretCredentialOptions{ClientOptions: opts}
	cred, err := NewClientSecretCredential(liveSP.tenantID, liveSP.clientID, "invalid secret", &o)
	if err != nil {
		t.Fatalf("failed to construct credential: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), testTRO)
	if !reflect.ValueOf(tk).IsZero() {
		t.Fatal("expected a zero value AccessToken")
	}
	if e, ok := err.(*AuthenticationFailedError); ok {
		if e.RawResponse == nil {
			t.Fatal("expected a non-nil RawResponse")
		}
	} else {
		t.Fatalf("expected AuthenticationFailedError, received %T", err)
	}
	if !strings.HasPrefix(err.Error(), credNameSecret) {
		t.Fatal("missing credential type prefix")
	}
}
